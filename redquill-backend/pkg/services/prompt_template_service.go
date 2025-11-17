// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_template_service.go
/@Description: Prompt template service implementation
/*/

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"redquill-backend/pkg/models"
	"redquill-backend/pkg/utils/llm"
)

// PromptTemplateService Prompt模板服务
type PromptTemplateService struct {
	client *mongo.Client
	dbName string
}

// NewPromptTemplateService 创建Prompt模板服务
func NewPromptTemplateService(client *mongo.Client, dbName string) *PromptTemplateService {
	return &PromptTemplateService{
		client: client,
		dbName: dbName,
	}
}

// PostPromptTemplates 创建Prompt模板
func (s *PromptTemplateService) PostPromptTemplates(ctx context.Context, name, templateType, phase, content string, variables []string, description, creatorID, creator string) (models.PromptTemplate, error) {
	coll := s.client.Database(s.dbName).Collection("prompt_templates")

	now := time.Now()
	template := models.PromptTemplate{
		Name:        name,
		Type:        templateType,
		Phase:       phase,
		Content:     content,
		Variables:   variables,
		Description: description,
		UsageCount:  0,
		CreatorID:   creatorID,
		Creator:     creator,
		Ctime:       now.Unix(),
		Mtime:       now.Unix(),
	}

	res, err := coll.InsertOne(ctx, template)
	if err != nil {
		return models.PromptTemplate{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		template.ID = oid.Hex()
	}

	return template, nil
}

// GetPromptTemplates 获取Prompt模板
func (s *PromptTemplateService) GetPromptTemplates(ctx context.Context, templateType string) ([]models.PromptTemplate, error) {
	coll := s.client.Database(s.dbName).Collection("prompt_templates")

	cursor, err := coll.Find(ctx, bson.M{"type": templateType})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []models.PromptTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// GenerateWithLLM 使用LLM生成内容
func (s *PromptTemplateService) GenerateWithLLM(ctx context.Context, req models.GenerationRequest) (models.GenerationResponse, error) {
	// 获取LLM模型配置
	llmModel, err := s.getLLMModel(ctx, req.LLMModelID)
	if err != nil {
		return models.GenerationResponse{
			Success: false,
			Message: "LLM model not found",
			Error:   err.Error(),
		}, nil
	}

	// 获取Prompt模板
	template, err := s.getPromptTemplate(ctx, req.TemplateType)
	if err != nil {
		return models.GenerationResponse{
			Success: false,
			Message: "Prompt template not found",
			Error:   err.Error(),
		}, nil
	}

	// 构建完整的Prompt
	fullPrompt, err := s.buildPrompt(template.Content, req.InputData)
	if err != nil {
		return models.GenerationResponse{
			Success: false,
			Message: "Failed to build prompt",
			Error:   err.Error(),
		}, nil
	}

	// 创建LLM客户端
	timeout := llmModel.Config.Timeout
	if timeout <= 0 {
		timeout = 300 // 默认5分钟超时
	}
	llmConfig := llm.LLMConfig{
		Provider: llmModel.Config.Provider,
		BaseURL:  llmModel.Config.BaseURL,
		APIKey:   llmModel.Config.APIKey,
		Model:    llmModel.Config.ModelName,
		Timeout:  time.Duration(timeout) * time.Second,
	}

	client, err := llm.NewClient(llmConfig)
	if err != nil {
		return models.GenerationResponse{
			Success: false,
			Message: "Failed to create LLM client",
			Error:   err.Error(),
		}, nil
	}

	// 构建消息
	messages := []llm.Message{
		{
			Role:    "user",
			Content: fullPrompt,
		},
	}

	// 执行生成
	chatReq := llm.ChatRequest{
		Model:       llmModel.Config.ModelName,
		Messages:    messages,
		Stream:      req.Stream,
		Temperature: llmModel.Config.Temperature,
		MaxTokens:   llmModel.Config.MaxTokens,
	}

	var response string
	var tokenCount int64

	if req.Stream {
		// 流式生成
		stream, err := client.ChatStream(ctx, chatReq)
		if err != nil {
			return models.GenerationResponse{
				Success: false,
				Message: "Stream generation failed",
				Error:   err.Error(),
			}, nil
		}

		for chunk := range stream {
			if chunk.Error != nil {
				return models.GenerationResponse{
					Success: false,
					Message: "Stream generation error",
					Error:   chunk.Error.Error(),
				}, nil
			}

			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					response += choice.Delta.Content
				}
			}

			// 统计token使用量
			if chunk.Usage != nil {
				tokenCount += chunk.Usage.TotalTokens
			}
		}
	} else {
		// 同步生成
		resp, err := client.Chat(ctx, chatReq)
		if err != nil {
			return models.GenerationResponse{
				Success: false,
				Message: "Chat generation failed",
				Error:   err.Error(),
			}, nil
		}

		if len(resp.Choices) > 0 {
			response = resp.Choices[0].Message.Content
		}

		// 统计token使用量
		if resp.Usage != nil {
			tokenCount = resp.Usage.TotalTokens
		}
	}

	// 解析响应为结构化数据
	structuredData, err := s.parseResponse(response, req.TemplateType)
	if err != nil {
		return models.GenerationResponse{
			Success: false,
			Message: "Failed to parse response",
			Error:   err.Error(),
		}, nil
	}

	// 更新模型使用次数
	s.updateLLMModelUsage(ctx, req.LLMModelID)

	// 更新模板使用次数
	s.updateTemplateUsage(ctx, req.TemplateType)

	return models.GenerationResponse{
		Success:    true,
		Message:    "Generation successful",
		Data:       structuredData,
		UsageCount: llmModel.UsageCount + 1,
		TokenCount: tokenCount,
	}, nil
}

// getLLMModel 获取LLM模型
func (s *PromptTemplateService) getLLMModel(ctx context.Context, modelID string) (models.LLMModel, error) {
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, err := primitive.ObjectIDFromHex(modelID)
	if err != nil {
		return models.LLMModel{}, errors.New("invalid model id")
	}

	var llmModel models.LLMModel
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&llmModel); err != nil {
		return models.LLMModel{}, err
	}

	return llmModel, nil
}

// getPromptTemplate 获取Prompt模板
func (s *PromptTemplateService) getPromptTemplate(ctx context.Context, templateType string) (models.PromptTemplate, error) {
	coll := s.client.Database(s.dbName).Collection("prompt_templates")

	var template models.PromptTemplate
	if err := coll.FindOne(ctx, bson.M{"type": templateType}).Decode(&template); err != nil {
		return models.PromptTemplate{}, err
	}

	return template, nil
}

// buildPrompt 构建完整Prompt
func (s *PromptTemplateService) buildPrompt(templateContent string, inputData map[string]interface{}) (string, error) {
	prompt := templateContent

	// 替换变量
	for key, value := range inputData {
		placeholder := fmt.Sprintf("{%s}", key)
		valueStr := fmt.Sprintf("%v", value)
		prompt = strings.ReplaceAll(prompt, placeholder, valueStr)
	}

	return prompt, nil
}

// parseResponse 解析响应为结构化数据
func (s *PromptTemplateService) parseResponse(response, templateType string) (map[string]interface{}, error) {
	// 根据模板类型解析不同的JSON结构
	var result map[string]interface{}
	
	// 尝试解析JSON响应
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// 如果不是JSON，返回原始响应
		return map[string]interface{}{
			"raw_response":  response,
			"template_type": templateType,
		}, nil
	}
	
	// 根据模板类型返回相应的数据结构
	switch templateType {
	case "story_core":
		if concepts, ok := result["concepts"]; ok {
			return map[string]interface{}{
				"concepts": concepts,
			}, nil
		}
	case "worldview":
		return result, nil
	case "character":
		return result, nil
	case "chapter":
		return result, nil
	default:
		return result, nil
	}
	
	return result, nil
}

// updateLLMModelUsage 更新LLM模型使用次数
func (s *PromptTemplateService) updateLLMModelUsage(ctx context.Context, modelID string) {
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, _ := primitive.ObjectIDFromHex(modelID)
	coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$inc": bson.M{"usage_count": 1},
		"$set": bson.M{"mtime": time.Now().Unix()},
	})
}

// updateTemplateUsage 更新模板使用次数
func (s *PromptTemplateService) updateTemplateUsage(ctx context.Context, templateType string) {
	coll := s.client.Database(s.dbName).Collection("prompt_templates")
	coll.UpdateOne(ctx, bson.M{"type": templateType}, bson.M{
		"$inc": bson.M{"usage_count": 1},
		"$set": bson.M{"mtime": time.Now().Unix()},
	})
}

// StreamChunk 流式数据块
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   error  `json:"error,omitempty"`
}

// GenerateWithLLMStream 流式LLM生成
func (s *PromptTemplateService) GenerateWithLLMStream(ctx context.Context, req models.GenerationRequest) (<-chan StreamChunk, error) {
	// 获取LLM模型配置
	llmModel, err := s.getLLMModel(ctx, req.LLMModelID)
	if err != nil {
		ch := make(chan StreamChunk, 1)
		ch <- StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 获取Prompt模板
	template, err := s.getPromptTemplate(ctx, req.TemplateType)
	if err != nil {
		ch := make(chan StreamChunk, 1)
		ch <- StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 构建完整的Prompt
	fullPrompt, err := s.buildPrompt(template.Content, req.InputData)
	if err != nil {
		ch := make(chan StreamChunk, 1)
		ch <- StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 创建LLM客户端
	timeout := llmModel.Config.Timeout
	if timeout <= 0 {
		timeout = 300 // 默认5分钟超时
	}
	llmConfig := llm.LLMConfig{
		Provider: llmModel.Config.Provider,
		BaseURL:  llmModel.Config.BaseURL,
		APIKey:   llmModel.Config.APIKey,
		Model:    llmModel.Config.ModelName,
		Timeout:  time.Duration(timeout) * time.Second,
	}

	client, err := llm.NewClient(llmConfig)
	if err != nil {
		ch := make(chan StreamChunk, 1)
		ch <- StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 构建消息
	messages := []llm.Message{
		{
			Role:    "user",
			Content: fullPrompt,
		},
	}

	// 执行流式生成
	chatReq := llm.ChatRequest{
		Model:       llmModel.Config.ModelName,
		Messages:    messages,
		Stream:      true,
		Temperature: llmModel.Config.Temperature,
		MaxTokens:   llmModel.Config.MaxTokens,
	}

	stream, err := client.ChatStream(ctx, chatReq)
	if err != nil {
		ch := make(chan StreamChunk, 1)
		ch <- StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 转换流式响应
	result := make(chan StreamChunk, 100)
	go func() {
		defer close(result)
		defer s.updateLLMModelUsage(ctx, req.LLMModelID)
		defer s.updateTemplateUsage(ctx, req.TemplateType)

		for chunk := range stream {
			if chunk.Error != nil {
				result <- StreamChunk{Error: chunk.Error}
				return
			}

			// 提取内容
			content := ""
			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					content += choice.Delta.Content
				}
			}

			if content != "" {
				result <- StreamChunk{
					Content: content,
					Done:    false,
				}
			}

			// 检查是否完成
			if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != "" {
				result <- StreamChunk{
					Content: "",
					Done:    true,
				}
				return
			}
		}
	}()

	return result, nil
}
