// Package llm
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: errors.go
/@Description: LLM error handling
/*/

package llm

import (
	"fmt"
	"net/http"
)

// Metrics 监控指标
type Metrics struct {
	RequestCount    int64 `json:"request_count"`
	RequestDuration int64 `json:"request_duration_ms"`
	ErrorCount      int64 `json:"error_count"`
	TokenCount      int64 `json:"token_count"`
}

// IsRetryableError 判断错误是否可重试
func IsRetryableError(err error) bool {
	if llmErr, ok := err.(*LLMError); ok {
		switch llmErr.Type {
		case string(ErrorTypeRateLimit), string(ErrorTypeNetwork), string(ErrorTypeServer):
			return true
		}
	}
	return false
}

// GetHTTPStatusFromError 从错误获取HTTP状态码
func GetHTTPStatusFromError(err error) int {
	if llmErr, ok := err.(*LLMError); ok {
		switch llmErr.Type {
		case string(ErrorTypeInvalidRequest):
			return http.StatusBadRequest
		case string(ErrorTypeAuth):
			return http.StatusUnauthorized
		case string(ErrorTypeRateLimit):
			return http.StatusTooManyRequests
		case string(ErrorTypeServer):
			return http.StatusInternalServerError
		case string(ErrorTypeNetwork):
			return http.StatusServiceUnavailable
		}
	}
	return http.StatusInternalServerError
}

// WrapError 包装错误
func WrapError(err error, message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeServer),
		Message: fmt.Sprintf("%s: %v", message, err),
	}
}

// NewInvalidRequestError 创建无效请求错误
func NewInvalidRequestError(message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeInvalidRequest),
		Message: message,
	}
}

// NewAuthError 创建认证错误
func NewAuthError(message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeAuth),
		Message: message,
	}
}

// NewRateLimitError 创建限流错误
func NewRateLimitError(message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeRateLimit),
		Message: message,
	}
}

// NewServerError 创建服务器错误
func NewServerError(message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeServer),
		Message: message,
	}
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string) *LLMError {
	return &LLMError{
		Type:    string(ErrorTypeNetwork),
		Message: message,
	}
}