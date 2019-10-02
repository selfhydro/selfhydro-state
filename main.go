package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/ambientTemperature", getTemperatureHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"https://selfhydro.com"},
	}).Handler(r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
