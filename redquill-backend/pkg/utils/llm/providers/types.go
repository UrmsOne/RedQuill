// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: types.go
/@Description: Provider types and interfaces
/*/

package providers

import (
	"context"
	"net/http"
)

// LLMConfig LLM配置
type LLMConfig struct {
	Provider   string            `json:"provider"` // openai, azure, ollama, deepseek, doubao, qwen
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key"`
	Model      string            `json:"model"`
	Headers    map[string]string `json:"headers,omitempty"`
	Timeout    int               `json:"timeout"`
	MaxRetries int               `json:"max_retries"`
	RetryDelay int               `json:"retry_delay"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream,omitempty"`
	Temperature      float64   `json:"temperature,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
}

// Message 消息类型
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Created int64    `json:"created"`
}

// Choice 选择项
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Delta        Message `json:"delta,omitempty"` // 流式响应使用
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// StreamChunk 流式响应块
type StreamChunk struct {
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Choices []Choice  `json:"choices"`
	Usage   *Usage    `json:"usage,omitempty"`
	Error   *LLMError `json:"error,omitempty"`
}

// Model 模型信息
type Model struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner,omitempty"`
	Context int    `json:"context,omitempty"`
}

// LLMError LLM错误
type LLMError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func (e *LLMError) Error() string {
	return e.Message
}

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeInvalidRequest ErrorType = "invalid_request"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
	ErrorTypeAuth           ErrorType = "auth"
	ErrorTypeServer         ErrorType = "server"
	ErrorTypeNetwork        ErrorType = "network"
)

// Provider 提供商接口
type Provider interface {
	// Chat 同步响应
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ChatStream 流式响应
	ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)

	// Health 健康检查
	Health(ctx context.Context) error

	// Models 获取支持的模型列表
	Models(ctx context.Context) ([]Model, error)
}

// StreamProcessor 流式处理器
type StreamProcessor struct {
	client *http.Client
}

// NewStreamProcessor 创建流式处理器
func NewStreamProcessor(client *http.Client) *StreamProcessor {
	return &StreamProcessor{
		client: client,
	}
}
