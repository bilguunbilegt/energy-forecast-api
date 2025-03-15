package main

// Request structure for prediction
type PredictionRequest struct {
	Population  float64 `json:"population"`
	Temperature float64 `json:"temperature"`
}

// Response structure for prediction
type PredictionResponse struct {
	EnergyKWh float64 `json:"predicted_energy_kwh"`
}
