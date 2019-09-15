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
	ambientTemperature := &AmbientTemperature{
		DeviceID:    "selfhydro",
		Temperature: 12.1,
	}
	json.NewEncoder(w).Encode(ambientTemperature)
}
