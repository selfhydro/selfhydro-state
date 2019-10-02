package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type AmbientTemperature struct {
	DeviceID    string  `json:"deviceID"`
	Temperature float64 `json:"temperature"`
}

func getTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	stateRepository := NewStateRepository()
	log.Print("received request for ambient temperature")
	ambientTemperature := stateRepository.GetAmbientTemperature("selfhydro")
	json.NewEncoder(w).Encode(ambientTemperature)
}
