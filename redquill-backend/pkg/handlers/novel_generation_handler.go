// Package handlers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: novel_generation_handler.go
/@Description: Novel generation handlers implementation
/*/

package handlers

import (
	"net/http"
	"redquill-backend/pkg/models"
	"redquill-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// GenerateStoryCoreHandler 生成故事核心
func GenerateStoryCoreHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID    string                 `json:"novel_id" binding:"required"`
			LLMModelID string                 `json:"llm_model_id" binding:"required"`
			InputData  map[string]interface{} `json:"input_data" binding:"required"`
			Stream     bool                   `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 如果请求流式响应
		if req.Stream {
			GenerateStoryCoreStreamHandler(client, dbName, req.NovelID, req.LLMModelID, req.InputData)(c)
			return
		}

		storyCore, err := services.NewNovelGenerationService(client, dbName).GenerateStoryCore(
			c.Request.Context(),
			req.NovelID,
			req.LLMModelID,
			req.InputData,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, storyCore)
	}
}

// GenerateStoryCoreStreamHandler 流式生成故事核心
func GenerateStoryCoreStreamHandler(client *mongo.Client, dbName string, novelID, llmModelID string, inputData map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置流式响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Cache-Control")

		// 创建流式生成请求
		generationReq := models.GenerationRequest{
			NovelID:      novelID,
			LLMModelID:   llmModelID,
			InputData:    inputData,
			TemplateType: "story_core",
			Stream:       true,
		}

		// 调用流式生成服务
		templateService := services.NewPromptTemplateService(client, dbName)
		response, err := templateService.GenerateWithLLMStream(c.Request.Context(), generationReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 发送流式数据
		chunkCount := 0
		for chunk := range response {
			chunkCount++
			if chunk.Error != nil {
				c.SSEvent("error", gin.H{"error": chunk.Error.Error()})
				return
			}

			// 发送数据块
			c.SSEvent("data", gin.H{
				"content": chunk.Content,
				"done":    chunk.Done,
			})
			c.Writer.Flush()

			if chunk.Done {
				break
			}
		}
	}
}

// GenerateWorldviewHandler 生成世界观
func GenerateWorldviewHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID    string                 `json:"novel_id" binding:"required"`
			LLMModelID string                 `json:"llm_model_id" binding:"required"`
			InputData  map[string]interface{} `json:"input_data" binding:"required"`
			Stream     bool                   `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 如果请求流式响应
		if req.Stream {
			GenerateWorldviewStreamHandler(client, dbName, req.NovelID, req.LLMModelID, req.InputData)(c)
			return
		}

		worldview, err := services.NewNovelGenerationService(client, dbName).GenerateWorldview(
			c.Request.Context(),
			req.NovelID,
			req.LLMModelID,
			req.InputData,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, worldview)
	}
}

// GenerateCharacterHandler 生成角色
func GenerateCharacterHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID    string                 `json:"novel_id" binding:"required"`
			LLMModelID string                 `json:"llm_model_id" binding:"required"`
			InputData  map[string]interface{} `json:"input_data" binding:"required"`
			Stream     bool                   `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 如果请求流式响应
		if req.Stream {
			GenerateCharacterStreamHandler(client, dbName, req.NovelID, req.LLMModelID, req.InputData)(c)
			return
		}

		character, err := services.NewNovelGenerationService(client, dbName).GenerateCharacter(
			c.Request.Context(),
			req.NovelID,
			req.LLMModelID,
			req.InputData,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, character)
	}
}

// GenerateChapterHandler 生成章节
func GenerateChapterHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID    string                 `json:"novel_id" binding:"required"`
			LLMModelID string                 `json:"llm_model_id" binding:"required"`
			InputData  map[string]interface{} `json:"input_data" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chapter, err := services.NewNovelGenerationService(client, dbName).GenerateChapter(
			c.Request.Context(),
			req.NovelID,
			req.LLMModelID,
			req.InputData,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, chapter)
	}
}

// GenerateWithLLMHandler 通用LLM生成
func GenerateWithLLMHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID      string                 `json:"novel_id" binding:"required"`
			LLMModelID   string                 `json:"llm_model_id" binding:"required"`
			InputData    map[string]interface{} `json:"input_data" binding:"required"`
			TemplateType string                 `json:"template_type" binding:"required"`
			Stream       bool                   `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		generationReq := models.GenerationRequest{
			NovelID:      req.NovelID,
			LLMModelID:   req.LLMModelID,
			InputData:    req.InputData,
			TemplateType: req.TemplateType,
			Stream:       req.Stream,
		}

		response, err := services.NewPromptTemplateService(client, dbName).GenerateWithLLM(
			c.Request.Context(),
			generationReq,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// GenerateWorldviewStreamHandler 流式生成世界观
func GenerateWorldviewStreamHandler(client *mongo.Client, dbName string, novelID, llmModelID string, inputData map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置流式响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Cache-Control")

		// 创建流式生成请求
		generationReq := models.GenerationRequest{
			NovelID:      novelID,
			LLMModelID:   llmModelID,
			InputData:    inputData,
			TemplateType: "worldview",
			Stream:       true,
		}

		// 调用流式生成服务
		templateService := services.NewPromptTemplateService(client, dbName)
		response, err := templateService.GenerateWithLLMStream(c.Request.Context(), generationReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 发送流式数据
		chunkCount := 0
		for chunk := range response {
			chunkCount++
			if chunk.Error != nil {
				c.SSEvent("error", gin.H{"error": chunk.Error.Error()})
				return
			}

			// 发送数据块
			c.SSEvent("data", gin.H{
				"content": chunk.Content,
				"done":    chunk.Done,
			})
			c.Writer.Flush()

			if chunk.Done {
				break
			}
		}
	}
}

// GenerateCharacterStreamHandler 流式生成角色
func GenerateCharacterStreamHandler(client *mongo.Client, dbName string, novelID, llmModelID string, inputData map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置流式响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Cache-Control")

		// 创建流式生成请求
		generationReq := models.GenerationRequest{
			NovelID:      novelID,
			LLMModelID:   llmModelID,
			InputData:    inputData,
			TemplateType: "character",
			Stream:       true,
		}

		// 调用流式生成服务
		templateService := services.NewPromptTemplateService(client, dbName)
		response, err := templateService.GenerateWithLLMStream(c.Request.Context(), generationReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 发送流式数据
		chunkCount := 0
		for chunk := range response {
			chunkCount++
			if chunk.Error != nil {
				c.SSEvent("error", gin.H{"error": chunk.Error.Error()})
				return
			}

			// 发送数据块
			c.SSEvent("data", gin.H{
				"content": chunk.Content,
				"done":    chunk.Done,
			})
			c.Writer.Flush()

			if chunk.Done {
				break
			}
		}
	}
}

// GenerateCharactersFromOutlineHandler 根据大纲流式生成角色
func GenerateCharactersFromOutlineHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID          string                 `json:"novel_id" binding:"required"`
			LLMModelID       string                 `json:"llm_model_id" binding:"required"`
			OutlineContent   string                 `json:"outline_content" binding:"required"`
			StoryCore        string                 `json:"story_core" binding:"required"`
			Worldview        string                 `json:"worldview" binding:"required"`
			UserRequirements string                 `json:"user_requirements"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 设置流式响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Cache-Control")

		// 构建输入数据
		inputData := map[string]interface{}{
			"outline_content":   req.OutlineContent,
			"story_core":        req.StoryCore,
			"worldview":         req.Worldview,
			"user_requirements": req.UserRequirements,
		}

		// 创建流式生成请求
		generationReq := models.GenerationRequest{
			NovelID:      req.NovelID,
			LLMModelID:   req.LLMModelID,
			InputData:    inputData,
			TemplateType: "batch_character",
			Stream:       true,
		}

		// 调用流式生成服务
		templateService := services.NewPromptTemplateService(client, dbName)
		response, err := templateService.GenerateWithLLMStream(c.Request.Context(), generationReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 发送流式数据
		chunkCount := 0
		for chunk := range response {
			chunkCount++
			if chunk.Error != nil {
				c.SSEvent("error", gin.H{"error": chunk.Error.Error()})
				return
			}

			// 发送数据块
			c.SSEvent("data", gin.H{
				"content": chunk.Content,
				"done":    chunk.Done,
			})
			c.Writer.Flush()

			if chunk.Done {
				break
			}
		}
	}
}
