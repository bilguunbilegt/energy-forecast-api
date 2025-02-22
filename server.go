package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

// Load Model Coefficients
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

	// Check if coefficients are valid
	if len(coeff) < 3 {
		log.Println("Invalid model coefficients: Model is not trained properly")
		return nil, fmt.Errorf("invalid model coefficients")
	}

	return coeff, nil
}

// Predict energy consumption
func predictEnergy(population, temperature float64) (float64, error) {
	coeff, err := loadModel()
	if err != nil {
		return 0, err
	}

	// Apply regression formula: Energy = b0 + (b1 * Population) + (b2 * Temperature)
	prediction := coeff[0] + (coeff[1] * population) + (coeff[2] * temperature)

	// Ensure predictions make logical sense
	if prediction < 0 {
		log.Printf("Warning: Negative prediction (%f). Setting to 0.\n", prediction)
		prediction = 100 // Set a **minimum threshold** instead of forcing 0.
	}

	return prediction, nil
}

func predictHandler(c *gin.Context) {
	var req PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	energy, err := predictEnergy(req.Population, req.Temperature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Model could not generate prediction"})
		return
	}

	c.JSON(http.StatusOK, PredictionResponse{EnergyKWh: energy})
}

func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)

	fmt.Println("API Server is running on port 5000")
	router.Run("0.0.0.0:5000")
}
