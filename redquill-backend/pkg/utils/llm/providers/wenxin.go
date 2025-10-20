// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: wenxin.go
/@Description: 文心一言 provider implementation
/*/

package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// WenxinProvider 文心一言提供商
type WenxinProvider struct {
	config LLMConfig
	client *http.Client
	stream *StreamProcessor
}

// NewWenxinProvider 创建文心一言提供商
func NewWenxinProvider(config LLMConfig, client *http.Client) *WenxinProvider {
	return &WenxinProvider{
		config: config,
		client: client,
		stream: NewStreamProcessor(client),
	}
}

// Chat 同步聊天
func (p *WenxinProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = false

	// 文心一言API需要转换请求格式
	wenxinReq := p.convertToWenxinRequest(req)
	reqBody, err := json.Marshal(wenxinReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}

	// 文心一言需要获取access_token
	accessToken, err := p.getAccessToken(ctx)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeAuth),
			Message: fmt.Sprintf("get access token error: %v", err),
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("create request error: %v", err),
		}
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

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

	// 文心一言API响应格式转换
	var wenxinResp WenxinResponse
	if err := json.Unmarshal(body, &wenxinResp); err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeServer),
			Message: fmt.Sprintf("unmarshal response error: %v", err),
		}
	}

	// 转换为统一格式
	chatResp := p.convertFromWenxinResponse(wenxinResp)
	return &chatResp, nil
}

// ChatStream 流式聊天
func (p *WenxinProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	req.Stream = true

	// 文心一言API需要转换请求格式
	wenxinReq := p.convertToWenxinRequest(req)
	reqBody, err := json.Marshal(wenxinReq)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeInvalidRequest),
			Message: fmt.Sprintf("marshal request error: %v", err),
		}
	}

	// 文心一言需要获取access_token
	accessToken, err := p.getAccessToken(ctx)
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeAuth),
			Message: fmt.Sprintf("get access token error: %v", err),
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, &LLMError{
			Type:    string(ErrorTypeNetwork),
			Message: fmt.Sprintf("create request error: %v", err),
		}
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
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
func (p *WenxinProvider) Health(ctx context.Context) error {
	// 文心一言通过获取access_token来检查健康状态
	_, err := p.getAccessToken(ctx)
	return err
}

// Models 获取模型列表
func (p *WenxinProvider) Models(ctx context.Context) ([]Model, error) {
	// 文心一言返回预定义的模型列表
	return []Model{
		{ID: "ernie-bot", Name: "文心一言", Owner: "百度"},
		{ID: "ernie-bot-turbo", Name: "文心一言 Turbo", Owner: "百度"},
		{ID: "ernie-bot-4", Name: "文心一言 4.0", Owner: "百度"},
	}, nil
}

// getAccessToken 获取文心一言access_token
func (p *WenxinProvider) getAccessToken(ctx context.Context) (string, error) {
	// 这里需要根据文心一言的API文档实现获取access_token的逻辑
	// 通常需要client_id和client_secret
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", p.config.APIKey)    // 这里假设APIKey是client_id
	params.Set("client_secret", p.config.Model) // 这里假设Model是client_secret

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/oauth/2.0/token", bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("get access token error: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}

	return tokenResp.AccessToken, nil
}

// WenxinRequest 文心一言请求格式
type WenxinRequest struct {
	Messages     []Message `json:"messages"`
	Stream       bool      `json:"stream,omitempty"`
	Temperature  float64   `json:"temperature,omitempty"`
	TopP         float64   `json:"top_p,omitempty"`
	PenaltyScore float64   `json:"penalty_score,omitempty"`
}

// WenxinResponse 文心一言响应格式
type WenxinResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Result  string      `json:"result"`
	IsEnd   bool        `json:"is_end"`
	Usage   WenxinUsage `json:"usage"`
}

// WenxinUsage 文心一言使用统计
type WenxinUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// convertToWenxinRequest 转换为文心一言请求格式
func (p *WenxinProvider) convertToWenxinRequest(req ChatRequest) WenxinRequest {
	return WenxinRequest{
		Messages:     req.Messages,
		Stream:       req.Stream,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
		PenaltyScore: req.FrequencyPenalty,
	}
}

// convertFromWenxinResponse 从文心一言响应转换为统一格式
func (p *WenxinProvider) convertFromWenxinResponse(resp WenxinResponse) ChatResponse {
	choices := []Choice{
		{
			Index: 0,
			Message: Message{
				Role:    "assistant",
				Content: resp.Result,
			},
			FinishReason: "stop",
		},
	}

	usage := Usage{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}

	return ChatResponse{
		ID:      resp.ID,
		Model:   p.config.Model,
		Choices: choices,
		Usage:   usage,
		Created: resp.Created,
	}
}
