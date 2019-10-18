package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WaterTemperature struct {
	DeviceID    string    `json:"deviceID"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

func getWaterTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	systemID := getParameter(r.URL, "systemid")
	stateRepository := NewStateRepository()
	log.Print("received request for water temperature")
	waterTemperature := stateRepository.GetWaterTemperature(systemID)
	json.NewEncoder(w).Encode(waterTemperature)
}

func getParameter(url *url.URL, key string) string {
	parameter, ok := url.Query()[key]
	if !ok {
		return ""
	}
	fmt.Println(parameter)
	return parameter[0]
}
