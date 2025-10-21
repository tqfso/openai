package proxy

import (
	"bufio"
	"bytes"
	"common/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ChatCompletionUsage 定义了使用量的结构
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse 定义了响应的结构
type ChatCompletionResponse struct {
	Usage *ChatCompletionUsage `json:"usage,omitempty"`
}

// ChatCompletionChunk 定义了流式响应chunk的结构
type ChatCompletionChunk struct {
	Usage *ChatCompletionUsage `json:"usage,omitempty"`
}

type ChatCompletionsHandler struct {
	Handler
}

func NewChatCompletionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &ChatCompletionsHandler{}
		h.SetTaskHandler(h)
		h.OnRequest(c)
	}
}

func (h *ChatCompletionsHandler) OnAfter(resp *http.Response) error {

	isStream := strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream")

	if isStream {
		// 流式响应处理
		originalBody := resp.Body
		reader, writer := io.Pipe()
		resp.Body = reader

		go func() {
			defer writer.Close()
			defer originalBody.Close()

			var buf bytes.Buffer
			scanner := bufio.NewScanner(originalBody)
			var lastChunk ChatCompletionChunk

			for scanner.Scan() {
				line := scanner.Text()
				buf.WriteString(line)
				buf.WriteString("\n")

				if strings.HasPrefix(line, "data: ") {
					// 尝试解析并记录最后一个包含usage的chunk
					if strings.Contains(line, `"usage"`) {
						chunk := line[6:] // 去掉 "data: " 前缀
						if err := json.Unmarshal([]byte(chunk), &lastChunk); err == nil {
							if lastChunk.Usage != nil {
								logger.Info("Stream Usage",
									logger.String("Model", h.ModelName),
									logger.Int("Prompt", lastChunk.Usage.PromptTokens),
									logger.Int("Completion", lastChunk.Usage.CompletionTokens),
									logger.Int("TotalTokens", lastChunk.Usage.TotalTokens),
								)
							}
						}
					}
				}

				// 立即写入每一行数据
				if _, err := writer.Write(buf.Bytes()); err != nil {
					logger.Error("Failed to write stream data", logger.Err(err))
					return
				}
				buf.Reset()
			}

			if err := scanner.Err(); err != nil {
				logger.Error("Scanner error", logger.Err(err))
			}
		}()

	} else {
		// 非流式响应处理
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return NewResponseError(http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %v", err))
		}
		resp.Body.Close()

		// 解析使用量信息但不修改原始数据
		var completionResponse ChatCompletionResponse
		if err := json.Unmarshal(data, &completionResponse); err == nil {
			if completionResponse.Usage != nil {
				logger.Info("Usage",
					logger.String("Model", h.ModelName),
					logger.Int("Prompt", completionResponse.Usage.PromptTokens),
					logger.Int("Completion", completionResponse.Usage.CompletionTokens),
					logger.Int("TotalTokens", completionResponse.Usage.TotalTokens),
				)
			}
		}

		// 使用原始数据重新设置响应体
		resp.Body = io.NopCloser(bytes.NewBuffer(data))
	}

	return nil
}
