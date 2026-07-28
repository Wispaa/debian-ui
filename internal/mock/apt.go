package mock

import (
	"strings"
	"sync"
	"time"
)

type PackageStatus string

const (
	PkgInstalled       PackageStatus = "Installed"
	PkgUpdateAvailable PackageStatus = "Update Available"
	PkgAvailable       PackageStatus = "Available"
)

type Package struct {
	Name             string
	InstalledVersion string
	LatestVersion    string
	Status           PackageStatus
}

var (
	aptMu       sync.Mutex
	mockPackages = map[string]*Package{
		"curl":         {Name: "curl", InstalledVersion: "7.88.1-10", LatestVersion: "7.88.1-10", Status: PkgInstalled},
		"sing-box":     {Name: "sing-box", InstalledVersion: "1.8.0", LatestVersion: "1.9.0-rc", Status: PkgUpdateAvailable},
		"amnezia-wg":   {Name: "amnezia-wg", InstalledVersion: "1.0.2", LatestVersion: "1.1.0", Status: PkgUpdateAvailable},
		"htop":         {Name: "htop", InstalledVersion: "3.2.2-2", LatestVersion: "3.2.2-2", Status: PkgInstalled},
		"vim":          {Name: "vim", InstalledVersion: "", LatestVersion: "2:9.0.1378-2", Status: PkgAvailable},
		"tcpdump":      {Name: "tcpdump", InstalledVersion: "", LatestVersion: "4.99.3-1", Status: PkgAvailable},
		"nginx":        {Name: "nginx", InstalledVersion: "1.22.1-9", LatestVersion: "1.24.0", Status: PkgUpdateAvailable},
		"python3":      {Name: "python3", InstalledVersion: "3.11.2-1", LatestVersion: "3.11.2-1", Status: PkgInstalled},
	}
)

func GetPackages(searchQuery string) []*Package {
	aptMu.Lock()
	defer aptMu.Unlock()

	var result []*Package
	for _, pkg := range mockPackages {
		if searchQuery != "" {
			if !strings.Contains(strings.ToLower(pkg.Name), strings.ToLower(searchQuery)) {
				continue
			}
		}
		// Create a copy to prevent race conditions during read
		result = append(result, &Package{
			Name:             pkg.Name,
			InstalledVersion: pkg.InstalledVersion,
			LatestVersion:    pkg.LatestVersion,
			Status:           pkg.Status,
		})
	}
	return result
}

func ActionPackage(name, action string) *Package {
	time.Sleep(1000 * time.Millisecond) // Simulate dpkg delay
	aptMu.Lock()
	defer aptMu.Unlock()

	pkg, ok := mockPackages[name]
	if !ok {
		return nil
	}

	switch action {
	case "install":
		pkg.InstalledVersion = pkg.LatestVersion
		pkg.Status = PkgInstalled
	case "upgrade":
		pkg.InstalledVersion = pkg.LatestVersion
		pkg.Status = PkgInstalled
	case "remove":
		pkg.InstalledVersion = ""
		pkg.Status = PkgAvailable
	}

	return &Package{
		Name:             pkg.Name,
		InstalledVersion: pkg.InstalledVersion,
		LatestVersion:    pkg.LatestVersion,
		Status:           pkg.Status,
	}
}

// Simulate running apt update && apt upgrade
func SimulateFullUpgrade() {
	time.Sleep(2 * time.Second)
	aptMu.Lock()
	defer aptMu.Unlock()

	for _, pkg := range mockPackages {
		if pkg.Status == PkgUpdateAvailable {
			pkg.InstalledVersion = pkg.LatestVersion
			pkg.Status = PkgInstalled
		}
	}
}
