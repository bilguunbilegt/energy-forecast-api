package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// Define Kafka topic and broker
const (
	kafkaBroker = "localhost:9092" // Change to your EC2 or AWS MSK address
	kafkaTopic  = "energy-predictions"
)

// Struct for prediction request
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// Kafka Producer Function
func sendToKafka(data PredictionRequest) error {
	// Create a Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
	})

	// Encode message
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Send message
	err = writer.WriteMessages(context.Background(),
		kafka.Message{Value: msg},
	)

	writer.Close()
	return err
}

// Prediction API Handler
func predictHandler(c *gin.Context) {
	var request PredictionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Send data to Kafka
	err := sendToKafka(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prediction request sent to Kafka"})
}

func main() {
	router := gin.Default()
	router.POST("/predict", predictHandler)
	fmt.Println("API Server running on 5000")
	router.Run("0.0.0.0:5000")
}
