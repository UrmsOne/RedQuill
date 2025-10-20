// Package llm
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: client.go
/@Description: LLM client factory and main client implementation
/*/

package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"redquill-backend/pkg/utils/llm/providers"
)

// Client LLM客户端实现
type Client struct {
	provider providers.Provider
}

// NewClient 创建新的LLM客户端
func NewClient(config LLMConfig) (*Client, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	var provider providers.Provider

	switch config.Provider {
	case "openai":
		provider = providers.NewOpenAIProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "azure":
		provider = providers.NewAzureProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "ollama":
		provider = providers.NewOllamaProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "deepseek":
		provider = providers.NewDeepSeekProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "doubao":
		provider = providers.NewDoubaoProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "qwen":
		provider = providers.NewQwenProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	case "wenxin":
		provider = providers.NewWenxinProvider(providers.LLMConfig{
			Provider:   config.Provider,
			BaseURL:    config.BaseURL,
			APIKey:     config.APIKey,
			Model:      config.Model,
			Headers:    config.Headers,
			Timeout:    int(config.Timeout.Seconds()),
			MaxRetries: config.MaxRetries,
			RetryDelay: int(config.RetryDelay.Seconds()),
		}, client)
	default:
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("unsupported provider: %s", config.Provider),
		}
	}

	return &Client{
		provider: provider,
	}, nil
}

// NewClientWithConfig 使用多厂商配置创建客户端
func NewClientWithConfig(config *MultiProviderConfig, providerName string) (*Client, error) {
	if providerName == "" {
		providerName = config.Default
	}

	providerConfig, exists := config.Providers[providerName]
	if !exists {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("provider not found: %s", providerName),
		}
	}

	return NewClient(providerConfig)
}

// Chat 同步聊天
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// 转换请求格式
	providerReq := providers.ChatRequest{
		Model:            req.Model,
		Messages:         convertMessages(req.Messages),
		Stream:           req.Stream,
		Temperature:      req.Temperature,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		FrequencyPenalty: req.FrequencyPenalty,
		PresencePenalty:  req.PresencePenalty,
	}

	resp, err := c.provider.Chat(ctx, providerReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	return &ChatResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Choices: convertChoices(resp.Choices),
		Usage:   convertUsage(resp.Usage),
		Created: resp.Created,
	}, nil
}

// ChatStream 流式聊天
func (c *Client) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	// 转换请求格式
	providerReq := providers.ChatRequest{
		Model:            req.Model,
		Messages:         convertMessages(req.Messages),
		Stream:           req.Stream,
		Temperature:      req.Temperature,
		MaxTokens:        req.MaxTokens,
		TopP:             req.TopP,
		FrequencyPenalty: req.FrequencyPenalty,
		PresencePenalty:  req.PresencePenalty,
	}

	stream, err := c.provider.ChatStream(ctx, providerReq)
	if err != nil {
		return nil, err
	}

	// 转换流式响应
	result := make(chan StreamChunk, 100)
	go func() {
		defer close(result)
		for chunk := range stream {
			select {
			case result <- StreamChunk{
				ID:      chunk.ID,
				Model:   chunk.Model,
				Choices: convertChoices(chunk.Choices),
				Usage:   convertUsagePtr(chunk.Usage),
				Error:   convertError(chunk.Error),
			}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return result, nil
}

// Health 健康检查
func (c *Client) Health(ctx context.Context) error {
	return c.provider.Health(ctx)
}

// Models 获取模型列表
func (c *Client) Models(ctx context.Context) ([]Model, error) {
	models, err := c.provider.Models(ctx)
	if err != nil {
		return nil, err
	}

	// 转换模型格式
	result := make([]Model, len(models))
	for i, model := range models {
		result[i] = Model{
			ID:      model.ID,
			Name:    model.Name,
			Owner:   model.Owner,
			Context: model.Context,
		}
	}

	return result, nil
}

// NewClientFromEnv 从环境变量创建客户端
func NewClientFromEnv() (*Client, error) {
	config := LoadConfigFromEnv()
	return NewClient(*config)
}

// NewClientFromFile 从配置文件创建客户端
func NewClientFromFile(filename string, providerName string) (*Client, error) {
	config, err := LoadConfigFromFile(filename)
	if err != nil {
		return nil, err
	}

	return NewClientWithConfig(config, providerName)
}

// 转换函数
func convertMessages(messages []Message) []providers.Message {
	result := make([]providers.Message, len(messages))
	for i, msg := range messages {
		result[i] = providers.Message{
			Role:    msg.Role,
			Content: msg.Content,
			Name:    msg.Name,
		}
	}
	return result
}

func convertChoices(choices []providers.Choice) []Choice {
	result := make([]Choice, len(choices))
	for i, choice := range choices {
		result[i] = Choice{
			Index:        choice.Index,
			Message:      convertMessage(choice.Message),
			FinishReason: choice.FinishReason,
			Delta:        convertMessage(choice.Delta),
		}
	}
	return result
}

func convertMessage(msg providers.Message) Message {
	return Message{
		Role:    msg.Role,
		Content: msg.Content,
		Name:    msg.Name,
	}
}

func convertUsage(usage providers.Usage) *Usage {
	return &Usage{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}
}

func convertUsagePtr(usage *providers.Usage) *Usage {
	if usage == nil {
		return nil
	}
	result := convertUsage(*usage)
	return result
}

func convertError(err *providers.LLMError) *LLMError {
	if err == nil {
		return nil
	}
	return &LLMError{
		Type:    err.Type,
		Message: err.Message,
		Code:    err.Code,
		Details: err.Details,
	}
}
