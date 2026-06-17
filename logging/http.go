package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxBodyLogBytes = 4096

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)

		requestBody := readRequestBody(c)
		start := time.Now()

		writer := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = writer

		logger := Default().With(
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
		)
		ctx := WithContext(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		logger.InfoContext(ctx, "http request started",
			slog.String("raw_path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Any("request_body", payloadForLog(requestBody)),
		)

		c.Next()

		latency := time.Since(start)
		responseBody := writer.body.Bytes()
		attrs := []any{
			slog.Int("status", c.Writer.Status()),
			slog.Int("response_size", writer.body.Len()),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.Any("response_body", payloadForLog(responseBody)),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}
		if c.Writer.Status() >= http.StatusInternalServerError {
			logger.ErrorContext(ctx, "http request completed", attrs...)
			return
		}
		if c.Writer.Status() >= http.StatusBadRequest {
			logger.WarnContext(ctx, "http request completed", attrs...)
			return
		}
		logger.InfoContext(ctx, "http request completed", attrs...)
	}
}

func readRequestBody(c *gin.Context) []byte {
	if c.Request == nil || c.Request.Body == nil {
		return nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return []byte(`"failed to read request body"`)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}

func payloadForLog(body []byte) any {
	if len(body) == 0 {
		return nil
	}
	if len(body) > maxBodyLogBytes {
		return map[string]any{
			"truncated": true,
			"size":      len(body),
			"preview":   string(body[:maxBodyLogBytes]),
		}
	}
	var parsed any
	if err := json.Unmarshal(body, &parsed); err == nil {
		return parsed
	}
	return string(body)
}
