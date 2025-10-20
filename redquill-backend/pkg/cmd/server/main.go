package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"redquill-backend/pkg/config"
	"redquill-backend/pkg/server"
	"syscall"

	"redquill-backend/pkg/utils"
)

func main() {
	// Load configuration (env + defaults)
	cfg := config.Load()

	// Init MongoDB connection
	mongoClient, err := utils.Connect(context.Background(), cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() { _ = mongoClient.Disconnect(context.Background()) }()

	// Build and run HTTP server
	hs := server.NewHTTPServer(cfg, mongoClient)

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := hs.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	hs.Shutdown()
}
