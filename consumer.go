package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "kafka:9092"  // ✅ Fixed to work inside Docker Compose
	kafkaTopic  = "energy-predictions"
)

type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

func consumeFromKafka() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBroker},
		Topic:       kafkaTopic,
		GroupID:     "energy-consumer",
		StartOffset: kafka.FirstOffset,  // ✅ Ensures old messages are read
	})

	fmt.Println("🟢 Kafka Consumer is running...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("❌ Error reading message:", err)
			continue
		}

		var request PredictionRequest
		if err := json.Unmarshal(msg.Value, &request); err != nil {
			log.Println("❌ Error decoding JSON:", err)
			continue
		}

		fmt.Printf("📩 Received Prediction: Population: %.0f, Temp: %.1f\n",
			request.Population, request.Temperature)
	}

	reader.Close()
}

func main() {
	consumeFromKafka()
}
