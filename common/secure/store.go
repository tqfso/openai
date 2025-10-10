package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 硬编码加密密钥（长度必须是 16, 24 或 32 字节）
var aesKey = []byte("Ztwv2hV14r2smkZ3WXYQFFj6sY6gjguu") // 32字节 = AES-256

// 加密 plaintext（base64 返回）
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	plainBytes := pkcs7Pad([]byte(plaintext), block.BlockSize())
	cipherText := make([]byte, len(plainBytes))

	// 使用零 IV（可换成随机 IV 存数据库）
	iv := make([]byte, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 解密 base64 密文
func Decrypt(cipherTextBase64 string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}
	if len(cipherText)%block.BlockSize() != 0 {
		return "", errors.New("invalid ciphertext size")
	}

	iv := make([]byte, block.BlockSize())
	mode := cipher.NewCBCDecrypter(block, iv)
	plainPadded := make([]byte, len(cipherText))
	mode.CryptBlocks(plainPadded, cipherText)

	plainBytes, err := pkcs7Unpad(plainPadded, block.BlockSize())
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// --- Padding helpers (PKCS7) ---

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padded data")
	}
	padLen := int(data[len(data)-1])
	if padLen > blockSize || padLen == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padLen; i < len(data); i++ {
		if data[i] != byte(padLen) {
			return nil, errors.New("invalid padding content")
		}
	}
	return data[:len(data)-padLen], nil
}
