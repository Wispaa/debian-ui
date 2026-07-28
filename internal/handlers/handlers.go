package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"router-ui/internal/mock"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

type DashboardData struct {
	Services   []*mock.Service
	Interfaces []*mock.Interface
	Firewall   *mock.FirewallInfo
}

type ConfigData struct {
	Title   string
	Status  string
	Config  string
	SaveURL string
}

type LogsData struct {
	Logs []mock.LogEntry
}

type Handler struct {
	sm    *mock.ServiceManager
	tmpl  *template.Template
}

func NewHandler(sm *mock.ServiceManager) *Handler {
	// Parse templates from embedded FS
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	return &Handler{
		sm:   sm,
		tmpl: tmpl,
	}
}

func (h *Handler) HandleLayout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := DashboardData{
		Services:   h.sm.GetAll(),
		Interfaces: mock.GetNetworkInterfaces(),
		Firewall:   mock.GetFirewallInfo(),
	}

	err := h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	data := DashboardData{
		Services:   h.sm.GetAll(),
		Interfaces: mock.GetNetworkInterfaces(),
		Firewall:   mock.GetFirewallInfo(),
	}

	err := h.tmpl.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/config/")

	var data ConfigData

	if name == "interfaces" {
		data = ConfigData{
			Title:   "Network Interfaces",
			Config:  mock.GetNetworkConfig(),
			SaveURL: "/config/interfaces/save",
		}
	} else {
		svc := h.sm.Get(name)
		if svc == nil {
			http.NotFound(w, r)
			return
		}
		data = ConfigData{
			Title:   svc.DisplayName,
			Status:  string(svc.Status),
			Config:  svc.Config,
			SaveURL: "/config/" + svc.Name + "/save",
		}
	}

	err := h.tmpl.ExecuteTemplate(w, "config_editor", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleConfigSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	configStr := r.Form.Get("config")
	// The path is /config/{name}/save
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	name := parts[2]

	if name == "interfaces" {
		mock.UpdateNetworkConfig(configStr)
	} else {
		if h.sm.UpdateConfig(name, configStr) == nil {
			http.NotFound(w, r)
			return
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "save_success", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleLogsPage(w http.ResponseWriter, r *http.Request) {
	// Generate initial logs without filters
	logs := mock.GenerateMockLogs("", "", "")

	data := LogsData{
		Logs: logs,
	}

	err := h.tmpl.ExecuteTemplate(w, "logs_page", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleMockLogs(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	presetDate := r.URL.Query().Get("preset_date")
	customDate := r.URL.Query().Get("custom_date")
	customTime := r.URL.Query().Get("custom_time")
	search := r.URL.Query().Get("search")

	// Determine what date string to pass to the mock (simplified for mockup)
	dateFilter := presetDate
	if presetDate == "Custom Date" {
		dateFilter = customDate + " " + customTime
	}

	// Generate mock logs based on filters
	logs := mock.GenerateMockLogs(service, dateFilter, search)

	data := LogsData{
		Logs: logs,
	}

	err := h.tmpl.ExecuteTemplate(w, "log_entries", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleServiceAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Path should be /service/{name}/{action}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/service/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	name := parts[0]
	action := parts[1]

	updatedService := h.sm.Action(name, action)
	if updatedService == nil {
		http.NotFound(w, r)
		return
	}

	// Render only the updated service_card fragment for HTMX to swap
	err := h.tmpl.ExecuteTemplate(w, "service_card", updatedService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
