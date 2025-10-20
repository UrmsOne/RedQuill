package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HealthHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		mongoErr := client.Ping(ctx, nil)

		status := http.StatusOK
		resp := gin.H{"status": "ok"}
		if mongoErr != nil {
			status = http.StatusServiceUnavailable
			resp["mongo"] = mongoErr.Error()
		}
		c.JSON(status, resp)
	}
}


