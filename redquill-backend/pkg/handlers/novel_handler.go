// Package handlers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: novel_handler.go
/@Description: Novel handlers' implementation
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

// PostNovelsHandler 创建小说
func PostNovelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title            string                  `json:"title" binding:"required"`
			Status           string                  `json:"status" binding:"required"`
			CurrentPhase     string                  `json:"current_phase" binding:"required"`
			ProjectBlueprint models.ProjectBlueprint `json:"project_blueprint" binding:"required"`
			AIContext        models.AIContext        `json:"ai_context"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取当前用户信息（从JWT中）
		authorID := c.GetString("uid")

		novel, err := services.NewNovelService(client, dbName).PostNovels(
			c.Request.Context(),
			req.Title,
			authorID,
			req.Status,
			req.CurrentPhase,
			req.ProjectBlueprint,
			req.AIContext,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, novel)
	}
}

// GetNovelsHandler 获取小说详情
func GetNovelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		novel, err := services.NewNovelService(client, dbName).GetNovels(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, novel)
	}
}

// PutNovelsHandler 更新小说
func PutNovelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Title            *string                  `json:"title"`
			Status           *string                  `json:"status"`
			CurrentPhase     *string                  `json:"current_phase"`
			ProjectBlueprint *models.ProjectBlueprint `json:"project_blueprint"`
			AIContext        *models.AIContext        `json:"ai_context"`
			ExtraInfo        *map[string]interface{}  `json:"extra_info"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		novel, err := services.NewNovelService(client, dbName).PutNovels(
			c.Request.Context(),
			id,
			req.Title,
			req.Status,
			req.CurrentPhase,
			req.ProjectBlueprint,
			req.AIContext,
			req.ExtraInfo,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, novel)
	}
}

// DeleteNovelsHandler 删除小说
func DeleteNovelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := services.NewNovelService(client, dbName).DeleteNovels(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}

// ListNovelsHandler 分页查询小说列表
func ListNovelsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, size, sortExpr, q := common.ParseCommonQueryParams(c.Request.URL.Query())
		result, err := services.NewNovelService(client, dbName).ListNovels(c.Request.Context(), page, size, sortExpr, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// PostStoryCoresHandler 创建故事核心
func PostStoryCoresHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID             string `json:"novel_id" binding:"required"`
			Title               string `json:"title" binding:"required"`
			CoreConflict        string `json:"core_conflict" binding:"required"`
			Theme               string `json:"theme" binding:"required"`
			Innovation          string `json:"innovation" binding:"required"`
			CommercialPotential string `json:"commercial_potential" binding:"required"`
			TargetAudience      string `json:"target_audience" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		storyCore, err := services.NewNovelService(client, dbName).PostStoryCores(
			c.Request.Context(),
			req.NovelID,
			req.Title,
			req.CoreConflict,
			req.Theme,
			req.Innovation,
			req.CommercialPotential,
			req.TargetAudience,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, storyCore)
	}
}

// GetStoryCoresHandler 获取故事核心列表
func GetStoryCoresHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		storyCores, err := services.NewNovelService(client, dbName).GetStoryCores(c.Request.Context(), novelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, storyCores)
	}
}

// PostWorldviewsHandler 创建世界观
func PostWorldviewsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID          string                  `json:"novel_id" binding:"required"`
			PowerSystem      models.PowerSystem      `json:"power_system" binding:"required"`
			SocietyStructure models.SocietyStructure `json:"society_structure" binding:"required"`
			Geography        models.Geography        `json:"geography" binding:"required"`
			SpecialRules     []string                `json:"special_rules"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		worldview, err := services.NewNovelService(client, dbName).PostWorldviews(
			c.Request.Context(),
			req.NovelID,
			req.PowerSystem,
			req.SocietyStructure,
			req.Geography,
			req.SpecialRules,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, worldview)
	}
}

// GetWorldviewsHandler 获取世界观
func GetWorldviewsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		worldview, err := services.NewNovelService(client, dbName).GetWorldviews(c.Request.Context(), novelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, worldview)
	}
}

// PostCharactersHandler 创建角色
func PostCharactersHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID        string                `json:"novel_id" binding:"required"`
			Name           string                `json:"name" binding:"required"`
			Type           string                `json:"type" binding:"required"`
			CoreAttributes models.CoreAttributes `json:"core_attributes" binding:"required"`
			SoulProfile    models.SoulProfile    `json:"soul_profile" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		character, err := services.NewNovelService(client, dbName).PostCharacters(
			c.Request.Context(),
			req.NovelID,
			req.Name,
			req.Type,
			req.CoreAttributes,
			req.SoulProfile,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, character)
	}
}

// GetCharactersHandler 获取角色列表
func GetCharactersHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		characters, err := services.NewNovelService(client, dbName).GetCharacters(c.Request.Context(), novelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, characters)
	}
}

// PostChaptersHandler 创建章节
func PostChaptersHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID              string                `json:"novel_id" binding:"required"`
			ChapterNumber        int                   `json:"chapter_number" binding:"required"`
			Title                string                `json:"title" binding:"required"`
			Content              string                `json:"content" binding:"required"`
			Summary              string                `json:"summary" binding:"required"`
			Outline              models.ChapterOutline `json:"outline" binding:"required"`
			QualityMetrics       models.QualityMetrics `json:"quality_metrics" binding:"required"`
			CharacterDevelopment map[string]string     `json:"character_development"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chapter, err := services.NewNovelService(client, dbName).PostChapters(
			c.Request.Context(),
			req.NovelID,
			req.ChapterNumber,
			req.Title,
			req.Content,
			req.Summary,
			req.Outline,
			req.QualityMetrics,
			req.CharacterDevelopment,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, chapter)
	}
}

// GetChaptersHandler 获取章节列表
func GetChaptersHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		chapters, err := services.NewNovelService(client, dbName).GetChapters(c.Request.Context(), novelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, chapters)
	}
}

// GetChapterHandler 获取单个章节
func GetChapterHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		chapter, err := services.NewNovelService(client, dbName).GetChapter(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, chapter)
	}
}

// PostWritingSessionsHandler 创建创作会话
func PostWritingSessionsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NovelID        string                `json:"novel_id" binding:"required"`
			CurrentChapter int                   `json:"current_chapter" binding:"required"`
			SessionContext models.SessionContext `json:"session_context" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session, err := services.NewNovelService(client, dbName).PostWritingSessions(
			c.Request.Context(),
			req.NovelID,
			req.CurrentChapter,
			req.SessionContext,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

// GetWritingSessionsHandler 获取创作会话
func GetWritingSessionsHandler(client *mongo.Client, dbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelID := c.Param("novel_id")
		session, err := services.NewNovelService(client, dbName).GetWritingSessions(c.Request.Context(), novelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}
