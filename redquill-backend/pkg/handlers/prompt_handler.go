// Package handlers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: prompt_handler.go
/@Description: Prompt handlers' implementation
/*/

package handlers

import (
	"net/http"
	"redquill-backend/pkg/common"
	"redquill-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostPromptsHandler 创建Prompt
func PostPromptsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string   `json:"name" binding:"required"`
			Type        string   `json:"type" binding:"required"`
			Category    string   `json:"category"`
			Description string   `json:"description"`
			Content     string   `json:"content" binding:"required"`
			Variables   []string `json:"variables"`
			Tags        []string `json:"tags"`
			Public      bool     `json:"public"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取当前用户信息（从JWT中）
		creatorID := c.GetString("uid")
		username := c.GetString("username")

		prompt, err := services.NewPromptService(client, dbName).PostPrompts(
			c.Request.Context(),
			req.Name,
			req.Type,
			req.Category,
			req.Description,
			req.Content,
			req.Variables,
			req.Tags,
			req.Public,
			creatorID,
			username,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, prompt)
	}
}

// GetPromptsHandler 获取Prompt详情
func GetPromptsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		prompt, err := services.NewPromptService(client, dbName).GetPrompts(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, prompt)
	}
}

// PutPromptsHandler 更新Prompt
func PutPromptsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Name        *string   `json:"name"`
			Type        *string   `json:"type"`
			Category    *string   `json:"category"`
			Description *string   `json:"description"`
			Content     *string   `json:"content"`
			Variables   *[]string `json:"variables"`
			Tags        *[]string `json:"tags"`
			Public      *bool     `json:"public"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prompt, err := services.NewPromptService(client, dbName).PutPrompts(
			c.Request.Context(),
			id,
			req.Name,
			req.Type,
			req.Category,
			req.Description,
			req.Content,
			req.Variables,
			req.Tags,
			req.Public,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, prompt)
	}
}

// DeletePromptsHandler 删除Prompt
func DeletePromptsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := services.NewPromptService(client, dbName).DeletePrompts(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}

// ListPromptsHandler 分页查询Prompt列表
func ListPromptsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())
		result, err := services.NewPromptService(client, dbName).ListPrompts(c.Request.Context(), page, size, sortExpr, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
