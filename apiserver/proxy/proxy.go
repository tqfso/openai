package proxy

import (
	"apiserver/model"
	"apiserver/user"
	"bytes"
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
	OnBefore() error                   // 前置处理
	OnAfter(resp *http.Response) error // 后置处理
}

type Handler struct {
	GinContext *gin.Context
	Task       TaskInterface
	ModelName  string
	ApiKey     string
	ApiKeyInfo *user.ApiKeyInfo
	TargetURL  *url.URL
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
			if err == nil {
				return
			}
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
	var reqBody map[string]interface{}
	if err := json.Unmarshal(data, &reqBody); err != nil {
		return NewResponseError(http.StatusBadRequest, "Invalid request body")
	}

	// 从请求体中获取模型名称
	if modelName, ok := reqBody["model"].(string); ok {
		h.ModelName = modelName
	} else {
		return NewResponseError(http.StatusBadRequest, "Model name is required")
	}

	// 重新设置请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	return nil
}

func (h *Handler) selectTarget() *ResponseError {
	target := model.SelectTarget(h.ModelName)
	if target == nil {
		return NewResponseError(http.StatusBadGateway, "Not found target")
	}

	h.TargetURL, _ = url.Parse(fmt.Sprintf("http://%s:%d", target.IP, target.Port))

	return nil
}
