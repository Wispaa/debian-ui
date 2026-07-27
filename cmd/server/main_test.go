package main

import (
	"net/http"
	"net/http/httptest"
	"router-ui/internal/handlers"
	"router-ui/internal/mock"
	"strings"
	"testing"
)

func TestHandleDashboard(t *testing.T) {
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HandleDashboard)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Router Dashboard") {
		t.Errorf("handler returned unexpected body, does not contain 'Router Dashboard'")
	}
	if !strings.Contains(body, "sing-box") {
		t.Errorf("handler returned unexpected body, does not contain 'sing-box'")
	}
}

func TestHandleServiceAction(t *testing.T) {
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	// Test starting a stopped service (amnezia-wg is initially stopped)
	req, err := http.NewRequest("POST", "/service/amnezia-wg/start", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HandleServiceAction)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	// Should render the partial fragment for amnezia-wg as Running
	if !strings.Contains(body, `id="service-amnezia-wg"`) {
		t.Errorf("expected fragment to contain id='service-amnezia-wg'")
	}
	if !strings.Contains(body, "Running") {
		t.Errorf("expected fragment to indicate the service is 'Running'")
	}

	// Also verify the internal state was updated
	svc := sm.Get("amnezia-wg")
	if svc.Status != mock.StatusRunning {
		t.Errorf("expected mock state to be Running, got %v", svc.Status)
	}
}
