package mock

import (
	"strings"
)

type SystemUnit struct {
	UnitName    string
	LoadState   string
	ActiveState string
	SubState    string
	Description string
}

func GetSystemUnits(searchQuery string, stateFilter string) []SystemUnit {
	allUnits := []SystemUnit{
		{"sing-box.service", "loaded", "active", "running", "sing-box routing platform"},
		{"amnezia-wg.service", "loaded", "inactive", "dead", "Amnezia WireGuard tunnel"},
		{"wireguard.service", "loaded", "active", "running", "WireGuard VPN service"},
		{"openvpn.service", "loaded", "failed", "failed", "OpenVPN connection daemon"},
		{"dnscrypt-proxy.service", "loaded", "active", "running", "DNSCrypt client proxy"},
		{"firewalld.service", "loaded", "active", "running", "firewalld - dynamic firewall daemon"},
		{"sshd.service", "loaded", "active", "running", "OpenSSH server daemon"},
		{"nginx.service", "loaded", "active", "running", "A high performance web server and a reverse proxy server"},
		{"cron.service", "loaded", "active", "running", "Regular background program processing daemon"},
		{"rsyslog.service", "loaded", "active", "running", "System Logging Service"},
		{"dbus.service", "loaded", "active", "running", "D-Bus System Message Bus"},
		{"systemd-journald.service", "loaded", "active", "running", "Journal Service"},
		{"networkd-dispatcher.service", "loaded", "active", "running", "Dispatcher daemon for systemd-networkd"},
		{"unattended-upgrades.service", "loaded", "inactive", "dead", "Unattended Upgrades Shutdown"},
		{"bluetooth.service", "loaded", "inactive", "dead", "Bluetooth service"},
	}

	var filtered []SystemUnit

	for _, unit := range allUnits {
		// Apply State Filter
		if stateFilter != "" && stateFilter != "All" {
			if !strings.EqualFold(unit.ActiveState, stateFilter) {
				continue
			}
		}

		// Apply Search Filter
		if searchQuery != "" {
			query := strings.ToLower(searchQuery)
			if !strings.Contains(strings.ToLower(unit.UnitName), query) &&
				!strings.Contains(strings.ToLower(unit.Description), query) {
				continue
			}
		}

		filtered = append(filtered, unit)
	}

	return filtered
}
