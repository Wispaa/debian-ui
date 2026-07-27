package mock

type InterfaceType string

const (
	TypePhysical InterfaceType = "Physical"
	TypeVirtual  InterfaceType = "Virtual"
)

type Interface struct {
	Name   string
	Type   InterfaceType
	IP     string
	Status string // "Up", "Down"
}

func GetNetworkInterfaces() []*Interface {
	return []*Interface{
		{Name: "eth0", Type: TypePhysical, IP: "203.0.113.5 (WAN)", Status: "Up"},
		{Name: "eth1", Type: TypePhysical, IP: "192.168.1.1 (LAN)", Status: "Up"},
		{Name: "awg0", Type: TypeVirtual, IP: "10.0.0.1", Status: "Down"},
		{Name: "tun0", Type: TypeVirtual, IP: "10.8.0.1", Status: "Up"},
	}
}
