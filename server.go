package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bilguunbilegt/energy-forecast-api/training" // Ensure correct module import
	"github.com/gin-gonic/gin"
)

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

// Health Check Handler
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	// Train the model at startup
	log.Println("Starting model training...")
	training.TrainModel()
	log.Println("Model training completed.")

	router := gin.Default()

	// Register API routes
	router.POST("/predict", predictHandler) // Uses predictHandler from predictHandler.go
	router.GET("/health", healthCheck)      // Health check endpoint

	fmt.Println("API Server is running on port 5000")
	log.Fatal(router.Run("0.0.0.0:5000")) // logs errors
}
