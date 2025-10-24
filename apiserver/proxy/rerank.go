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

type RerankHandler struct {
	Handler
}

func NewRerankHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &RerankHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *RerankHandler) OnAfter(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewResponseError(http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %v", err))
	}
	resp.Body.Close()

	var completionResponse ChatCompletionResponse
	if err := json.Unmarshal(data, &completionResponse); err == nil {
		h.HandleUsage(completionResponse.Usage)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}

func (h *RerankHandler) HandleUsage(usage *ChatCompletionUsage) {
	if usage == nil {
		return
	}

	logger.Info("Rerank Usage", logger.String("Model", h.ModelName), logger.Int("TotalTokens", usage.TotalTokens))
}
