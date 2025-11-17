// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: qwen.go
/@Description: 千问 provider implementation
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

// QwenProvider 千问提供商
type QwenProvider struct {
	config LLMConfig
	client *http.Client
	stream *StreamProcessor
}

// NewQwenProvider 创建千问提供商
func NewQwenProvider(config LLMConfig, client *http.Client) *QwenProvider {
	return &QwenProvider{
		config: config,
		client: client,
		stream: NewStreamProcessor(client),
	}
}

// Chat 同步聊天
func (p *QwenProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	
	// 千问API需要转换请求格式
	qwenReq := p.convertToQwenRequest(req)
	reqBody, err := json.Marshal(qwenReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
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
	
	// 千问API响应格式转换
	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("unmarshal response error: %v", err),
		}
	}
	
	// 转换为统一格式
	chatResp := p.convertFromQwenResponse(qwenResp)
	return &chatResp, nil
}

// ChatStream 流式聊天
func (p *QwenProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	req.Stream = true
	
	// 千问API需要转换请求格式
	qwenReq := p.convertToQwenRequest(req)
	reqBody, err := json.Marshal(qwenReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/v1/chat/completions", bytes.NewBuffer(reqBody))
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
func (p *QwenProvider) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/v1/models", nil)
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
func (p *QwenProvider) Models(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/v1/models", nil)
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

// QwenRequest 千问请求格式
type QwenRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// QwenResponse 千问响应格式
type QwenResponse struct {
	ID      string        `json:"id"`
	Model   string        `json:"model"`
	Choices []QwenChoice  `json:"choices"`
	Usage   Usage         `json:"usage"`
	Created int64         `json:"created"`
}

// QwenChoice 千问选择项
type QwenChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// convertToQwenRequest 转换为千问请求格式
func (p *QwenProvider) convertToQwenRequest(req ChatRequest) QwenRequest {
	return QwenRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}
}

// convertFromQwenResponse 从千问响应转换为统一格式
func (p *QwenProvider) convertFromQwenResponse(resp QwenResponse) ChatResponse {
	choices := make([]Choice, len(resp.Choices))
	for i, choice := range resp.Choices {
		choices[i] = Choice{
			Index:        choice.Index,
			Message:      choice.Message,
			FinishReason: choice.FinishReason,
		}
	}
	
	return ChatResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Choices: choices,
		Usage:   resp.Usage,
		Created: resp.Created,
	}
}
