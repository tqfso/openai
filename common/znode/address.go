package znode

// GenAddress generates an address from a public key
// The address format is:
//   - 1 byte prefix (81 = 'Z' in base58)
//   - 20 bytes of BLAKE2s-160 hash of the BLAKE2b-512 hash of the public key
//   - 4 bytes checksum (first 4 bytes of BLAKE2s-256 hash of the above 21 bytes)
//
// Total address size is 25 bytes
// Parameter:
//   - publicKey: The input public key (33 bytes)
//   - address: Output buffer for the generated address (must be at least 25 bytes)
func GenAddress(publicKey []byte, address []byte) {
	// Validate input parameters
	if len(publicKey) < SIZE_Publickey || len(address) < SIZE_Address {
		return
	}

	s1 := make([]byte, 64) // Buffer for BLAKE2b-512 output
	s := make([]byte, 64)  // Buffer for BLAKE2s-256 output

	// Calculate 160-bit hash of the public key
	// Step 1: Generate 64-byte hash from 33-byte public key
	blake2b_512A(publicKey, SIZE_Publickey, s1)

	// Step 2: Generate 20-byte address from 64-byte hash
	blake2s_160A(s1, 64, address[1:])

	// Step 3: Add identifier => base58(81) = 'Z' (ZOL identifier)
	address[0] = 81

	// Calculate checksum
	// Step 4: Hash the first 21 bytes (prefix + 20-byte address)
	blake2s_256A(address[:21], 21, s)

	// Step 5: Add 4-byte checksum
	copy(address[21:], s[:4]) // Address = unbase58('Z') + 20-byte address code + 4-byte checksum
}

// NodeIDToAddress converts a node ID to an address
// The address format is:
//   - 1 byte prefix (81 = 'Z' in base58)
//   - 20 bytes of the node ID
//   - 4 bytes checksum (first 4 bytes of BLAKE2s-256 hash of the above 21 bytes)
//
// Total address size is 25 bytes
// Parameters:
//   - nodeID: The input node ID (20 bytes)
//   - address: Output buffer for the generated address (must be at least 25 bytes)
func NodeIDToAddress(nodeID []byte, address []byte) {
	// Validate input parameters
	if len(nodeID) < SIZE_NodeID || len(address) < SIZE_Address {
		return
	}

	s := make([]byte, 64) // Buffer for BLAKE2s-256 output

	// Step 1: Copy nodeID to address (after the prefix byte)
	copy(address[1:], nodeID[:SIZE_NodeID])

	// Step 2: Add identifier => base58(81) = 'Z' (ZOL identifier)
	address[0] = 81

	// Calculate checksum
	// Step 3: Hash the first 21 bytes (prefix + 20-byte node ID)
	blake2s_256A(address[:21], 21, s)

	// Step 4: Add 4-byte checksum
	copy(address[21:], s[:4]) // Address = unbase58('Z') + 20-byte node ID + 4-byte checksum
}
