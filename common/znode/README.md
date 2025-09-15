# BLAKE2 Implementation

This codebase implements the BLAKE2 cryptographic hash function family, including BLAKE2s (optimized for 32-bit platforms) and BLAKE2b (optimized for 64-bit platforms), along with address generation utilities.

## File Structure

The code is divided into the following modules, each focusing on specific functionality:

1. **constants.go**
   - Contains constant definitions
   - Initialization vectors
   - Permutation tables

2. **utils.go**
   - Utility functions
   - Data structures

3. **blake2s.go**
   - Complete implementation of the BLAKE2s hash algorithm
   - 32-bit version, suitable for memory-constrained environments

4. **blake2b.go**
   - Complete implementation of the BLAKE2b hash algorithm
   - 64-bit version, providing higher security

5. **address.go**
   - Address generation functionality
   - Support for creating addresses from public keys and node IDs

## Main Algorithm Flow

### Hash Computation Flow
```
Initialize state → Process data (can be called multiple times) → Finalize and output hash
```

### Address Generation Flow
```
Input (public key or node ID) → Hash (using BLAKE2) → Add prefix → Calculate checksum → Final address
```

## Address Format

The generated addresses have a specific format:
- 1 byte prefix (81, which corresponds to 'Z' in base58 encoding)
- 20 bytes of content (hash of public key or node ID)
- 4 bytes checksum (hash of the previous 21 bytes)

Total address size is 25 bytes, which provides a good balance between security and usability.

## Usage Examples

```go
// Calculate a 160-bit hash using BLAKE2s
func calculateHash() {
    data := []byte("Hello, World!")
    digest := make([]byte, 20)
    
    blake2s_160A(data, len(data), digest)
    
    fmt.Printf("BLAKE2s-160 Hash: %x\n", digest)
}

// Generate an address from a public key
func createAddress() {
    publicKey := make([]byte, SIZE_Publickey) // 33 bytes
    address := make([]byte, SIZE_Address)     // 25 bytes
    
    // Fill public key data...
    
    GenAddress(publicKey, address)
    
    fmt.Printf("Generated Address: %x\n", address)
}
```

## Performance Considerations

- BLAKE2 is designed for high performance, with BLAKE2s optimized for 32-bit systems and BLAKE2b for 64-bit systems
- The implementation processes data in blocks, minimizing memory usage for large inputs
- Direct processing of full blocks from input avoids unnecessary copying when possible