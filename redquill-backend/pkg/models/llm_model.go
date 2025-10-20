// Package models
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: llm_model.go
/@Description: LLM model data structure
/*/

package models

// LLMModel LLM模型
type LLMModel struct {
	ID               string         `json:"id" bson:"_id,omitempty"`
	Name             string         `json:"name" bson:"name"`
	ModelID          string         `json:"model_id" bson:"model_id"`
	DisplayName      string         `json:"display_name" bson:"display_name"`
	Description      string         `json:"description" bson:"description"`
	Capabilities     []string       `json:"capabilities" bson:"capabilities"`
	TemperatureRange []float64      `json:"temperature_range" bson:"temperature_range"`
	CostPerToken     float64        `json:"cost_per_token" bson:"cost_per_token"`
	Status           string         `json:"status" bson:"status"`
	Config           LLMModelConfig `json:"config" bson:"config"`
	UsageCount       int64          `json:"usage_count" bson:"usage_count"`
	CreatorID        string         `json:"creator_id" bson:"creator_id"`
	Creator          string         `json:"creator" bson:"creator"`
	Ctime            int64          `json:"ctime" bson:"ctime"`
	Mtime            int64          `json:"mtime" bson:"mtime"`
}

// LLMModelConfig LLM模型配置
type LLMModelConfig struct {
	Provider    string  `json:"provider" bson:"provider"`
	APIKey      string  `json:"api_key" bson:"api_key"`
	BaseURL     string  `json:"base_url" bson:"base_url"`
	ModelName   string  `json:"model_name" bson:"model_name"`
	Temperature float64 `json:"temperature" bson:"temperature"`
	MaxTokens   int     `json:"max_tokens" bson:"max_tokens"`
	Timeout     int     `json:"timeout" bson:"timeout"`
}

// LLMModelTestRequest LLM模型测试请求
type LLMModelTestRequest struct {
	Messages []LLMTestMessage `json:"messages"`
	Stream   bool             `json:"stream,omitempty"`
}

// LLMMessage LLM消息
type LLMTestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMModelTestResponse LLM模型测试响应
type LLMModelTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// LLMModelServiceRequest LLM模型服务请求
type LLMModelServiceRequest struct {
	Messages []LLMTestMessage `json:"messages" binding:"required"`
	Stream   bool             `json:"stream,omitempty"`
	UserID   string           `json:"user_id,omitempty"` // 调用用户ID，用于统计
}

// LLMModelServiceResponse LLM模型服务响应
type LLMModelServiceResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Data        string `json:"data,omitempty"`
	Error       string `json:"error,omitempty"`
	UsageCount  int64  `json:"usage_count,omitempty"` // 当前模型使用次数
	TokenCount  int64  `json:"token_count,omitempty"` // 本次调用消耗的token数
}
