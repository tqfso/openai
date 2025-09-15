package znode

import (
	"encoding/binary"
)

// blake2sG implements the G mixing function for BLAKE2s
// It mixes four 32-bit words of the internal state according to the BLAKE2 specification
// Parameters:
//   - r, i: Current round and indexing parameters for selecting message words
//   - a, b, c, d: Pointers to the four words to be mixed
//   - m: Array of message words
func blake2sG(r, i int, a, b, c, d *uint32, m []uint32) {
	*a = *a + *b + m[blake2sSigma[r][2*i]]   // a = a + b + m[sigma]
	*d = rotr32(*d^*a, 16)                   // d = (d ^ a) >>> 16
	*c = *c + *d                             // c = c + d
	*b = rotr32(*b^*c, 12)                   // b = (b ^ c) >>> 12
	*a = *a + *b + m[blake2sSigma[r][2*i+1]] // a = a + b + m[sigma]
	*d = rotr32(*d^*a, 8)                    // d = (d ^ a) >>> 8
	*c = *c + *d                             // c = c + d
	*b = rotr32(*b^*c, 7)                    // b = (b ^ c) >>> 7
}

// blake2sRound performs a full mixing round on the internal state for BLAKE2s
// A round consists of eight G function calls in a specific pattern
// Parameters:
//   - r: Current round number (used for message word selection)
//   - v: Pointer to the internal state (16 words)
//   - m: Array of message words
func blake2sRound(r int, v *[16]uint32, m []uint32) {
	// Column step - mix the four columns
	blake2sG(r, 0, &v[0], &v[4], &v[8], &v[12], m)
	blake2sG(r, 1, &v[1], &v[5], &v[9], &v[13], m)
	blake2sG(r, 2, &v[2], &v[6], &v[10], &v[14], m)
	blake2sG(r, 3, &v[3], &v[7], &v[11], &v[15], m)

	// Diagonal step - mix the four diagonals
	blake2sG(r, 4, &v[0], &v[5], &v[10], &v[15], m)
	blake2sG(r, 5, &v[1], &v[6], &v[11], &v[12], m)
	blake2sG(r, 6, &v[2], &v[7], &v[8], &v[13], m)
	blake2sG(r, 7, &v[3], &v[4], &v[9], &v[14], m)
}

// blake2sInit0 initializes a BLAKE2s state with default values
// This sets up the state with the IV values and zeros for other fields
// Parameter:
//   - S: Pointer to the state to initialize
func blake2sInit0(S *blake2sState) {
	// Initialize with IV values
	for i := 0; i < 8; i++ {
		S.h[i] = blake2sIV[i]
	}

	// Zero other fields
	S.t[0] = 0
	S.t[1] = 0
	S.f[0] = 0
	S.f[1] = 0
	S.buflen = 0
	S.lastNode = 0
}

// blake2sCompress performs the compression function of BLAKE2s
// It processes a full block of data and updates the internal state
// Parameters:
//   - S: Pointer to the current state
//   - block: The 64-byte block of data to process
func blake2sCompress(S *blake2sState, block []byte) {
	m := make([]uint32, 16) // Message words
	v := [16]uint32{}       // Working state

	// Convert block bytes to 32-bit words (little-endian)
	for i := 0; i < 16; i++ {
		m[i] = binary.LittleEndian.Uint32(block[i*4:])
	}

	// Initialize working state v[0..15]
	for i := 0; i < 8; i++ {
		v[i] = S.h[i] // First half from current state
	}

	// Second half from IV, with counter and flags mixed in
	v[8] = blake2sIV[0]
	v[9] = blake2sIV[1]
	v[10] = blake2sIV[2]
	v[11] = blake2sIV[3]
	v[12] = S.t[0] ^ blake2sIV[4] // Mix in counter low bits
	v[13] = S.t[1] ^ blake2sIV[5] // Mix in counter high bits
	v[14] = S.f[0] ^ blake2sIV[6] // Mix in finalization flag 0
	v[15] = S.f[1] ^ blake2sIV[7] // Mix in finalization flag 1

	// Perform 10 mixing rounds
	for r := 0; r < 10; r++ {
		blake2sRound(r, &v, m)
	}

	// Update the state with the result
	for i := 0; i < 8; i++ {
		S.h[i] = S.h[i] ^ v[i] ^ v[i+8]
	}
}

// blake2sIncrementCounter adds the specified number of bytes to the counter
// Parameters:
//   - S: Pointer to the state
//   - inc: Number of bytes to add to the counter
func blake2sIncrementCounter(S *blake2sState, inc uint32) {
	S.t[0] += inc
	// Handle carry to high word if necessary
	if S.t[0] < inc {
		S.t[1]++
	}
}

// blake2sIsLastBlock checks if the final block flag is set
// Parameter:
//   - S: Pointer to the state
//
// Returns true if the finalization flag is set
func blake2sIsLastBlock(S *blake2sState) bool {
	return S.f[0] != 0
}

// blake2sSetLastBlock sets the finalization flags
// This indicates that the current block is the last one
// Parameter:
//   - S: Pointer to the state
func blake2sSetLastBlock(S *blake2sState) {
	// Set last node flag if necessary
	if S.lastNode != 0 {
		S.f[1] = 0xFFFFFFFF
	}
	// Always set last block flag
	S.f[0] = 0xFFFFFFFF
}

