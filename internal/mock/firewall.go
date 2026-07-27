package mock

type Zone struct {
	Name      string
	Active    bool
	OpenPorts []string
}

type FirewallInfo struct {
	Zones []*Zone
}

func GetFirewallInfo() *FirewallInfo {
	return &FirewallInfo{
		Zones: []*Zone{
			{Name: "public", Active: true, OpenPorts: []string{"22/tcp", "80/tcp", "443/tcp"}},
			{Name: "trusted", Active: true, OpenPorts: []string{"Any"}},
		},
	}
}
