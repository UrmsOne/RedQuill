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

	}
}