// blake2sUpdate processes input data in chunks
// It feeds data into the hash state block by block
// Parameters:
//   - S: Pointer to the state
//   - in: Input data to process
//
// Returns 0 on success
func blake2sUpdate(S *blake2sState, in []byte) int {
	inlen := len(in)

	if inlen > 0 {
		left := S.buflen                          // Bytes already in buffer
		fill := uint32(BLAKE2S_BLOCKBYTES) - left // Space left in buffer

		// Handle case where input fills the buffer
		if uint32(inlen) > fill {
			S.buflen = 0 // Buffer will be emptied

			// Fill buffer with start of input
			copy(S.buf[left:], in[:fill])

			// Process the now full buffer
			blake2sIncrementCounter(S, BLAKE2S_BLOCKBYTES)
			blake2sCompress(S, S.buf[:])

			// Move past the processed data
			in = in[fill:]
			inlen -= int(fill)

			// Process full blocks directly from input
			for inlen > BLAKE2S_BLOCKBYTES {
				blake2sIncrementCounter(S, BLAKE2S_BLOCKBYTES)
				blake2sCompress(S, in[:BLAKE2S_BLOCKBYTES])
				in = in[BLAKE2S_BLOCKBYTES:]
				inlen -= BLAKE2S_BLOCKBYTES
			}
		}

		// Store any remaining input in the buffer
		copy(S.buf[S.buflen:], in[:inlen])
		S.buflen += uint32(inlen)
	}

	return 0 // Success
}

// blake2sFinal finalizes the hash computation
// It processes any remaining data, sets finalization flags, and outputs the digest
// Parameters:
//   - S: Pointer to the state
//   - out: Output buffer for the digest
//   - outlen: Size of the desired output in bytes
//
// Returns 0 on success, negative on error
func blake2sFinal(S *blake2sState, out []byte, outlen int) int {
	// Check for valid output parameters
	if out == nil || outlen < int(S.outlen) || len(out) < outlen {
		return -1
	}

	// Check if already finalized
	if blake2sIsLastBlock(S) {
		return -1
	}

	// Count remaining bytes in buffer
	blake2sIncrementCounter(S, S.buflen)

	// Set finalization flags
	blake2sSetLastBlock(S)

	// Pad buffer with zeros
	for i := S.buflen; i < BLAKE2S_BLOCKBYTES; i++ {
		S.buf[i] = 0
	}

	// Process the final block
	blake2sCompress(S, S.buf[:])

	// Ensure we don't write past end of output buffer
	hashSize := int(S.outlen)
	if outlen < hashSize {
		hashSize = outlen
	}

	// Write digest to output buffer, little-endian format
	// Only write up to min(outlen, S.outlen) bytes
	for i := 0; i < (hashSize+3)/4; i++ {
		// Make sure we don't write past the end of the output buffer
		if i*4 < hashSize {
			remaining := hashSize - i*4
			if remaining >= 4 {
				binary.LittleEndian.PutUint32(out[i*4:], S.h[i])
			} else {
				// Partial word at the end
				var wordBytes [4]byte
				binary.LittleEndian.PutUint32(wordBytes[:], S.h[i])
				copy(out[i*4:], wordBytes[:remaining])
			}
		}
	}

	return 0 // Success
}

// blake2sInit initializes a BLAKE2s state with the desired output size
// Parameters:
//   - S: Pointer to the state to initialize
//   - outlen: Desired digest length in bytes (1-32)
//
// Returns 0 on success, negative on error
func blake2sInit(S *blake2sState, outlen uint8) int {
	// Check for valid output length
	if outlen == 0 || outlen > BLAKE2S_OUTBYTES {
		return -1
	}

	// Initialize with default values
	blake2sInit0(S)

	// Set desired output length and customize h[0]
	S.outlen = outlen
	S.h[0] ^= uint32(outlen) | (1 << 16) | (1 << 24)

	return 0 // Success
}

// blake2s performs the complete BLAKE2s hash computation
// This is a convenience function that handles initialization, data processing, and finalization
// Parameters:
//   - msg: Input message to hash
//   - out: Output buffer for the digest
//   - outlen: Size of the desired output in bytes (1-32)
//
// Returns 0 on success, negative on error
func blake2s(msg []byte, out []byte, outlen int) int {
	// Check for valid output parameters
	if outlen <= 0 || outlen > BLAKE2S_OUTBYTES || out == nil || len(out) < outlen {
		return -1
	}

	var S blake2sState

	// Initialize state
	if blake2sInit(&S, uint8(outlen)) < 0 {
		return -1
	}

	// Process message
	if blake2sUpdate(&S, msg) < 0 {
		return -1
	}

	// Finalize and get digest
	if blake2sFinal(&S, out, outlen) < 0 {
		return -1
	}

	return 0 // Success
}

// blake2s_160A computes a BLAKE2s-160 hash of the input data
// This is a convenience wrapper for BLAKE2s with a 20-byte output
// Parameters:
//   - data: Input data buffer
//   - data_len: Length of data to hash
//   - digest: Output buffer for the hash result (must be at least 20 bytes)
//
// Returns 0 on success, negative on error
func blake2s_160A(data []byte, data_len int, digest []byte) int {
	// Validate input
	if data_len < 0 || data_len > len(data) || len(digest) < 20 {
		return -1
	}
	return blake2s(data[:data_len], digest, 20)
}

// blake2s_256A computes a BLAKE2s-256 hash of the input data
// This is a convenience wrapper for BLAKE2s with a 32-byte output
// Parameters:
//   - data: Input data buffer
//   - data_len: Length of data to hash
//   - digest: Output buffer for the hash result (must be at least 32 bytes)
//
// Returns 0 on success, negative on error
func blake2s_256A(data []byte, data_len int, digest []byte) int {
	// Validate input
	if data_len < 0 || data_len > len(data) || len(digest) < 32 {
		return -1
	}
	return blake2s(data[:data_len], digest, 32)
}
