package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	// Diagnostics
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})

	// Participants
	mux.HandleFunc("GET /v1/participants", mockListHandler)
	mux.HandleFunc("GET /v1/participants/{id}", mockGetHandler)

	// Settlement Intents
	mux.HandleFunc("POST /v1/settlement-intents", mockCreateHandler)
	mux.HandleFunc("GET /v1/settlement-intents/{id}", mockGetHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/accept", mockSuccessHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/reject", mockSuccessHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/reserve", mockSuccessHandler)
	mux.HandleFunc("POST /v1/settlement-intents/{id}/execute", mockSuccessHandler)

	// Reservations
	mux.HandleFunc("GET /v1/reservations/{id}", mockGetHandler)

	// Settlement Instructions
	mux.HandleFunc("GET /v1/settlement-instructions/{id}", mockGetHandler)

	// Settlement Proofs
	mux.HandleFunc("GET /v1/settlement-proofs/{id}", mockGetHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting gateway on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func mockSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func mockGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": "mock-id", "status": "FINAL"})
}

func mockListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]string{{"id": "mock-id"}})
}

func mockCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": "mock-id"})
}
