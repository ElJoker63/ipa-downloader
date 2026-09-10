package machine

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ElJoker63/ipa-downloader/v2/pkg/util/operatingsystem"
	"golang.org/x/term"
)

//go:generate go run go.uber.org/mock/mockgen -source=machine.go -destination=machine_mock.go -package machine
type Machine interface {
	MacAddress() (string, error)
	HomeDirectory() string
	ReadPassword(fd int) ([]byte, error)
}

type machine struct {
	os operatingsystem.OperatingSystem
}

type Args struct {
	OS operatingsystem.OperatingSystem
}

func New(args Args) Machine {
	return &machine{
		os: args.OS,
	}
}

func (*machine) MacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	if len(interfaces) == 0 {
		return "", fmt.Errorf("could not find network interfaces: %w", err)
	}

	isVirtual := func(name string) bool {
		nameLower := strings.ToLower(name)
		for _, kw := range []string{"radmin", "vethernet", "wsl", "tap", "tun", "pseudo", "vmware", "virtual", "loopback"} {
			if strings.Contains(nameLower, kw) {
				return true
			}
		}
		return false
	}

	// 1st pass: Active physical interface (Up, not loopback, 6-byte MAC, globally unique OUI, not virtual)
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 6 && (iface.HardwareAddr[0]&0x02 == 0) && !isVirtual(iface.Name) {
			return iface.HardwareAddr.String(), nil
		}
	}

	// 2nd pass: Any active interface with a 6-byte MAC that isn't virtual
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 6 && !isVirtual(iface.Name) {
			return iface.HardwareAddr.String(), nil
		}
	}

	// 3rd pass: Any non-virtual interface with a valid MAC
	for _, iface := range interfaces {
		if iface.HardwareAddr.String() != "" && !isVirtual(iface.Name) {
			return iface.HardwareAddr.String(), nil
		}
	}

	// Fallback: any interface with a valid MAC
	for _, iface := range interfaces {
		addr := iface.HardwareAddr.String()
		if addr != "" {
			return addr, nil
		}
	}

	return "", fmt.Errorf("could not find network interfaces with a valid mac address: %w", err)
}

func (m *machine) HomeDirectory() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(m.os.Getenv("HOMEDRIVE"), m.os.Getenv("HOMEPATH"))
	}

	return m.os.Getenv("HOME")
}

func (*machine) ReadPassword(fd int) ([]byte, error) {
	data, err := term.ReadPassword(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to read password: %w", err)
	}

	return data, nil
}
