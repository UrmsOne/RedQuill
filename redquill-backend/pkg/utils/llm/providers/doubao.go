// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: doubao.go
/@Description: 豆包 provider implementation
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

// DoubaoProvider 豆包提供商
type DoubaoProvider struct {
	config LLMConfig
	client *http.Client
	stream *StreamProcessor
}

// NewDoubaoProvider 创建豆包提供商
func NewDoubaoProvider(config LLMConfig, client *http.Client) *DoubaoProvider {
	return &DoubaoProvider{
		config: config,
		client: client,
		stream: NewStreamProcessor(client),
	}
}

// Chat 同步聊天
func (p *DoubaoProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	
	// 豆包API需要转换请求格式
	doubaoReq := p.convertToDoubaoRequest(req)
	reqBody, err := json.Marshal(doubaoReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/api/v1/chat/completions", bytes.NewBuffer(reqBody))
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
	
	// 豆包API响应格式转换
	var doubaoResp DoubaoResponse
	if err := json.Unmarshal(body, &doubaoResp); err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("unmarshal response error: %v", err),
		}
	}
	
	// 转换为统一格式
	chatResp := p.convertFromDoubaoResponse(doubaoResp)
	return &chatResp, nil
}

// ChatStream 流式聊天
func (p *DoubaoProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	req.Stream = true
	
	// 豆包API需要转换请求格式
	doubaoReq := p.convertToDoubaoRequest(req)
	reqBody, err := json.Marshal(doubaoReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/api/v1/chat/completions", bytes.NewBuffer(reqBody))
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
func (p *DoubaoProvider) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/api/v1/models", nil)
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
func (p *DoubaoProvider) Models(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/api/v1/models", nil)
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

// DoubaoRequest 豆包请求格式
type DoubaoRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// DoubaoResponse 豆包响应格式
type DoubaoResponse struct {
	ID      string           `json:"id"`
	Model   string           `json:"model"`
	Choices []DoubaoChoice   `json:"choices"`
	Usage   Usage            `json:"usage"`
	Created int64            `json:"created"`
}

// DoubaoChoice 豆包选择项
type DoubaoChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// convertToDoubaoRequest 转换为豆包请求格式
func (p *DoubaoProvider) convertToDoubaoRequest(req ChatRequest) DoubaoRequest {
	return DoubaoRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}
}

// convertFromDoubaoResponse 从豆包响应转换为统一格式
func (p *DoubaoProvider) convertFromDoubaoResponse(resp DoubaoResponse) ChatResponse {
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
