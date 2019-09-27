package main

import (
	"encoding/json"
	"net/http"
)

type AmbientTemperature struct {
	DeviceID    string  `json:"deviceID"`
	Temperature float64 `json:"temperature"`
}

func getTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	stateRepository := NewStateRepository()
	ambientTemperature := stateRepository.GetAmbientTemperature("selfhydro")
	json.NewEncoder(w).Encode(ambientTemperature)
}
