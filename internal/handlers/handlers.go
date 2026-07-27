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

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
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
