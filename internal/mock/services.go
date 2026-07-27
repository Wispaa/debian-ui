package mock

import (
	"sync"
	"time"
)

type ServiceStatus string

const (
	StatusRunning ServiceStatus = "Running"
	StatusStopped ServiceStatus = "Stopped"
)

type Service struct {
	Name        string
	DisplayName string
	Status      ServiceStatus
}

type ServiceManager struct {
	mu       sync.Mutex
	services map[string]*Service
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: map[string]*Service{
			"sing-box":       {Name: "sing-box", DisplayName: "sing-box", Status: StatusRunning},
			"amnezia-wg":     {Name: "amnezia-wg", DisplayName: "Amnezia WG", Status: StatusStopped},
			"wireguard":      {Name: "wireguard", DisplayName: "WireGuard", Status: StatusRunning},
			"openvpn":        {Name: "openvpn", DisplayName: "OpenVPN", Status: StatusStopped},
			"dnscrypt-proxy": {Name: "dnscrypt-proxy", DisplayName: "DNSCrypt Proxy", Status: StatusRunning},
			"firewalld":      {Name: "firewalld", DisplayName: "Firewalld", Status: StatusRunning},
		},
	}
}

func (sm *ServiceManager) GetAll() []*Service {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var list []*Service
	order := []string{"sing-box", "amnezia-wg", "wireguard", "openvpn", "dnscrypt-proxy", "firewalld"}
	for _, name := range order {
		if s, ok := sm.services[name]; ok {
			list = append(list, &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status})
		}
	}
	return list
}

func (sm *ServiceManager) Get(name string) *Service {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.services[name]; ok {
		return &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status}
	}
	return nil
}

func (sm *ServiceManager) Action(name string, action string) *Service {
	// Simulate OS execution delay
	time.Sleep(1500 * time.Millisecond)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	s, ok := sm.services[name]
	if !ok {
		return nil
	}

	switch action {
	case "start":
		s.Status = StatusRunning
	case "stop":
		s.Status = StatusStopped
	case "restart":
		s.Status = StatusRunning
	}

	return &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status}
}
