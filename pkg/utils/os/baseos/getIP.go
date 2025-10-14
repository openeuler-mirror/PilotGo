package baseos

import (
	"errors"
	"net"
	"strings"
)

type HostIPInfo struct {
	Interface string
	IP        string
	IsIPv4    bool
	IsPrivate bool
}

func GetAllHostIPs() ([]HostIPInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var results []HostIPInfo
	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil {
				continue
			}
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}

			isV4 := ip.To4() != nil
			isPrivate := isPrivateIP(ip)

			results = append(results, HostIPInfo{
				Interface: iface.Name,
				IP:        ip.String(),
				IsIPv4:    isV4,
				IsPrivate: isPrivate,
			})
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no valid non-loopback addresses found")
	}
	return results, nil
}

func extractIP(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		s := addr.String()
		if strings.Contains(s, "/") {
			parts := strings.SplitN(s, "/", 2)
			ip := net.ParseIP(parts[0])
			return ip
		}
		return net.ParseIP(s)
	}
}

func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		if v4[0] == 10 {
			return true
		}

		if v4[0] == 192 && v4[1] == 168 {
			return true
		}
		if v4[0] == 172 && (v4[1] >= 16 && v4[1] <= 31) {
			return true
		}
		if v4[0] == 169 && v4[1] == 254 {
			return true
		}
		return false
	}

	if len(ip) == net.IPv6len {
		b := ip[0]
		if (b & 0xfe) == 0xfc {
			return true
		}
	}
	return false
}
