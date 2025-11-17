// Package models
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_template_model.go
/@Description: Prompt template model data structure
/*/

package models

// PromptTemplate Prompt模板
type PromptTemplate struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"` // story_core|worldview|character|chapter|quality_review
	Phase       string `json:"phase" bson:"phase"` // 创作阶段
	Content     string `json:"content" bson:"content"`
	Variables   []string `json:"variables" bson:"variables"`
	Description string `json:"description" bson:"description"`
	UsageCount  int64  `json:"usage_count" bson:"usage_count"`
	CreatorID   string `json:"creator_id" bson:"creator_id"`
	Creator     string `json:"creator" bson:"creator"`
	Ctime       int64  `json:"ctime" bson:"ctime"`
	Mtime       int64  `json:"mtime" bson:"mtime"`
}

// GenerationRequest 生成请求
type GenerationRequest struct {
	NovelID      string                 `json:"novel_id" binding:"required"`
	LLMModelID   string                 `json:"llm_model_id" binding:"required"`
	InputData    map[string]interface{} `json:"input_data" binding:"required"`
	TemplateType string                 `json:"template_type" binding:"required"`
	Stream       bool                   `json:"stream,omitempty"`
}

// GenerationResponse 生成响应
type GenerationResponse struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	UsageCount  int64                  `json:"usage_count,omitempty"`
	TokenCount  int64                  `json:"token_count,omitempty"`
}
