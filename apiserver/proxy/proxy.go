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
	"mime/multipart"
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
	contentType := c.ContentType()

	// 处理 form-data 格式
	if strings.Contains(contentType, "multipart/form-data") {
		// 解析 multipart form，32MB 限制
		err := c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			return NewResponseError(http.StatusBadRequest, "Failed to parse form data")
		}

		// 获取model字段
		modelName := c.Request.FormValue("model")
		if modelName == "" {
			return NewResponseError(http.StatusBadRequest, "Model name is required")
		}
		h.ModelName = modelName
		h.RequestBody = make(map[string]any)
		h.RequestBody["model"] = modelName

		// 获取其他字段		
		if form := c.Request.MultipartForm; form != nil {
			for key, values := range form.Value {
				if len(values) > 0 && key != "model" {
					h.RequestBody[key] = values[0]
				}
			}
			h.RequestBody["_files"] = form.File
		}
		return nil
	}

	// 处理 JSON 格式
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
	contentType := c.ContentType()

	if strings.Contains(contentType, "multipart/form-data") {
		// 重新构建 multipart 请求
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 写入普通字段
		for key, value := range h.RequestBody {
			if key != "_files" && value != nil {
				writer.WriteField(key, fmt.Sprintf("%v", value))
			}
		}

		// 写入文件
		if files, ok := h.RequestBody["_files"].(map[string][]*multipart.FileHeader); ok {
			for key, fileHeaders := range files {
				for _, fileHeader := range fileHeaders {
					file, err := fileHeader.Open()
					if err != nil {
						continue
					}
					part, err := writer.CreateFormFile(key, fileHeader.Filename)
					if err != nil {
						file.Close()
						continue
					}
					io.Copy(part, file)
					file.Close()
				}
			}
		}

		writer.Close()
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())
		c.Request.ContentLength = int64(body.Len())
		c.Request.Body = io.NopCloser(body)
		return
	}

	// 对于 JSON 格式，重新设置请求体
	data, _ := json.Marshal(h.RequestBody)
	c.Request.Header.Set("Content-Type", "application/json")
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
