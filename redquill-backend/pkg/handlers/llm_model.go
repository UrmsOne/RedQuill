// Package handlers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: llm_model.go
/@Description: LLM model handlers implementation
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

// PostLLMModelsHandler 创建LLM模型
func PostLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string                `json:"name" binding:"required"`
			ModelID          string                `json:"model_id"`
			DisplayName      string                `json:"display_name" binding:"required"`
			Description      string                `json:"description"`
			Capabilities     []string              `json:"capabilities"`
			TemperatureRange []float64             `json:"temperature_range"`
			CostPerToken     float64               `json:"cost_per_token"`
			Status           string                `json:"status" binding:"required"`
			Config           models.LLMModelConfig `json:"config" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取当前用户信息（从JWT中）
		creatorID := c.GetString("uid")
		creator := c.GetString("username")

		llmModel, err := services.NewLLMModelService(client, dbName).PostLLMModels(
			c.Request.Context(),
			req.Name,
			req.ModelID,
			req.DisplayName,
			req.Description,
			req.Capabilities,
			req.TemperatureRange,
			req.CostPerToken,
			req.Status,
			req.Config,
			creatorID,
			creator,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, llmModel)
	}
}

// GetLLMModelsHandler 获取LLM模型详情
func GetLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		llmModel, err := services.NewLLMModelService(client, dbName).GetLLMModels(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, llmModel)
	}
}

// PutLLMModelsHandler 更新LLM模型
func PutLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Name             *string                `json:"name"`
			ModelID          *string                `json:"model_id"`
			DisplayName      *string                `json:"display_name"`
			Description      *string                `json:"description"`
			Capabilities     *[]string              `json:"capabilities"`
			TemperatureRange *[]float64             `json:"temperature_range"`
			CostPerToken     *float64               `json:"cost_per_token"`
			Status           *string                `json:"status"`
			Config           *models.LLMModelConfig `json:"config"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		llmModel, err := services.NewLLMModelService(client, dbName).PutLLMModels(
			c.Request.Context(),
			id,
			req.Name,
			req.ModelID,
			req.DisplayName,
			req.Description,
			req.Capabilities,
			req.TemperatureRange,
			req.CostPerToken,
			req.Status,
			req.Config,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, llmModel)
	}
}

// DeleteLLMModelsHandler 删除LLM模型
func DeleteLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := services.NewLLMModelService(client, dbName).DeleteLLMModels(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}

// ListLLMModelsHandler 分页查询LLM模型列表
func ListLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())
		result, err := services.NewLLMModelService(client, dbName).ListLLMModels(c.Request.Context(), page, size, sortExpr, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// TestLLMModelsHandler 测试LLM模型
func TestLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Messages []models.LLMTestMessage `json:"messages" binding:"required"`
			Stream   bool                    `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 转换消息格式
		messages := make([]models.LLMTestMessage, len(req.Messages))
		for i, msg := range req.Messages {
			messages[i] = models.LLMTestMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}

		testReq := models.LLMModelTestRequest{
			Messages: messages,
			Stream:   req.Stream,
		}

		result, err := services.NewLLMModelService(client, dbName).TestLLMModel(c.Request.Context(), id, testReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// ServiceLLMModelsHandler 使用LLM模型提供服务
func ServiceLLMModelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Messages []models.LLMTestMessage `json:"messages" binding:"required"`
			Stream   bool                    `json:"stream,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取当前用户信息（从JWT中）
		userID := c.GetString("uid")

		// 构建服务请求
		serviceReq := models.LLMModelServiceRequest{
			Messages: req.Messages,
			Stream:   req.Stream,
			UserID:   userID,
		}

		result, err := services.NewLLMModelService(client, dbName).ServiceLLMModel(c.Request.Context(), id, serviceReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
