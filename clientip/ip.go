package clientip

import "net"

func (p *Parser) ipIsPrivate(ip net.IP) bool {
	for _, ipNet := range p.privateIPNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func (p *Parser) extractPublicIPs(ips []net.IP) (publicIPs []net.IP) {
	for _, ip := range ips {
		if p.ipIsPrivate(ip) {
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

//nolint:gomnd
func privateIPNets() [8]net.IPNet {
	return [8]net.IPNet{
		{ // localhost
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPv4Mask(255, 0, 0, 0),
		},
		{ // 24-bit block
			IP:   net.IP{10, 0, 0, 0},
			Mask: net.IPv4Mask(255, 0, 0, 0),
		},
		{ // 20-bit block
			IP:   net.IP{172, 16, 0, 0},
			Mask: net.IPv4Mask(255, 240, 0, 0),
		},
		{ // 16-bit block
			IP:   net.IP{192, 168, 0, 0},
			Mask: net.IPv4Mask(255, 255, 0, 0),
		},
		{ // link local address
			IP:   net.IP{169, 254, 0, 0},
			Mask: net.IPv4Mask(255, 255, 0, 0),
		},
		{ // localhost IPv6
			IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{ // unique local address IPv6
			IP:   net.IP{0xfc, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Mask: net.IPMask{254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{ // link local address IPv6
			IP:   net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Mask: net.IPMask{255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
}
