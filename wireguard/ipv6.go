package wireguard

import (
	"fmt"
	"net"
)

func ipV6String(ip net.IP) string {
	if ip.To4() != nil && len(ip) == net.IPv6len {
		return fmt.Sprintf("::ffff:%x:%x", uint16(ip[12])<<8|uint16(ip[13]), uint16(ip[14])<<8|uint16(ip[15]))
	}
	return ip.String()
}

func netV6String(ipnet net.IPNet) string {
	if ones, total := ipnet.Mask.Size(); total == 128 && ipnet.IP.To4() != nil && len(ipnet.IP) == net.IPv6len {
		return fmt.Sprintf("%s/%d", ipV6String(ipnet.IP), ones)
	}
	return ipnet.String()
}

func addrV6String(addr net.Addr) string {
	ip, port, _ := net.SplitHostPort(addr.String())
	return fmt.Sprintf("[%s]:%s", ipV6String(net.ParseIP(ip)), port)
}
