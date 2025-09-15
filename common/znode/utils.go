package znode

// rotr32 performs a bitwise right rotation on a 32-bit word
// Parameters:
//   - w: The word to rotate
//   - c: Number of bits to rotate by
//
// Returns the rotated word
func rotr32(w uint32, c uint) uint32 {
	return (w >> c) | (w << (32 - c))
}

// rotr64 performs a bitwise right rotation on a 64-bit word
// Parameters:
//   - w: The word to rotate
//   - c: Number of bits to rotate by
//
// Returns the rotated word
func rotr64(w uint64, c uint) uint64 {
	return (w >> c) | (w << (64 - c))
}

// blake2sState holds the state for the BLAKE2s hashing algorithm
type blake2sState struct {
	h        [8]uint32                // Internal chain state
	t        [2]uint32                // Counter for bytes processed
	f        [2]uint32                // Finalization flags
	buf      [BLAKE2S_BLOCKBYTES]byte // Buffer for unprocessed data
	buflen   uint32                   // Amount of bytes in buffer
	outlen   uint8                    // Desired output size in bytes
	lastNode byte                     // Last node flag for tree hashing
}

// blake2bState holds the state for the BLAKE2b hashing algorithm
type blake2bState struct {
	h        [8]uint64                // Internal chain state
	t        [2]uint64                // Counter for bytes processed
	f        [2]uint64                // Finalization flags
	buf      [BLAKE2B_BLOCKBYTES]byte // Buffer for unprocessed data
	buflen   uint64                   // Amount of bytes in buffer
	outlen   uint64                   // Desired output size in bytes
	lastNode byte                     // Last node flag for tree hashing
}
