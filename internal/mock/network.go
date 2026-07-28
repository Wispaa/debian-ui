package mock

import (
	"sync"
	"time"
)

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

var (
	interfacesConfig = "auto eth0\niface eth0 inet dhcp\n\nauto eth1\niface eth1 inet static\n    address 192.168.1.1\n    netmask 255.255.255.0"
	netMu            sync.Mutex
)

func GetNetworkInterfaces() []*Interface {
	return []*Interface{
		{Name: "eth0", Type: TypePhysical, IP: "203.0.113.5 (WAN)", Status: "Up"},
		{Name: "eth1", Type: TypePhysical, IP: "192.168.1.1 (LAN)", Status: "Up"},
		{Name: "awg0", Type: TypeVirtual, IP: "10.0.0.1", Status: "Down"},
		{Name: "tun0", Type: TypeVirtual, IP: "10.8.0.1", Status: "Up"},
	}
}

func GetNetworkConfig() string {
	netMu.Lock()
	defer netMu.Unlock()
	return interfacesConfig
}

func UpdateNetworkConfig(config string) {
	time.Sleep(1000 * time.Millisecond)
	netMu.Lock()
	defer netMu.Unlock()
	interfacesConfig = config
}
