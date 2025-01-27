package wireguard

import (
	"context"
	"fmt"
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Manager struct {
	interfaceName  string
	privateKey     wgtypes.Key
	listenPort     uint16
	mtu            uint16
	netV4, netV6   net.IPNet
	addrV4, addrV6 net.IP
	dnsV4, dnsV6   net.UDPAddr
}

func NewManager(
	ctx context.Context,
	interfaceName,
	privateKey string,
	listenPort uint16,
	mtu uint16,
	netV4, netV6 net.IPNet,
	addrV4, addrV6 net.IP,
	dnsV4, dnsV6 net.UDPAddr,
) (*Manager, error) {
	key, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	wgManager := &Manager{
		interfaceName: interfaceName,
		privateKey:    key,
		listenPort:    listenPort,
		mtu:           mtu,
		netV4:         netV4,
		netV6:         netV6,
		addrV4:        addrV4,
		addrV6:        addrV6,
		dnsV4:         dnsV4,
		dnsV6:         dnsV6,
	}

	err = wgManager.start(ctx)
	if err != nil {
		return nil, err
	}

	return wgManager, nil
}

func (m *Manager) PublicKey() string {
	return m.privateKey.PublicKey().String()
}

func (m *Manager) Add(publicKey, address string) error {
	wgClient, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to create wireguard client: %w", err)
	}
	defer wgClient.Close()

	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	err = wgClient.ConfigureDevice(m.interfaceName, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{{
			PublicKey: key,
			AllowedIPs: []net.IPNet{{
				IP:   net.ParseIP(address),
				Mask: net.CIDRMask(32, 32),
			}},
		}},
		ReplacePeers: false,
	})
	if err != nil {
		return fmt.Errorf("failed to add peer to wireguard device: %w", err)
	}

	return nil
}

func (m *Manager) Remove(publicKey string) error {
	wgClient, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to create wireguard client: %w", err)
	}
	defer wgClient.Close()

	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	err = wgClient.ConfigureDevice(m.interfaceName, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{{
			PublicKey: key,
			Remove:    true,
		}},
		ReplacePeers: false,
	})
	if err != nil {
		return fmt.Errorf("failed to remove peer from wireguard device: %w", err)
	}

	return nil
}
