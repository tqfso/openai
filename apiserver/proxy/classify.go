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

type ClassifyHandler struct {
	Handler
}

func NewClassifyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ClassifyHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ClassifyHandler) OnAfter(resp *http.Response) error {
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

func (h *ClassifyHandler) HandleUsage(usage *ChatCompletionUsage) {
	if usage == nil {
		return
	}

	logger.Info("Classify Usage", logger.String("Model", h.ModelName), logger.Int("PromptTokens", usage.PromptTokens), logger.Int("TotalTokens", usage.TotalTokens))
}
