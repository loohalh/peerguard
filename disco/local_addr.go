package disco

import (
	"log/slog"
	"net"
	"strings"
)

type cidrs []*net.IPNet

func (cs *cidrs) Contains(ip net.IP) bool {
	for _, cidr := range *cs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

type interfaceNamePrefixs []string

func (inp *interfaceNamePrefixs) HasPrefix(interfaceName string) bool {
	for _, ifaceNamePrefix := range *inp {
		if strings.HasPrefix(interfaceName, ifaceNamePrefix) {
			return true
		}
	}
	return false
}

var (
	ignoredLocalCIDRs                cidrs
	ignoredLocalInterfaceNamePrefixs interfaceNamePrefixs
)

func SetIgnoredLocalCIDRs(cidrs ...string) {
	for _, cidr := range cidrs {
		_, _cidr, err := net.ParseCIDR(cidr)
		if err != nil {
			slog.Debug("SetIgnoredLocalCIDRs parse cidr failed", "err", err)
			continue
		}
		ignoredLocalCIDRs = append(ignoredLocalCIDRs, _cidr)
	}
}

func SetIgnoredLocalInterfaceNamePrefixs(prefixs ...string) {
	ignoredLocalInterfaceNamePrefixs = prefixs
}

func ListLocalIPs() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, iface := range ifaces {
		if iface.Flags&net.FlagRunning != net.FlagRunning {
			continue
		}
		if ignoredLocalInterfaceNamePrefixs.HasPrefix(iface.Name) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ignoredLocalCIDRs.Contains(ipnet.IP) {
					continue
				}
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}
