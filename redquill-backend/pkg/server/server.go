package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"redquill-backend/pkg/config"
	"redquill-backend/pkg/middleware"
	"redquill-backend/pkg/routes"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPServer struct {
	engine *gin.Engine
	server *http.Server
}

func NewHTTPServer(cfg config.Config, mongoClient *mongo.Client) *HTTPServer {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(gin.Logger())

	routes.Register(engine, cfg, mongoClient)

	hs := &HTTPServer{
		engine: engine,
	}
	hs.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           engine,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return hs
}

func (s *HTTPServer) Run() error {
	log.Printf("HTTP server listening on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
