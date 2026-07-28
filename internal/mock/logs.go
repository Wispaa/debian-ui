package mock

import (
	"fmt"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp string
	Service   string
	Level     string // INFO, WARN, ERROR
	Message   string
}

// GenerateMockLogs generates a batch of mock logs for the given parameters.
func GenerateMockLogs(serviceFilter, dateFilter, searchFilter string) []LogEntry {
	var logs []LogEntry

	services := []string{"sing-box", "amnezia-wg", "wireguard", "openvpn", "dnscrypt-proxy", "firewalld", "systemd"}

	// If a specific service is requested from the dropdown, only generate logs for it
	if serviceFilter != "" && serviceFilter != "All Services" {
		services = []string{serviceFilter}
	}

	levels := []string{"INFO", "INFO", "INFO", "WARN", "ERROR"}

	now := time.Now()

	// Generate some fake logs
	for i := 0; i < 50; i++ {
		// Go back in time a bit for each log
		t := now.Add(-time.Duration(50-i) * time.Minute)

		svc := services[i%len(services)]
		level := levels[i%len(levels)]

		msg := fmt.Sprintf("Mock log message for %s at sequence %d", svc, i)
		if level == "WARN" {
			msg = fmt.Sprintf("Warning: something might be wrong with %s", svc)
		} else if level == "ERROR" {
			msg = fmt.Sprintf("Error: failed to connect or start %s", svc)
		}

		entry := LogEntry{
			Timestamp: t.Format("Jan 02 15:04:05"),
			Service:   svc,
			Level:     level,
			Message:   msg,
		}

		// Apply search filter if present
		if searchFilter != "" {
			if !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(searchFilter)) &&
			   !strings.Contains(strings.ToLower(entry.Service), strings.ToLower(searchFilter)) {
				continue // Skip if it doesn't match the search query
			}
		}

		logs = append(logs, entry)
	}

	return logs
}
