// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: openai.go
/@Description: OpenAI provider implementation
/*/

package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIProvider OpenAI提供商
type OpenAIProvider struct {
	config LLMConfig
	client *http.Client
	stream *StreamProcessor
}

// NewOpenAIProvider 创建OpenAI提供商
func NewOpenAIProvider(config LLMConfig, client *http.Client) *OpenAIProvider {
	return &OpenAIProvider{
		config: config,
		client: client,
		stream: NewStreamProcessor(client),
	}
}

// Chat 同步聊天
func (p *OpenAIProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("create request error: %v", err),
		}
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	
	// 添加自定义头
	for k, v := range p.config.Headers {
		httpReq.Header.Set(k, v)
	}
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("request error: %v", err),
		}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("read response error: %v", err),
		}
	}
	
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error LLMError `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, &errorResp.Error
		}
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
		}
	}
	
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("unmarshal response error: %v", err),
		}
	}
	
	return &chatResp, nil
}

// ChatStream 流式聊天
func (p *OpenAIProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	req.Stream = true
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("create request error: %v", err),
		}
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")
	
	// 添加自定义头
	for k, v := range p.config.Headers {
		httpReq.Header.Set(k, v)
	}
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("request error: %v", err),
		}
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
		}
	}
	
	return p.stream.ProcessSSEStream(ctx, resp)
}

// Health 健康检查
func (p *OpenAIProvider) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/models", nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: HTTP %d", resp.StatusCode)
	}
	
	return nil
}

// Models 获取模型列表
func (p *OpenAIProvider) Models(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/models", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get models failed: HTTP %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var modelsResp struct {
		Data []Model `json:"data"`
	}
	
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, err
	}
	
	return modelsResp.Data, nil
}