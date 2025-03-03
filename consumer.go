package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Define Kafka topic and broker
const (
	kafkaBroker = "localhost:9092"
	kafkaTopic  = "energy-predictions"
)

// Struct for prediction request
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// Kafka Consumer Function
func consumeFromKafka() {
	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: "energy-consumer",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		// Parse request
		var request PredictionRequest
		if err := json.Unmarshal(msg.Value, &request); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		// Call prediction model (Replace this with actual logic)
		fmt.Printf("Processing Prediction: Population: %.0f, Temp: %.1f\n",
			request.Population, request.Temperature)
	}

	reader.Close()
}

func main() {
	fmt.Println("Kafka Consumer is running...")
	consumeFromKafka()
}
