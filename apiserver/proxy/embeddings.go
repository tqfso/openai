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

type EmbeddingsHandler struct {
	Handler
}

func NewEmbeddingsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &EmbeddingsHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *EmbeddingsHandler) OnAfter(resp *http.Response) error {
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

func (h *EmbeddingsHandler) HandleUsage(usage *ChatCompletionUsage) {
	if usage == nil {
		return
	}

	logger.Info("Usage",
		logger.String("Model", h.ModelName),
		logger.Int("PromptTokens", usage.PromptTokens),
		logger.Int("CompletionTokens", usage.CompletionTokens),
		logger.Int("TotalTokens", usage.TotalTokens))
}
