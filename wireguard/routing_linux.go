//go:build linux

package wireguard

import (
	"errors"
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

func (m *Manager) addRouting() error {
	if m.netV4 != nil {
		ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
		if err != nil {
			return fmt.Errorf("failed to create iptables client: %w", err)
		}
		err = ipt.Insert("nat", "POSTROUTING", 1, "-s", m.netV4.String(), "-j", "MASQUERADE")
		if err != nil {
			return fmt.Errorf("failed to add masquerade rule for ipv4: %w", err)
		}

		err = ipt.Insert("nat", "PREROUTING", 1, "-d", m.addrV4.String(), "-p", "udp", "--dport", "53", "-j", "DNAT", "--to-destination", m.dnsV4.String())
		if err != nil {
			return fmt.Errorf("failed to add dns rule for ipv4: %w", err)
		}
	}

	if m.netV6 != nil {
		ip6t, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
		if err != nil {
			return fmt.Errorf("failed to create iptables client: %w", err)
		}

		err = ip6t.Insert("nat", "POSTROUTING", 1, "-s", netV6String(*m.netV6), "-j", "MASQUERADE")
		if err != nil {
			return fmt.Errorf("failed to add masquerade rule for ipv6: %w", err)
		}

		err = ip6t.Insert("nat", "PREROUTING", 1, "-d", ipV6String(m.addrV6), "-p", "udp", "--dport", "53", "-j", "DNAT", "--to-destination", addrV6String(m.dnsV6))
		if err != nil {
			return fmt.Errorf("failed to add dns rule for ipv6: %w", err)
		}
	}

	return nil
}

func (m *Manager) removeRouting() error {
	var errs []error

	if m.netV4 != nil {
		ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create iptables client: %w", err))
		}

		if ipt != nil {
			err = ipt.DeleteIfExists("nat", "POSTROUTING", "-s", m.netV4.String(), "-j", "MASQUERADE")
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to delete masquerade rule for ipv4: %w", err))
			}

			err = ipt.DeleteIfExists("nat", "PREROUTING", "-d", m.addrV4.String(), "-p", "udp", "--dport", "53", "-j", "DNAT", "--to-destination", m.dnsV4.String())
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to delete dns rule for ipv4: %w", err))
			}
		}
	}

	if m.netV6 != nil {
		ip6t, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create iptables client: %w", err))
		}

		if ip6t != nil {
			err = ip6t.DeleteIfExists("nat", "POSTROUTING", "-s", m.netV6.String(), "-j", "MASQUERADE")
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to delete masquerade rule for ipv6: %w", err))
			}

			err = ip6t.DeleteIfExists("nat", "PREROUTING", "-d", m.addrV4.String(), "-p", "udp", "--dport", "53", "-j", "DNAT", "--to-destination", addrV6String(m.dnsV6))
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to delete dns rule for ipv6: %w", err))
			}
		}
	}

	return errors.Join(errs...)
}
