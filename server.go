package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// Request structure
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// SageMaker Response Structure
type SageMakerResponse struct {
	PredictedEnergy float64 `json:"predictions"`
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

	if len(coeff) < 3 {
		log.Println("Invalid model coefficients: Model is not trained properly")
		return nil, fmt.Errorf("invalid model coefficients")
	}

	return coeff, nil
}

// Send Data to SageMaker for Prediction
func predictWithSageMaker(population, temperature float64) (float64, error) {
	sagemakerURL := "https://runtime.sagemaker.us-east-1.amazonaws.com/endpoints/your-model-endpoint/invocations"

	payload := map[string]float64{
		"population":  population,
		"temperature": temperature,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", sagemakerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Amz-Target", "SageMaker.InvokeEndpoint")
	req.Header.Set("Authorization", "Bearer YOUR_AWS_AUTH_TOKEN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var sagemakerResponse SageMakerResponse
	err = json.Unmarshal(body, &sagemakerResponse)
	if err != nil {
		return 0, err
	}

	return sagemakerResponse.PredictedEnergy, nil
}

func predictHandler(c *gin.Context) {
	var req PredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	energy, err := predictWithSageMaker(req.Population, req.Temperature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SageMaker could not generate prediction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"predicted_energy_kwh": energy})
}

func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)

	fmt.Println("API Server is running on port 5000")
	router.Run("0.0.0.0:5000")
}
