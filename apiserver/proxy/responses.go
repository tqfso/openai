package proxy

import (
	"bytes"
	"common/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type OpenAIResponse struct {
	Usage *ResponseUsage `json:"usage,omitempty"`
}

type ResponsesHandler struct {
	Handler
}

func NewResponsesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ResponsesHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ResponsesHandler) OnBefore() error {
	_, exist := h.RequestBody["stream"]
	if exist {
		return NewResponseError(http.StatusNotImplemented, "Streaming is not supported for responses API without Harmony")
	}
	return nil
}

func (h *ResponsesHandler) OnAfter(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewResponseError(http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %v", err))
	}
	resp.Body.Close()

	// 解析使用量信息但不修改原始数据
	var openaiResponse OpenAIResponse
	if err := json.Unmarshal(data, &openaiResponse); err == nil {
		h.HandleUsage(openaiResponse.Usage)
	}

	// 使用原始数据重新设置响应体
	resp.Body = io.NopCloser(bytes.NewBuffer(data))

	return nil
}

func (h *ResponsesHandler) HandleUsage(usage *ResponseUsage) {
	if usage == nil {
		return
	}

	logger.Info("Usage",
		logger.String("Model", h.ModelName),
		logger.Int("InputTokens", usage.InputTokens),
		logger.Int("OutputTokens", usage.OutputTokens),
		logger.Int("TotalTokens", usage.TotalTokens))
}
