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
	handler := http.HandlerFunc(h.HandleLayout)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Admin Dashboard") {
		t.Errorf("handler returned unexpected body, does not contain 'Admin Dashboard'")
	}
	if !strings.Contains(body, "sing-box") {
		t.Errorf("handler returned unexpected body, does not contain 'sing-box'")
	}
}

func TestHandleLogs(t *testing.T) {
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	// Test GET /page/logs
	req, err := http.NewRequest("GET", "/page/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HandleLogsPage)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "System Logs") {
		t.Errorf("handler returned unexpected body, does not contain 'System Logs'")
	}
	if !strings.Contains(body, "service-filter") {
		t.Errorf("handler returned unexpected body, does not contain 'service-filter'")
	}

	// Test GET /api/mock/logs with a filter
	req2, err := http.NewRequest("GET", "/api/mock/logs?service=sing-box", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(h.HandleMockLogs)
	handler2.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body2 := rr2.Body.String()
	// Since we filtered by sing-box, we shouldn't see amnezia-wg
	if strings.Contains(body2, "amnezia-wg") {
		t.Errorf("handler returned unexpected body, contains 'amnezia-wg' but should only have 'sing-box'")
	}
}

func TestHandleConfig(t *testing.T) {
	sm := mock.NewServiceManager()
	h := handlers.NewHandler(sm)

	req, err := http.NewRequest("GET", "/config/sing-box", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HandleConfig)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "sing-box Configuration") {
		t.Errorf("handler returned unexpected body, does not contain 'sing-box Configuration'")
	}
	if !strings.Contains(body, "<textarea") {
		t.Errorf("handler returned unexpected body, does not contain textarea")
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
