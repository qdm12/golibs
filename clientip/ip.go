package clientip

import "net"

func (e *extractor) ipIsPrivate(ip net.IP) bool {
	for _, ipNet := range e.privateIPNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func (e *extractor) extractPublicIPs(ips []net.IP) (publicIPs []net.IP) {
	for _, ip := range ips {
		if e.ipIsPrivate(ip) {
			continue
		}
		publicIPs = append(publicIPs, ip)
	}
	return publicIPs
}

func parseIPs(stringIPs []string) (ips []net.IP) {
	for _, s := range stringIPs {
		ip := net.ParseIP(s)
		if ip != nil {
			ips = append(ips, ip)
		}
	}
	return ips
}
