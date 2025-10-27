package proxy

import (
	"apiserver/model"
	"apiserver/user"
	"bytes"
	"common/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var sharedTransport = &http.Transport{
	MaxIdleConns:        200,
	MaxIdleConnsPerHost: 50,
	IdleConnTimeout:     90 * time.Second,
}

type TaskInterface interface {
	OnBefore() error                   // 转发请求前处理
	OnAfter(resp *http.Response) error // 转发响应前处理
}

type Handler struct {
	GinContext  *gin.Context
	Task        TaskInterface
	RequestBody map[string]any
	ModelName   string
	ApiKey      string
	ApiKeyInfo  *user.ApiKeyInfo
	TargetURL   *url.URL
}

func NewDefaultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &Handler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *Handler) SetTaskHandler(handler TaskInterface) {
	h.Task = handler
}

func (h *Handler) OnBefore() error {
	return nil
}

func (h *Handler) OnAfter(resp *http.Response) error {
	return nil
}

func (h *Handler) OnRequest(c *gin.Context) {

	h.GinContext = c

	// 检查API密钥
	if err := h.checkApiKey(); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	// 检查模型名称
	if err := h.checkModelName(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	// 选择转发目标
	if err := h.selectTarget(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	// 转发前处理
	if err := h.Task.OnBefore(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	// 重新设置请求体
	h.setbackBody()

	proxy := &httputil.ReverseProxy{
		Transport: sharedTransport,
		Director: func(req *http.Request) {
			req.Host = h.TargetURL.Host
			req.URL.Scheme = h.TargetURL.Scheme
			req.URL.Host = h.TargetURL.Host
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = c.Request.URL.RawQuery
		},
		ModifyResponse: func(resp *http.Response) error {
			return h.Task.OnAfter(resp)
		},
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			logger.Error("ReverseProxy", logger.String("HOST", req.Host), logger.String("URI", req.RequestURI), logger.Err(err))
			rw.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(rw).Encode(NewResponseError(http.StatusBadGateway, err.Error()))
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) GetRequestContext() context.Context {
	return h.GinContext.Request.Context()
}

func (h *Handler) checkApiKey() *ResponseError {
	c := h.GinContext

	auth := c.GetHeader("Authorization")
	if auth == "" {
		return NewResponseError(http.StatusUnauthorized, "Authorization required")
	}

	if !strings.HasPrefix(auth, "Bearer") {
		return NewResponseError(http.StatusUnauthorized, "Bearer required")
	}

	h.ApiKey = auth[7:]
	if h.ApiKey == "" {
		return NewResponseError(http.StatusUnauthorized, "API KEY required")
	}

	var err error
	h.ApiKeyInfo, err = user.FindKey(h.GetRequestContext(), h.ApiKey)
	if err != nil {
		return NewResponseError(http.StatusUnauthorized, err.Error())
	}

	return nil
}

func (h *Handler) checkModelName() *ResponseError {
	c := h.GinContext

	// 获取原始请求数据
	data, err := c.GetRawData()
	if err != nil {
		return NewResponseError(http.StatusBadRequest, "Failed to read request body")
	}

	// 解析JSON
	if err := json.Unmarshal(data, &h.RequestBody); err != nil {
		return NewResponseError(http.StatusBadRequest, "Invalid request body")
	}

	// 从请求体中获取模型名称
	if modelName, ok := h.RequestBody["model"].(string); ok {
		h.ModelName = modelName
	} else {
		return NewResponseError(http.StatusBadRequest, "Model name is required")
	}

	return nil
}

// 设置转发请求体
func (h *Handler) setbackBody() {
	c := h.GinContext
	data, _ := json.Marshal(h.RequestBody)
	c.Request.ContentLength = int64(len(data))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
}

// 选择转发目标
func (h *Handler) selectTarget() *ResponseError {
	target := model.SelectTarget(h.ModelName)
	if target == nil {
		return NewResponseError(http.StatusBadGateway, "Not found target")
	}

	h.TargetURL, _ = url.Parse(fmt.Sprintf("http://%s:%d", target.IP, target.Port))

	return nil
}
