package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sajari/regression"
)

// Data structure
type EnergyData struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
	EnergyKWh   float64 `json:"energy_kwh"`
}

// Training data in JSON format
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

	// Create a new regression model
	var r regression.Regression
	r.SetObserved("EnergyKWh")
	r.SetVar(0, "Population")
	r.SetVar(1, "Temperature")

	// Add training data
	for _, d := range energyData {
		r.Train(regression.DataPoint(d.EnergyKWh, []float64{d.Population, d.Temperature}))
	}

	// Train the model
	r.Run()

	// Save model coefficients to a file
	file, err := os.Create("model.json")
	if err != nil {
		log.Fatalf("Failed to save model: %v", err)
	}
	defer file.Close()

	modelData, _ := json.Marshal(r.Coeff)
	file.Write(modelData)

	fmt.Println("Model trained and saved to model.json")
}
