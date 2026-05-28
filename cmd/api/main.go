package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	port := ":7070"
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// healthHandler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // 1.header
	w.WriteHeader(http.StatusOK)                       // 2.status
	data := map[string]string{
		"status": "ok, health from Go Server ",
	}
	// w.Write([]byte(`{"status": "ok"}`)) // 3.body
	json.NewEncoder(w).Encode(data)
}
