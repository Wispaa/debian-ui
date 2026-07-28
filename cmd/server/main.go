package main

import (
	"log"
	"net/http"
	"strings"
	"router-ui/internal/handlers"
	"router-ui/internal/mock"
)

func main() {
	// Initialize mock state
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", h.HandleLayout)
	mux.HandleFunc("/page/dashboard", h.HandleDashboard)

	// Config routes
	mux.HandleFunc("/config/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/save") {
			h.HandleConfigSave(w, r)
		} else {
			h.HandleConfig(w, r)
		}
	})

	// Logs routes
	mux.HandleFunc("/page/logs", h.HandleLogsPage)
	mux.HandleFunc("/api/mock/logs", h.HandleMockLogs)

	// System routes
	mux.HandleFunc("/page/system", h.HandleSystemPage)
	mux.HandleFunc("/api/mock/system", h.HandleMockSystem)

	// APT routes
	mux.HandleFunc("/page/apt", h.HandleAptPage)
	mux.HandleFunc("/api/mock/apt/search", h.HandleMockAptSearch)
	mux.HandleFunc("/api/mock/apt/action", h.HandleMockAptAction)
	mux.HandleFunc("/api/mock/apt/upload", h.HandleMockAptUpload)
	mux.HandleFunc("/api/mock/apt/upgrade", h.HandleMockAptUpgrade)

	mux.HandleFunc("/service/", h.HandleServiceAction)

	log.Println("Starting mock router UI server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
