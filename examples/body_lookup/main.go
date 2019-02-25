package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type estimationRequest struct {
	Amount   float64
	Currency string
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to weaver!")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func getAmount() float64 {
	maxCap, _ := strconv.ParseFloat(getEnv("MAX_AMOUNT", "100"), 64)
	minCap, _ := strconv.ParseFloat(getEnv("MIN_AMOUNT", "100"), 64)
	randomValue := minCap + rand.Float64()*(maxCap-minCap)
	return math.Round(randomValue*100) / 100
}

func handleEstimate(w http.ResponseWriter, r *http.Request) {
	estimatedValue := estimationRequest{Amount: getAmount(), Currency: getEnv("CURRENCY", "IDR")}
	respEncoder := json.NewEncoder(w)
	respEncoder.Encode(estimatedValue)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/estimate", handleEstimate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
