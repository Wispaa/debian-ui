package mock

import (
	"strings"
	"sync"
	"time"
)

type DHCPLease struct {
	IPAddress  string
	IPSafe     string // Safe for HTML IDs
	MACAddress string
	Hostname   string
	ExpiresAt  string
	Type       string // "Dynamic" or "Static"
}

var (
	dhcpMu     sync.Mutex
	mockLeases = map[string]*DHCPLease{
		"192.168.1.100": {IPAddress: "192.168.1.100", MACAddress: "aa:bb:cc:dd:ee:01", Hostname: "Jules-Laptop", ExpiresAt: time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05"), Type: "Dynamic"},
		"192.168.1.101": {IPAddress: "192.168.1.101", MACAddress: "11:22:33:44:55:66", Hostname: "Smart-TV", ExpiresAt: time.Now().Add(12 * time.Hour).Format("2006-01-02 15:04:05"), Type: "Dynamic"},
		"192.168.1.50":  {IPAddress: "192.168.1.50", MACAddress: "ff:ee:dd:cc:bb:aa", Hostname: "NAS-Server", ExpiresAt: "Never", Type: "Static"},
		"192.168.1.105": {IPAddress: "192.168.1.105", MACAddress: "00:1A:2B:3C:4D:5E", Hostname: "Android-Phone", ExpiresAt: time.Now().Add(45 * time.Minute).Format("2006-01-02 15:04:05"), Type: "Dynamic"},
	}
)

func GetDHCPLeases(searchQuery string) []*DHCPLease {
	dhcpMu.Lock()
	defer dhcpMu.Unlock()

	var result []*DHCPLease
	for _, lease := range mockLeases {
		if searchQuery != "" {
			query := strings.ToLower(searchQuery)
			if !strings.Contains(strings.ToLower(lease.IPAddress), query) &&
				!strings.Contains(strings.ToLower(lease.MACAddress), query) &&
				!strings.Contains(strings.ToLower(lease.Hostname), query) {
				continue
			}
		}
		// Create a copy to prevent race conditions during read
		result = append(result, &DHCPLease{
			IPAddress:  lease.IPAddress,
			IPSafe:     strings.ReplaceAll(lease.IPAddress, ".", "-"),
			MACAddress: lease.MACAddress,
			Hostname:   lease.Hostname,
			ExpiresAt:  lease.ExpiresAt,
			Type:       lease.Type,
		})
	}
	return result
}

func ActionDHCPLease(ip, action string) *DHCPLease {
	time.Sleep(500 * time.Millisecond) // Simulate a tiny delay for processing
	dhcpMu.Lock()
	defer dhcpMu.Unlock()

	lease, ok := mockLeases[ip]
	if !ok {
		return nil
	}

	switch action {
	case "static":
		lease.Type = "Static"
		lease.ExpiresAt = "Never"
	case "revoke":
		delete(mockLeases, ip)
		return nil // Indicate deletion
	}

	return &DHCPLease{
		IPAddress:  lease.IPAddress,
		IPSafe:     strings.ReplaceAll(lease.IPAddress, ".", "-"),
		MACAddress: lease.MACAddress,
		Hostname:   lease.Hostname,
		ExpiresAt:  lease.ExpiresAt,
		Type:       lease.Type,
	}
}
