// Package services
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: llm_model.go
/@Description: LLM model service implementation
/*/

package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"redquill-backend/pkg/common"
	"redquill-backend/pkg/models"
	"redquill-backend/pkg/utils/llm"
)

// LLMModelService LLM模型服务
type LLMModelService struct {
	client *mongo.Client
	dbName string
}

// NewLLMModelService 创建LLM模型服务
func NewLLMModelService(client *mongo.Client, dbName string) *LLMModelService {
	return &LLMModelService{
		client: client,
		dbName: dbName,
	}
}

// PostLLMModels 创建LLM模型
func (s *LLMModelService) PostLLMModels(ctx context.Context, name, modelID, displayName, description string, capabilities []string, temperatureRange []float64, costPerToken float64, status string, config models.LLMModelConfig, creatorID, creator string) (models.LLMModel, error) {
	coll := s.client.Database(s.dbName).Collection("llm_models")

	// 检查模型名称是否已存在
	var existing models.LLMModel
	if err := coll.FindOne(ctx, bson.M{"name": name}).Decode(&existing); err == nil {
		return models.LLMModel{}, errors.New("model name already exists")
	}

	now := time.Now()
	llmModel := models.LLMModel{
		Name:             name,
		ModelID:          modelID,
		DisplayName:      displayName,
		Description:      description,
		Capabilities:     capabilities,
		TemperatureRange: temperatureRange,
		CostPerToken:     costPerToken,
		Status:           status,
		Config:           config,
		UsageCount:       0,
		Ctime:            now.Unix(),
		Mtime:            now.Unix(),
		CreatorID:        creatorID,
		Creator:          creator,
	}

	res, err := coll.InsertOne(ctx, llmModel)
	if err != nil {
		return models.LLMModel{}, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		llmModel.ID = oid.Hex()
	}

	return llmModel, nil
}

// GetLLMModels 获取LLM模型详情
func (s *LLMModelService) GetLLMModels(ctx context.Context, id string) (models.LLMModel, error) {
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.LLMModel{}, errors.New("invalid id")
	}

	var llmModel models.LLMModel
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&llmModel); err != nil {
		return models.LLMModel{}, err
	}

	return llmModel, nil
}

// PutLLMModels 更新LLM模型
func (s *LLMModelService) PutLLMModels(ctx context.Context, id string, name *string, modelID *string, displayName *string, description *string, capabilities *[]string, temperatureRange *[]float64, costPerToken *float64, status *string, config *models.LLMModelConfig) (models.LLMModel, error) {
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.LLMModel{}, errors.New("invalid id")
	}

	update := bson.M{"mtime": time.Now().Unix()}
	set := bson.M{}

	if name != nil {
		set["name"] = *name
	}
	if modelID != nil {
		set["model_id"] = *modelID
	}
	if displayName != nil {
		set["display_name"] = *displayName
	}
	if description != nil {
		set["description"] = *description
	}
	if capabilities != nil {
		set["capabilities"] = *capabilities
	}
	if temperatureRange != nil {
		set["temperature_range"] = *temperatureRange
	}
	if costPerToken != nil {
		set["cost_per_token"] = *costPerToken
	}
	if status != nil {
		set["status"] = *status
	}
	if config != nil {
		set["config"] = *config
	}

	for k, v := range set {
		update[k] = v
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out models.LLMModel
	if err := coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&out); err != nil {
		return models.LLMModel{}, err
	}

	return out, nil
}

// DeleteLLMModels 删除LLM模型
func (s *LLMModelService) DeleteLLMModels(ctx context.Context, id string) error {
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ListLLMModels 分页查询LLM模型列表
func (s *LLMModelService) ListLLMModels(ctx context.Context, page, pageSize int64, sortExpr, keyword string) (PagedLLMModels, error) {
	coll := s.client.Database(s.dbName).Collection("llm_models")

	// 构建过滤条件
	kwFilter := common.BuildKeywordFilter(keyword, []string{"name", "display_name", "description"})
	filter := common.MergeFilters(bson.M{}, kwFilter)

	// 构建选项
	sort := common.BuildSort(sortExpr)
	opts := common.BuildFindOptions(page, pageSize, sort, bson.M{})

	items, total, err := common.FindWithPagination[models.LLMModel](ctx, coll, filter, opts)
	if err != nil {
		return PagedLLMModels{}, err
	}

	totalPages := total / common.NormalizePageSize(pageSize)
	if total%common.NormalizePageSize(pageSize) != 0 {
		totalPages++
	}

	return PagedLLMModels{
		Items: items,
		Pagination: common.Pagination{
			Page:      common.NormalizePage(page),
			PageSize:  common.NormalizePageSize(pageSize),
			Total:     total,
			TotalPage: totalPages,
		},
	}, nil
}

// PagedLLMModels 分页LLM模型结果
type PagedLLMModels struct {
	Items      []models.LLMModel `json:"items"`
	Pagination common.Pagination `json:"pagination"`
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

// TestLLMModel 测试LLM模型
func (s *LLMModelService) TestLLMModel(ctx context.Context, id string, req models.LLMModelTestRequest) (models.LLMModelTestResponse, error) {
	// 获取模型信息
	llmModel, err := s.GetLLMModels(ctx, id)
	if err != nil {
		return models.LLMModelTestResponse{
			Success: false,
			Message: "Model not found",
			Error:   err.Error(),
		}, nil
	}

	// 检查模型状态
	if llmModel.Status != "active" {
		return models.LLMModelTestResponse{
			Success: false,
			Message: "Model is not active",
			Error:   "Model status is not active",
		}, nil
	}

	// 创建LLM客户端
	llmConfig := llm.LLMConfig{
		Provider: llmModel.Config.Provider,
		BaseURL:  llmModel.Config.BaseURL,
		APIKey:   llmModel.Config.APIKey,
		Model:    llmModel.Config.ModelName,
		Timeout:  time.Duration(llmModel.Config.Timeout) * time.Second,
	}

	client, err := llm.NewClient(llmConfig)
	if err != nil {
		return models.LLMModelTestResponse{
			Success: false,
			Message: "Failed to create LLM client",
			Error:   err.Error(),
		}, nil
	}

	// 转换消息格式
	messages := make([]llm.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 执行测试
	chatReq := llm.ChatRequest{
		Model:       llmModel.Config.ModelName,
		Messages:    messages,
		Stream:      req.Stream,
		Temperature: llmModel.Config.Temperature,
		MaxTokens:   llmModel.Config.MaxTokens,
	}

	// 同步测试
	resp, err := client.Chat(ctx, chatReq)
	if err != nil {
		return models.LLMModelTestResponse{
			Success: false,
			Message: "Chat test failed",
			Error:   err.Error(),
		}, nil
	}

	var response string
	if len(resp.Choices) > 0 {
		response = resp.Choices[0].Message.Content
	}

	return models.LLMModelTestResponse{
		Success: true,
		Message: "Chat test successful",
		Data:    response,
	}, nil
}

// ServiceLLMModel 使用LLM模型提供服务
func (s *LLMModelService) ServiceLLMModel(ctx context.Context, id string, req models.LLMModelServiceRequest) (models.LLMModelServiceResponse, error) {
	// 获取模型信息
	llmModel, err := s.GetLLMModels(ctx, id)
	if err != nil {
		return models.LLMModelServiceResponse{
			Success: false,
			Message: "Model not found",
			Error:   err.Error(),
		}, nil
	}

	// 检查模型状态
	if llmModel.Status != "active" {
		return models.LLMModelServiceResponse{
			Success: false,
			Message: "Model is not active",
			Error:   "Model status is not active",
		}, nil
	}

	// 创建LLM客户端
	llmConfig := llm.LLMConfig{
		Provider: llmModel.Config.Provider,
		BaseURL:  llmModel.Config.BaseURL,
		APIKey:   llmModel.Config.APIKey,
		Model:    llmModel.Config.ModelName,
		Timeout:  time.Duration(llmModel.Config.Timeout) * time.Second,
	}

	client, err := llm.NewClient(llmConfig)
	if err != nil {
		return models.LLMModelServiceResponse{
			Success: false,
			Message: "Failed to create LLM client",
			Error:   err.Error(),
		}, nil
	}

	// 转换消息格式
	messages := make([]llm.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 执行服务调用
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
		// 流式服务调用
		stream, err := client.ChatStream(ctx, chatReq)
		if err != nil {
			return models.LLMModelServiceResponse{
				Success: false,
				Message: "Stream service failed",
				Error:   err.Error(),
			}, nil
		}

		for chunk := range stream {
			if chunk.Error != nil {
				return models.LLMModelServiceResponse{
					Success: false,
					Message: "Stream service error",
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
		// 同步服务调用
		resp, err := client.Chat(ctx, chatReq)
		if err != nil {
			return models.LLMModelServiceResponse{
				Success: false,
				Message: "Chat service failed",
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

	// 更新模型使用次数
	coll := s.client.Database(s.dbName).Collection("llm_models")
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err = coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$inc": bson.M{"usage_count": 1},
		"$set": bson.M{"mtime": time.Now().Unix()},
	})

	if err != nil {
		// 记录错误但不影响服务响应
		// 可以考虑添加日志记录
	}

	return models.LLMModelServiceResponse{
		Success:    true,
		Message:    "Service call successful",
		Data:       response,
		UsageCount: llmModel.UsageCount + 1,
		TokenCount: tokenCount,
	}, nil
}

// TestLLMModelStream 流式测试LLM模型
func (s *LLMModelService) TestLLMModelStream(ctx context.Context, id string, req models.LLMModelTestRequest) (<-chan models.StreamChunk, error) {
	// 获取模型信息
	llmModel, err := s.GetLLMModels(ctx, id)
	if err != nil {
		ch := make(chan models.StreamChunk, 1)
		ch <- models.StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 检查模型状态
	if llmModel.Status != "active" {
		ch := make(chan models.StreamChunk, 1)
		ch <- models.StreamChunk{Error: errors.New("model is not active")}
		close(ch)
		return ch, nil
	}

	// 创建LLM客户端
	llmConfig := llm.LLMConfig{
		Provider: llmModel.Config.Provider,
		BaseURL:  llmModel.Config.BaseURL,
		APIKey:   llmModel.Config.APIKey,
		Model:    llmModel.Config.ModelName,
		Timeout:  time.Duration(llmModel.Config.Timeout) * time.Second,
	}

	client, err := llm.NewClient(llmConfig)
	if err != nil {
		ch := make(chan models.StreamChunk, 1)
		ch <- models.StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 转换消息格式
	messages := make([]llm.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 执行流式测试
	chatReq := llm.ChatRequest{
		Model:       llmModel.Config.ModelName,
		Messages:    messages,
		Stream:      true,
		Temperature: llmModel.Config.Temperature,
		MaxTokens:   llmModel.Config.MaxTokens,
	}

	stream, err := client.ChatStream(ctx, chatReq)
	if err != nil {
		ch := make(chan models.StreamChunk, 1)
		ch <- models.StreamChunk{Error: err}
		close(ch)
		return ch, nil
	}

	// 转换流式响应
	result := make(chan models.StreamChunk, 1000)
	go func() {
		defer close(result)

		for chunk := range stream {
			if chunk.Error != nil {
				result <- models.StreamChunk{Error: chunk.Error}
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
				result <- models.StreamChunk{
					Content: content,
					Done:    false,
				}
			}

			// 检查是否完成
			if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != "" {
				result <- models.StreamChunk{
					Content: "",
					Done:    true,
				}
				return
			}
		}
	}()

	return result, nil
}
