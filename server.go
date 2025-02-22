package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sajari/regression"
)

// Data structure for training
type EnergyData struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
	EnergyKWh   float64 `json:"energy_kwh"`
}

// Sample training data
var jsonData = `[
	{"population": 12670000, "temperature": 30, "energy_kwh": 2535},
	{"population": 12670000, "temperature": 33, "energy_kwh": 977},
	{"population": 12670000, "temperature": 45, "energy_kwh": 1754},
	{"population": 12796700, "temperature": 30, "energy_kwh": 594}
]`

func main() {
	var energyData []EnergyData
	err := json.Unmarshal([]byte(jsonData), &energyData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Create regression model
	var r regression.Regression
	r.SetObserved("EnergyKWh")
	r.SetVar(0, "Population")
	r.SetVar(1, "Temperature")

	// Train model
	for _, d := range energyData {
		r.Train(regression.DataPoint(d.EnergyKWh, []float64{d.Population, d.Temperature}))
	}

	// Run regression
	r.Run()

	// Extract numeric coefficients
	coefficients := make([]float64, len(r.Coefficients))
	for i, coeff := range r.Coefficients {
		coefficients[i] = coeff.Value
	}

	// Print model coefficients (for debugging)
	fmt.Println("Model Coefficients:", coefficients)

	// Save model coefficients
	modelFile, err := os.Create("model.json")
	if err != nil {
		log.Fatalf("Failed to save model: %v", err)
	}
	defer modelFile.Close()

	// Write JSON data properly
	modelData, err := json.MarshalIndent(coefficients, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode model data: %v", err)
	}

	_, err = modelFile.Write(modelData)
	if err != nil {
		log.Fatalf("Failed to write to model.json: %v", err)
	}

	fmt.Println("Model trained and saved as model.json")
}
