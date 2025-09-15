package auth

import (
	"common"
	"common/znode"
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

type ZUserToken struct {
	Version     byte   // 版本
	Client      byte   // 客户端
	Type        byte   // 类型
	ExpiredTime int64  // 过期时间
	DmappId     []byte // 应用ID
	PubKey      []byte // 用户公钥
	Sign        []byte // 签名
}

func (t ZUserToken) UserId() string {
	address := make([]byte, 25)
	znode.GenAddress(t.PubKey, address)
	return base58.Encode(address[:])
}

func (t ZUserToken) GetNodeId() string {
	nodeID := make([]byte, 20)
	address := make([]byte, 25)
	copy(nodeID, address[1:])
	return hex.EncodeToString(nodeID)
}

func (t ZUserToken) Check() error {
	expiredTime := time.Unix(t.ExpiredTime, 0)
	if time.Now().Before(expiredTime) {
		return fmt.Errorf("token is expired")
	}

	if t.Client != ClientH5 {
		return fmt.Errorf("token is not for user")
	}

	return nil
}

func (t *ZUserToken) Decode(token string) error {

	tokenBin := base58.Decode(token)
	if len(tokenBin) == 0 {
		return fmt.Errorf("invalid base58 token")
	}

	if len(tokenBin) < (66 + sha256.Size) {
		return fmt.Errorf("token too short")
	}

	t.Version = tokenBin[0]
	t.Client = tokenBin[1]
	t.Type = tokenBin[2]
	t.PubKey = tokenBin[3:36]
	t.DmappId = tokenBin[36:56]
	expiredTime := tokenBin[56:66]
	t.ExpiredTime, _ = strconv.ParseInt(string(expiredTime), 10, 64)

	t.Sign = tokenBin[66:]

	if len(t.Sign) != sha256.Size {
		return fmt.Errorf("sign length invalid, expected %d, got %d", sha256.Size, len(t.Sign))
	}

	return nil
}

func ZUserVerifyToken(token, appKey string) (*ZUserToken, error) {

	data := &ZUserToken{}
	if err := data.Decode(token); err != nil {
		return nil, err
	}

	if err := data.Check(); err != nil {
		return nil, err
	}

	signData := fmt.Sprintf("%d&%d&%d&%d&%s",
		data.Version,
		data.Client,
		data.Type,
		data.ExpiredTime,
		data.PubKey,
	)

	sign := ZDanSign(signData, appKey)
	if !hmac.Equal(data.Sign, sign[:]) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return data, nil
}

func ZUserAuthHander() gin.HandlerFunc {
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
