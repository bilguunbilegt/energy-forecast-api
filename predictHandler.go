package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Request structure
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// Response structure
type PredictionResponse struct {
	EnergyKWh float64 `json:"predicted_energy_kwh"`
}

// Prediction API Handler
func predictHandler(c *gin.Context) {
	var request PredictionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Directly call predictEnergy (no Kafka)
	energy, err := predictEnergy(request.Population, request.Temperature)
	if err != nil {
		log.Printf("Prediction Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prediction"})
		return
	}

	c.JSON(http.StatusOK, PredictionResponse{EnergyKWh: energy})
}

func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)
	fmt.Println("API Server running on port 5000")
	router.Run("0.0.0.0:5000")
}
