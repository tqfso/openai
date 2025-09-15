package znode

import (
	"encoding/binary"
)

// blake2bG implements the G mixing function for BLAKE2b
// It mixes four 64-bit words of the internal state according to the BLAKE2 specification
// Parameters:
//   - r, i: Current round and indexing parameters for selecting message words
//   - a, b, c, d: Pointers to the four words to be mixed
//   - m: Array of message words
func blake2bG(r, i int, a, b, c, d *uint64, m []uint64) {
	*a = *a + *b + m[blake2bSigma[r][2*i]]   // a = a + b + m[sigma]
	*d = rotr64(*d^*a, 32)                   // d = (d ^ a) >>> 32
	*c = *c + *d                             // c = c + d
	*b = rotr64(*b^*c, 24)                   // b = (b ^ c) >>> 24
	*a = *a + *b + m[blake2bSigma[r][2*i+1]] // a = a + b + m[sigma]
	*d = rotr64(*d^*a, 16)                   // d = (d ^ a) >>> 16
	*c = *c + *d                             // c = c + d
	*b = rotr64(*b^*c, 63)                   // b = (b ^ c) >>> 63
}

// blake2bRound performs a full mixing round on the internal state for BLAKE2b
// A round consists of eight G function calls in a specific pattern
// Parameters:
//   - r: Current round number (used for message word selection)
//   - v: Pointer to the internal state (16 words)
//   - m: Array of message words
func blake2bRound(r int, v *[16]uint64, m []uint64) {
	// Column step - mix the four columns
	blake2bG(r, 0, &v[0], &v[4], &v[8], &v[12], m)
	blake2bG(r, 1, &v[1], &v[5], &v[9], &v[13], m)
	blake2bG(r, 2, &v[2], &v[6], &v[10], &v[14], m)
	blake2bG(r, 3, &v[3], &v[7], &v[11], &v[15], m)

	// Diagonal step - mix the four diagonals
	blake2bG(r, 4, &v[0], &v[5], &v[10], &v[15], m)
	blake2bG(r, 5, &v[1], &v[6], &v[11], &v[12], m)
	blake2bG(r, 6, &v[2], &v[7], &v[8], &v[13], m)
	blake2bG(r, 7, &v[3], &v[4], &v[9], &v[14], m)
}

// blake2bInit0 initializes a BLAKE2b state with default values
// This sets up the state with the IV values and zeros for other fields
// Parameter:
//   - S: Pointer to the state to initialize
func blake2bInit0(S *blake2bState) {
	// Initialize with IV values
	for i := 0; i < 8; i++ {
		S.h[i] = blake2bIV[i]
	}

	// Zero other fields
	S.t[0] = 0
	S.t[1] = 0
	S.f[0] = 0
	S.f[1] = 0
	S.buflen = 0
	S.lastNode = 0
}

// blake2bCompress performs the compression function of BLAKE2b
// It processes a full block of data and updates the internal state
// Parameters:
//   - S: Pointer to the current state
//   - block: The 128-byte block of data to process
func blake2bCompress(S *blake2bState, block []byte) {
	m := make([]uint64, 16) // Message words
	v := [16]uint64{}       // Working state

	// Convert block bytes to 64-bit words (little-endian)
	for i := 0; i < 16; i++ {
		m[i] = binary.LittleEndian.Uint64(block[i*8:])
	}

	// Initialize working state v[0..15]
	for i := 0; i < 8; i++ {
		v[i] = S.h[i] // First half from current state
	}

	// Second half from IV, with counter and flags mixed in
	v[8] = blake2bIV[0]
	v[9] = blake2bIV[1]
	v[10] = blake2bIV[2]
	v[11] = blake2bIV[3]
	v[12] = S.t[0] ^ blake2bIV[4] // Mix in counter low bits
	v[13] = S.t[1] ^ blake2bIV[5] // Mix in counter high bits
	v[14] = S.f[0] ^ blake2bIV[6] // Mix in finalization flag 0
	v[15] = S.f[1] ^ blake2bIV[7] // Mix in finalization flag 1

	// Perform 12 mixing rounds
	for r := 0; r < 12; r++ {
		blake2bRound(r, &v, m)
	}

	// Update the state with the result
	for i := 0; i < 8; i++ {
		S.h[i] = S.h[i] ^ v[i] ^ v[i+8]
	}
}

// blake2bIncrementCounter adds the specified number of bytes to the counter
// Parameters:
//   - S: Pointer to the state
//   - inc: Number of bytes to add to the counter
func blake2bIncrementCounter(S *blake2bState, inc uint64) {
	S.t[0] += inc
	// Handle carry to high word if necessary
	if S.t[0] < inc {
		S.t[1]++
	}
}

// blake2bIsLastBlock checks if the final block flag is set
// Parameter:
//   - S: Pointer to the state
//
// Returns true if the finalization flag is set
func blake2bIsLastBlock(S *blake2bState) bool {
	return S.f[0] != 0
}

// blake2bSetLastBlock sets the finalization flags
// This indicates that the current block is the last one
// Parameter:
//   - S: Pointer to the state
func blake2bSetLastBlock(S *blake2bState) {
	// Set last node flag if necessary
	if S.lastNode != 0 {
		S.f[1] = 0xFFFFFFFFFFFFFFFF
	}
	// Always set last block flag
	S.f[0] = 0xFFFFFFFFFFFFFFFF
}

