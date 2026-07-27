package main

import (
	"log"
	"net/http"
	"router-ui/internal/handlers"
	"router-ui/internal/mock"
)

func main() {
	// Initialize mock state
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", h.HandleDashboard)
	mux.HandleFunc("/service/", h.HandleServiceAction)

	log.Println("Starting mock router UI server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
