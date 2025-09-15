package middleware

import (
	"common"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
)

type TokenClient byte

const (
	ClientH5    TokenClient = 1 // H5
	ClientThird TokenClient = 2 // 第三方
	ClientCloud TokenClient = 3 // 零极云
)

const (
	NormalToken  byte = 1
	RefreshToken byte = 2
)

type ZDanToken struct {
	Version     byte        // 版本
	Type        byte        // 类型
	DmappId     string      // 应用ID
	UserId      string      // 用户ID
	Client      TokenClient // 客户端
	ExpiredTime int64       // 过期时间
	Sign        []byte      // 签名
}

func (t ZDanToken) IsExpired() bool {
	expiredTime := time.Unix(t.ExpiredTime, 0)
	return time.Now().After(expiredTime)
}

func (t *ZDanToken) Decode(token string) error {

	tokenBin := base58.Decode(token)
	if len(tokenBin) == 0 {
		return fmt.Errorf("invalid base58 token")
	}

	if len(tokenBin) < (1 + 1 + 1 + 20 + 10 + sha256.Size) {
		return fmt.Errorf("token too short")
	}

	offset := 0
	t.Version = tokenBin[offset]
	offset += 1

	t.Client = TokenClient(tokenBin[offset])
	offset += 1

	t.Type = tokenBin[offset]
	offset += 1

	dmappIDBytes := tokenBin[offset : offset+20]
	t.DmappId = hex.EncodeToString(dmappIDBytes)
	offset += 20

	expireStrEnd := offset + 10
	expireStr := string(tokenBin[offset:expireStrEnd])
	t.ExpiredTime, _ = strconv.ParseInt(expireStr, 10, 64)
	offset = expireStrEnd

	if len(tokenBin[offset:]) != sha256.Size {
		return fmt.Errorf("signature length invalid, expected %d, got %d", sha256.Size, len(tokenBin[offset:]))
	}
	t.Sign = tokenBin[offset : offset+sha256.Size]

	return nil
}

func ZDanAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("ZCookie")
		if token == "" {
			c.JSON(http.StatusUnauthorized, common.Response{Code: common.AuthError, Msg: "ZCookie required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 签名
func ZDanSign(data, key string) [sha256.Size]byte {
	data = data + key
	return sha256.Sum256([]byte(data))
}

// 生成Token
func ZDanMakeToken(dmappHexId, appKey string) (token string, err error) {

	dmappId, err := hex.DecodeString(dmappHexId)
	if err != nil {
		return
	}

	expireTime := time.Now().Add(60 * 24 * time.Minute)
	strExpire := strconv.FormatInt(expireTime.Unix(), 10)
	strTokenVersion := strconv.Itoa(1)
	strClientType := strconv.Itoa(3)
	strNormalToken := strconv.Itoa(1)

	var signData string
	signData += strTokenVersion + "&"
	signData += strClientType + "&"
	signData += strNormalToken + "&"
	signData += strExpire

	sign := ZDanSign(signData, appKey)

	var tokenBin []byte
	tokenBin = append(tokenBin, 1)                    // 版本
	tokenBin = append(tokenBin, byte(ClientCloud))    // 类型
	tokenBin = append(tokenBin, 1)                    // 1普通token, 2刷新token
	tokenBin = append(tokenBin, dmappId[:]...)        // DMAPP ID
	tokenBin = append(tokenBin, []byte(strExpire)...) // 过期时间
	tokenBin = append(tokenBin, sign[:]...)           // 签名
	token = base58.Encode(tokenBin)

	return
}

// 验证TOKEN
func ZDanVerifyToken(token, appKey string) (*ZDanToken, error) {

	data := &ZDanToken{}
	if err := data.Decode(token); err != nil {
		return nil, err
	}

	if data.IsExpired() {
		return nil, fmt.Errorf("token is expired")
	}

	signData := fmt.Sprintf("%d&%d&%d&%d",
		data.Version,
		data.Client,
		data.Type,
		data.ExpiredTime,
	)

	sign := ZDanSign(signData, appKey)
	if !hmac.Equal(data.Sign, sign[:]) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return data, nil
}