// blake2bUpdate processes input data in chunks
// It feeds data into the hash state block by block
// Parameters:
//   - S: Pointer to the state
//   - in: Input data to process
//
// Returns 0 on success
func blake2bUpdate(S *blake2bState, in []byte) int {
	inlen := len(in)

	if inlen > 0 {
		left := S.buflen                          // Bytes already in buffer
		fill := uint64(BLAKE2B_BLOCKBYTES) - left // Space left in buffer

		// Handle case where input fills the buffer
		if uint64(inlen) > fill {
			S.buflen = 0 // Buffer will be emptied

			// Fill buffer with start of input
			copy(S.buf[left:], in[:fill])

			// Process the now full buffer
			blake2bIncrementCounter(S, BLAKE2B_BLOCKBYTES)
			blake2bCompress(S, S.buf[:])

			// Move past the processed data
			in = in[fill:]
			inlen -= int(fill)

			// Process full blocks directly from input
			for inlen > BLAKE2B_BLOCKBYTES {
				blake2bIncrementCounter(S, BLAKE2B_BLOCKBYTES)
				blake2bCompress(S, in[:BLAKE2B_BLOCKBYTES])
				in = in[BLAKE2B_BLOCKBYTES:]
				inlen -= BLAKE2B_BLOCKBYTES
			}
		}

		// Store any remaining input in the buffer
		copy(S.buf[S.buflen:], in[:inlen])
		S.buflen += uint64(inlen)
	}

	return 0 // Success
}

// blake2bFinal finalizes the hash computation
// It processes any remaining data, sets finalization flags, and outputs the digest
// Parameters:
//   - S: Pointer to the state
//   - out: Output buffer for the digest
//   - outlen: Size of the desired output in bytes
//
// Returns 0 on success, negative on error
func blake2bFinal(S *blake2bState, out []byte, outlen int) int {
	// Check for valid output parameters
	if out == nil || outlen < int(S.outlen) || len(out) < outlen {
		return -1
	}

	// Check if already finalized
	if blake2bIsLastBlock(S) {
		return -1
	}

	// Count remaining bytes in buffer
	blake2bIncrementCounter(S, S.buflen)

	// Set finalization flags
	blake2bSetLastBlock(S)

	// Pad buffer with zeros
	for i := S.buflen; i < BLAKE2B_BLOCKBYTES; i++ {
		S.buf[i] = 0
	}

	// Process the final block
	blake2bCompress(S, S.buf[:])

	// Ensure we don't write past end of output buffer
	hashSize := int(S.outlen)
	if outlen < hashSize {
		hashSize = outlen
	}

	// Write digest to output buffer, little-endian format
	// Only write up to min(outlen, S.outlen) bytes
	for i := 0; i < (hashSize+7)/8; i++ {
		// Make sure we don't write past the end of the output buffer
		if i*8 < hashSize {
			remaining := hashSize - i*8
			if remaining >= 8 {
				binary.LittleEndian.PutUint64(out[i*8:], S.h[i])
			} else {
				// Partial word at the end
				var wordBytes [8]byte
				binary.LittleEndian.PutUint64(wordBytes[:], S.h[i])
				copy(out[i*8:], wordBytes[:remaining])
			}
		}
	}

	return 0 // Success
}

// blake2bInit initializes a BLAKE2b state with the desired output size
// Parameters:
//   - S: Pointer to the state to initialize
//   - outlen: Desired digest length in bytes (1-64)
//
// Returns 0 on success, negative on error
func blake2bInit(S *blake2bState, outlen uint8) int {
	// Check for valid output length
	if outlen == 0 || outlen > BLAKE2B_OUTBYTES {
		return -1
	}

	// Initialize with default values
	blake2bInit0(S)

	// Set desired output length and customize h[0]
	S.outlen = uint64(outlen)
	S.h[0] ^= uint64(outlen) | (1 << 16) | (1 << 24)

	return 0 // Success
}

// blake2b performs the complete BLAKE2b hash computation
// This is a convenience function that handles initialization, data processing, and finalization
// Parameters:
//   - msg: Input message to hash
//   - out: Output buffer for the digest
//   - outlen: Size of the desired output in bytes (1-64)
//
// Returns 0 on success, negative on error
func blake2b(msg []byte, out []byte, outlen int) int {
	// Check for valid output parameters
	if outlen <= 0 || outlen > BLAKE2B_OUTBYTES || out == nil || len(out) < outlen {
		return -1
	}

	var S blake2bState

	// Initialize state
	if blake2bInit(&S, uint8(outlen)) < 0 {
		return -1
	}

	// Process message
	if blake2bUpdate(&S, msg) < 0 {
		return -1
	}

	// Finalize and get digest
	if blake2bFinal(&S, out, outlen) < 0 {
		return -1
	}

	return 0 // Success
}

// blake2b_512A computes a BLAKE2b-512 hash of the input data
// This is a convenience wrapper for BLAKE2b with a 64-byte output
// Parameters:
//   - data: Input data buffer
//   - data_len: Length of data to hash
//   - digest: Output buffer for the hash result (must be at least 64 bytes)
//
// Returns 0 on success, negative on error
func blake2b_512A(data []byte, data_len int, digest []byte) int {
	// Validate input
	if data_len < 0 || data_len > len(data) || len(digest) < 64 {
		return -1
	}
	return blake2b(data[:data_len], digest, 64)
}
