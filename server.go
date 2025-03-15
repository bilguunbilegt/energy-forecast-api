package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health Check Handler
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func main() {
	router := gin.Default()

	// Register API routes
	router.POST("/predict", predictHandler) // Uses predictHandler from predictHandler.go
	router.GET("/health", healthCheck)      // Health check endpoint

	fmt.Println("API Server is running on port 5000")
	log.Fatal(router.Run("0.0.0.0:5000")) // logs errors
}
