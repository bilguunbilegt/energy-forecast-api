package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Predict energy consumption based on population & temperature
func predictHandler(c *gin.Context) {
	var request PredictionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call predictEnergy function
	energy, err := predictEnergy(request.Population, request.Temperature)
	if err != nil {
		log.Printf("Prediction Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prediction"})
		return
	}

	c.JSON(http.StatusOK, PredictionResponse{EnergyKWh: energy})
}
