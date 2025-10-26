// Package handlers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/21 20:44
/@Name: outline_handler.go
/@Description: Outline handlers implementation
/*/

package handlers

import (
	"net/http"
	"redquill-backend/pkg/common"
	"redquill-backend/pkg/models"
	"redquill-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostOutlinesHandler 创建大纲
func PostOutlinesHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Outline
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		outline, err := services.NewNovelService(client, dbName).PostOutlines(
			c.Request.Context(),
			req,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outline)
	}
}

// GetOutlinesHandler 获取大纲列表
func GetOutlinesHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		if novelID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "novel_id is required"})
			return
		}

		outlines, err := services.NewNovelService(client, dbName).GetOutlines(
			c.Request.Context(),
			novelID,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, outlines)
	}
}

// GetOutlineHandler 获取单个大纲
func GetOutlineHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		outline, err := services.NewNovelService(client, dbName).GetOutline(
			c.Request.Context(),
			id,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, outline)
	}
}

// PutOutlinesHandler 更新大纲
func PutOutlinesHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		var req struct {
			Title     *string                `json:"title,omitempty"`
			Summary   *string                `json:"summary,omitempty"`
			Chapters  *[]models.ChapterInfo  `json:"chapters,omitempty"`
			StoryArcs *[]models.StoryArc     `json:"story_arcs,omitempty"`
			KeyThemes *[]string              `json:"key_themes,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		outline, err := services.NewNovelService(client, dbName).PutOutlines(
			c.Request.Context(),
			id,
			req.Title,
			req.Summary,
			req.Chapters,
			req.StoryArcs,
			req.KeyThemes,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, outline)
	}
}

// DeleteOutlinesHandler 删除大纲
func DeleteOutlinesHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		err := services.NewNovelService(client, dbName).DeleteOutlines(
			c.Request.Context(),
			id,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Outline deleted successfully"})
	}
}

// ListOutlinesHandler 分页查询大纲列表
func ListOutlinesHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())

		result, err := services.NewNovelService(client, dbName).ListOutlines(
			c.Request.Context(),
			page,
			size,
			sortExpr,
			q,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GenerateOutlineHandler 生成大纲
func GenerateOutlineHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
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
			GenerateOutlineStreamHandler(client, dbName, req.NovelID, req.LLMModelID, req.InputData)(c)
			return
		}

		outline, err := services.NewNovelGenerationService(client, dbName).GenerateOutline(
			c.Request.Context(),
			req.NovelID,
			req.LLMModelID,
			req.InputData,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, outline)
	}
}

// GenerateOutlineStreamHandler 流式生成大纲
func GenerateOutlineStreamHandler(client *mongo.Client, dbName string, novelID, llmModelID string, inputData map[string]interface{}) gin.HandlerFunc {
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
			TemplateType: "outline",
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
		for chunk := range response {
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
