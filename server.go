package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Structure for incoming JSON request
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// Structure for response
type PredictionResponse struct {
	EnergyKWh float64 `json:"predicted_energy_kwh"`
}

// Load Model Coefficients
func loadModel() ([]float64, error) {
	file, err := os.ReadFile("model.json")
	if err != nil {
		return nil, err
	}

	var coeff []float64
	err = json.Unmarshal(file, &coeff)
	if err != nil {
		return nil, err
	}

	return coeff, nil
}

// Prediction function
func predictEnergy(population, temperature float64) (float64, error) {
	coeff, err := loadModel()
	if err != nil {
		return 0, err
	}

	// Apply regression formula: Energy = b0 + (b1 * Population) + (b2 * Temperature)
	prediction := coeff[0] + (coeff[1] * population) + (coeff[2] * temperature)

	// Ensure non-negative prediction
	return math.Max(prediction, 0), nil
}

func predictHandler(c *gin.Context) {
	var req PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	energy, err := predictEnergy(req.Population, req.Temperature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in prediction"})
		return
	}

	c.JSON(http.StatusOK, PredictionResponse{EnergyKWh: energy})
}

func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)

	fmt.Println("Server running on port 5000")
	router.Run(":5000")
}
