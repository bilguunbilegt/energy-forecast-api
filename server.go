package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

// Load Model Coefficients from model.json
func loadModel() ([]float64, error) {
	file, err := os.ReadFile("model.json")
	if err != nil {
		log.Printf("Error loading model file: %v", err)
		return nil, err
	}

	var coeff []float64
	err = json.Unmarshal(file, &coeff)
	if err != nil {
		log.Printf("Error parsing model file: %v", err)
		return nil, err
	}

	if len(coeff) < 3 {
		log.Println("Invalid model coefficients: Model is not trained properly")
		return nil, fmt.Errorf("invalid model coefficients")
	}

	return coeff, nil
}

// Compute Energy Consumption Prediction
func predictEnergy(population, temperature float64) (float64, error) {
	coeff, err := loadModel()
	if err != nil {
		return 0, err
	}

	// Energy = b0 + (b1 * Population) + (b2 * Temperature)
	prediction := coeff[0] + (coeff[1] * population) + (coeff[2] * temperature)

	if prediction < 0 {
		log.Printf("Warning: Negative prediction (%f). Setting to 100.\n", prediction)
		prediction = 100
	}

	return prediction, nil
}

// API Endpoint for Predictions
func predictHandler(c *gin.Context) {
	var req PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	energy, err := predictEnergy(req.Population, req.Temperature)
	if err != nil {
		log.Printf("Prediction Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate prediction"})
		return
	}

	c.JSON(http.StatusOK, PredictionResponse{EnergyKWh: energy})
}

// Main Function to Start API
func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)

	fmt.Println("API Server is running on port 5000")
	router.Run("0.0.0.0:5000")
}
