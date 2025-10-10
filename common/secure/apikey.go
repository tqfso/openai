package secure

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateApiKey() (string, error) {
	randomBytes := make([]byte, 24) // 24 bytes → base64URL 编码后约 32 字符
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	keyPart := base64.RawURLEncoding.EncodeToString(randomBytes)
	// 确保去掉 padding 且长度接近 32 字符
	keyPart = strings.TrimRight(keyPart, "=")
	if len(keyPart) > 32 {
		keyPart = keyPart[:32]
	}

	return "sk-" + keyPart, nil
}
