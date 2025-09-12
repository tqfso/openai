package logger

import (
	"bufio"
	"io"
	"strings"
)

// GetWriter returns an io.Writer that redirects Gin's logs to zap logger.
// Usage: gin.DefaultWriter = logger.GetWriter()
func GetWriter() io.Writer {

	reader, writer := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.TrimSpace(text)

			if text == "" {
				continue
			}

			if strings.Contains(text, "[GIN-debug]") {
				logger.Debug(text)
			} else {
				logger.Info(text)
			}
		}
	}()

	return writer
}
