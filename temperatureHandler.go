package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type WaterTemperature struct {
	DeviceID    string  `json:"deviceID"`
	Temperature float64 `json:"temperature"`
}

func getTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	stateRepository := NewStateRepository()
	log.Print("received request for water temperature")
	waterTemperature := stateRepository.GetWaterTemperature("selfhydro")
	json.NewEncoder(w).Encode(waterTemperature)
}
