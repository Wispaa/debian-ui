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
	Config      string
}

type ServiceManager struct {
	mu       sync.Mutex
	services map[string]*Service
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: map[string]*Service{
			"sing-box":       {Name: "sing-box", DisplayName: "sing-box", Status: StatusRunning, Config: "{\n  \"log\": {\n    \"level\": \"info\"\n  },\n  \"inbounds\": []\n}"},
			"amnezia-wg":     {Name: "amnezia-wg", DisplayName: "Amnezia WG", Status: StatusStopped, Config: "[Interface]\nPrivateKey = ...\nAddress = 10.0.0.1/24"},
			"wireguard":      {Name: "wireguard", DisplayName: "WireGuard", Status: StatusRunning, Config: "[Interface]\nPrivateKey = ...\nListenPort = 51820"},
			"openvpn":        {Name: "openvpn", DisplayName: "OpenVPN", Status: StatusStopped, Config: "client\ndev tun\nproto udp\nremote 1.2.3.4 1194"},
			"dnscrypt-proxy": {Name: "dnscrypt-proxy", DisplayName: "DNSCrypt Proxy", Status: StatusRunning, Config: "server_names = ['cloudflare', 'google']\nlisten_addresses = ['127.0.0.1:53']"},
			"firewalld":      {Name: "firewalld", DisplayName: "Firewalld", Status: StatusRunning, Config: "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<zone>\n  <short>public</short>\n</zone>"},
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
			list = append(list, &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status, Config: s.Config})
		}
	}
	return list
}

func (sm *ServiceManager) Get(name string) *Service {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.services[name]; ok {
		return &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status, Config: s.Config}
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

	return &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status, Config: s.Config}
}

func (sm *ServiceManager) UpdateConfig(name string, config string) *Service {
	// Simulate OS file write delay
	time.Sleep(1000 * time.Millisecond)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	s, ok := sm.services[name]
	if !ok {
		return nil
	}

	s.Config = config
	return &Service{Name: s.Name, DisplayName: s.DisplayName, Status: s.Status, Config: s.Config}
}
