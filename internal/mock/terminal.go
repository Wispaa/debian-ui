package mock

import (
	"fmt"
	"strings"
	"time"
)

// ExecuteCommand takes a raw command string and returns a mock output.
func ExecuteCommand(cmd string) string {
	time.Sleep(200 * time.Millisecond) // Simulate a tiny bit of latency

	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}

	parts := strings.Split(cmd, " ")
	base := parts[0]

	switch base {
	case "help":
		return `Available mock commands:
  help      - Show this help message
  ping      - Mock ping a host
  ip        - Show mock network interfaces
  uname     - Show mock system info
  uptime    - Show mock uptime
  clear     - Clear the terminal screen (handled via JS)`
	case "ping":
		target := "8.8.8.8"
		if len(parts) > 1 {
			target = parts[1]
		}
		return fmt.Sprintf(`PING %s (%s) 56(84) bytes of data.
64 bytes from %s: icmp_seq=1 ttl=118 time=14.2 ms
64 bytes from %s: icmp_seq=2 ttl=118 time=12.1 ms
64 bytes from %s: icmp_seq=3 ttl=118 time=13.5 ms

--- %s ping statistics ---
3 packets transmitted, 3 received, 0%% packet loss, time 2003ms
rtt min/avg/max/mdev = 12.124/13.288/14.241/0.884 ms`, target, target, target, target, target, target)
	case "ip":
		if len(parts) > 1 && parts[1] == "a" {
			return `1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 52:54:00:12:34:56 brd ff:ff:ff:ff:ff:ff
    inet 203.0.113.5/24 brd 203.0.113.255 scope global eth0
       valid_lft forever preferred_lft forever`
		}
		return "Usage: ip a"
	case "uname":
		if len(parts) > 1 && parts[1] == "-a" {
			return "Linux debian-router 6.6.0-debian13-amd64 #1 SMP PREEMPT_DYNAMIC x86_64 GNU/Linux"
		}
		return "Linux"
	case "uptime":
		return " 14:32:15 up 12 days,  4:12,  1 user,  load average: 0.05, 0.03, 0.01"
	case "clear":
		// 'clear' is actually intercepted by the frontend to empty the div, but just in case it reaches here:
		return ""
	default:
		return fmt.Sprintf("-bash: %s: command not found", base)
	}
}
