package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"redquill-backend/pkg/config"
	"redquill-backend/pkg/handlers"
	"redquill-backend/pkg/middleware"
)

func Register(r *gin.Engine, cfg config.Config, mongoClient *mongo.Client) {
	// Health
	r.GET("/healthz", handlers.HealthHandler(mongoClient))

	// API v1
	v1 := r.Group("/api/v1")
	{
		// auth
		v1.POST("/login", handlers.LoginHandler(mongoClient, cfg.DBName, cfg))

		// users
		v1.POST("/user", handlers.PostUsersHandler(mongoClient, cfg.DBName)) // registration

		auth := v1.Group("")
		auth.Use(middleware.AuthRequired(cfg))
		auth.GET("/users", handlers.ListUsersHandler(mongoClient, cfg.DBName))
		auth.GET("/user/:id", handlers.GetUsersHandler(mongoClient, cfg.DBName))
		auth.PUT("/user/:id", handlers.PutUsersHandler(mongoClient, cfg.DBName))
		auth.DELETE("/user/:id", handlers.DeleteUsersHandler(mongoClient, cfg.DBName))

		// LLM models
		auth.POST("/llm-model", handlers.PostLLMModelsHandler(mongoClient, cfg.DBName))
		auth.GET("/llm-models", handlers.ListLLMModelsHandler(mongoClient, cfg.DBName))
		auth.GET("/llm-model/:id", handlers.GetLLMModelsHandler(mongoClient, cfg.DBName))
		auth.PUT("/llm-model/:id", handlers.PutLLMModelsHandler(mongoClient, cfg.DBName))
		auth.DELETE("/llm-model/:id", handlers.DeleteLLMModelsHandler(mongoClient, cfg.DBName))
		auth.POST("/llm-model/:id/test", handlers.TestLLMModelsHandler(mongoClient, cfg.DBName))
		auth.POST("/llm-model/:id/service", handlers.ServiceLLMModelsHandler(mongoClient, cfg.DBName))

		// Prompts
		auth.POST("/prompt", handlers.PostPromptsHandler(mongoClient, cfg.DBName))
		auth.GET("/prompts", handlers.ListPromptsHandler(mongoClient, cfg.DBName))
		auth.GET("/prompt/:id", handlers.GetPromptsHandler(mongoClient, cfg.DBName))
		auth.PUT("/prompt/:id", handlers.PutPromptsHandler(mongoClient, cfg.DBName))
		auth.DELETE("/prompt/:id", handlers.DeletePromptsHandler(mongoClient, cfg.DBName))

		// Novels
		auth.POST("/novel", handlers.PostNovelsHandler(mongoClient, cfg.DBName))
		auth.GET("/novels", handlers.ListNovelsHandler(mongoClient, cfg.DBName))
		auth.GET("/novel/:id", handlers.GetNovelsHandler(mongoClient, cfg.DBName))
		auth.PUT("/novel/:id", handlers.PutNovelsHandler(mongoClient, cfg.DBName))
		auth.DELETE("/novel/:id", handlers.DeleteNovelsHandler(mongoClient, cfg.DBName))

		// Story cores - 使用不同的路径前缀避免冲突
		auth.POST("/story-core", handlers.PostStoryCoresHandler(mongoClient, cfg.DBName))
		auth.GET("/story-cores/:novel_id", handlers.GetStoryCoresHandler(mongoClient, cfg.DBName))

		// Worldviews - 使用不同的路径前缀避免冲突
		auth.POST("/worldview", handlers.PostWorldviewsHandler(mongoClient, cfg.DBName))
		auth.GET("/worldview/:novel_id", handlers.GetWorldviewsHandler(mongoClient, cfg.DBName))

		// Characters - 使用不同的路径前缀避免冲突
		auth.POST("/character", handlers.PostCharactersHandler(mongoClient, cfg.DBName))
		auth.GET("/characters/:novel_id", handlers.GetCharactersHandler(mongoClient, cfg.DBName))

		// Chapters - 使用不同的路径前缀避免冲突
		auth.POST("/chapter", handlers.PostChaptersHandler(mongoClient, cfg.DBName))
		auth.GET("/chapters/:novel_id", handlers.GetChaptersHandler(mongoClient, cfg.DBName))
		auth.GET("/chapter/:id", handlers.GetChapterHandler(mongoClient, cfg.DBName))

		// Writing sessions - 使用不同的路径前缀避免冲突
		auth.POST("/writing-session", handlers.PostWritingSessionsHandler(mongoClient, cfg.DBName))
		auth.GET("/writing-session/:novel_id", handlers.GetWritingSessionsHandler(mongoClient, cfg.DBName))

		// Outlines - 大纲管理
		auth.POST("/outline", handlers.PostOutlinesHandler(mongoClient, cfg.DBName))
		auth.GET("/outlines", handlers.ListOutlinesHandler(mongoClient, cfg.DBName))
		auth.GET("/outline/:id", handlers.GetOutlineHandler(mongoClient, cfg.DBName))
		auth.PUT("/outline/:id", handlers.PutOutlinesHandler(mongoClient, cfg.DBName))
		auth.DELETE("/outline/:id", handlers.DeleteOutlinesHandler(mongoClient, cfg.DBName))
		auth.GET("/outlines/:novel_id", handlers.GetOutlinesHandler(mongoClient, cfg.DBName))

		// Novel generation - AI生成功能
		auth.POST("/generate/story-core", handlers.GenerateStoryCoreHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/worldview", handlers.GenerateWorldviewHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/character", handlers.GenerateCharacterHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/characters-from-outline", handlers.GenerateCharactersFromOutlineHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/outline", handlers.GenerateOutlineHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/chapter", handlers.GenerateChapterHandler(mongoClient, cfg.DBName))
		auth.POST("/generate/llm", handlers.GenerateWithLLMHandler(mongoClient, cfg.DBName))

	}
}
