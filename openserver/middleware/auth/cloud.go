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

const (
	ClientH5    byte = 1 // H5
	ClientThird byte = 2 // 第三方
	ClientCloud byte = 3 // 零极云
)

const (
	NormalToken  byte = 1
	RefreshToken byte = 2
)

type ZCloudToken struct {
	Version     byte   // 版本
	Client      byte   // 客户端
	Type        byte   // 类型
	ExpiredTime int64  // 过期时间
	DmappId     []byte // 应用ID
	Sign        []byte // 签名
}

func (t ZCloudToken) Check() error {
	expiredTime := time.Unix(t.ExpiredTime, 0)
	if time.Now().Before(expiredTime) {
		return fmt.Errorf("token is expired")
	}

	if t.Client != ClientCloud {
		return fmt.Errorf("token is not for cloud")
	}

	return nil
}

func (t *ZCloudToken) Decode(token string) error {

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

	t.Client = tokenBin[offset]
	offset += 1

	t.Type = tokenBin[offset]
	offset += 1

	t.DmappId = tokenBin[offset : offset+20]
	offset += 20

	expireStrEnd := offset + 10
	expireStr := string(tokenBin[offset:expireStrEnd])
	t.ExpiredTime, _ = strconv.ParseInt(expireStr, 10, 64)
	offset = expireStrEnd

	if len(tokenBin[offset:]) != sha256.Size {
		return fmt.Errorf("sign length invalid, expected %d, got %d", sha256.Size, len(tokenBin[offset:]))
	}
	t.Sign = tokenBin[offset : offset+sha256.Size]

	return nil
}

// 签名
func ZDanSign(data, key string) [sha256.Size]byte {
	data = data + key
	return sha256.Sum256([]byte(data))
}

// 生成零极云Token
func ZCloudMakeToken(dmappHexId, appKey string) (token string, err error) {

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

// 验证零极云TOKEN
func ZCloudVerifyToken(token, appKey string) (*ZCloudToken, error) {

	data := &ZCloudToken{}
	if err := data.Decode(token); err != nil {
		return nil, err
	}

	if err := data.Check(); err != nil {
		return nil, err
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

func ZCloudAuthHander() gin.HandlerFunc {
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
