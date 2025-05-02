//go:build linux

package wireguard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/pedramktb/go-base-lib/lifecycle"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const StopTimeout = 5 * time.Second

func (m *Manager) start(ctx context.Context) error {
	_ = m.stop(ctx)

	go func() {
		done, err := lifecycle.RegisterCloser(ctx)
		if err == nil {
			defer done()
		}
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), StopTimeout)
		defer cancel()
		_ = m.stop(ctx)
	}()

	wgClient, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to create wireguard client: %w", err)
	}
	defer wgClient.Close()

	err = exec.CommandContext(ctx, "ip", "link", "add", m.interfaceName, "type", "wireguard").Run()
	if err != nil {
		return fmt.Errorf("failed to create wireguard interface: %w", err)
	}

	listenPort := int(m.listenPort)

	err = wgClient.ConfigureDevice(m.interfaceName, wgtypes.Config{
		PrivateKey:   &m.privateKey,
		ListenPort:   &listenPort,
		ReplacePeers: true,
	})
	if err != nil {
		return fmt.Errorf("failed to configure wireguard device: %w", err)
	}

	if m.netV4 != nil {
		err = exec.CommandContext(ctx, "ip", "-4", "address", "add", m.netV4.String(), "dev", m.interfaceName).Run()
		if err != nil {
			return fmt.Errorf("failed to add ipv4 address to wireguard interface: %w", err)
		}
	}

	if m.netV6 != nil {
		err = exec.CommandContext(ctx, "ip", "-6", "address", "add", netV6String(*m.netV6), "dev", m.interfaceName).Run()
		if err != nil {
			return fmt.Errorf("failed to add ipv6 address to wireguard interface: %w", err)
		}
	}

	err = exec.CommandContext(ctx, "ip", "link", "set", "mtu", fmt.Sprint(m.mtu), "up", "dev", m.interfaceName).Run()
	if err != nil {
		return fmt.Errorf("failed to set mtu and start for wireguard interface: %w", err)
	}

	err = m.addRouting()
	if err != nil {
		return fmt.Errorf("failed to add iptables rules: %w", err)
	}

	return nil
}

func (m *Manager) stop(ctx context.Context) error {
	var errs []error

	err := m.removeRouting()
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to remove iptables rules: %w", err))
	}

	err = exec.CommandContext(ctx, "ip", "link", "delete", m.interfaceName).Run()
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to delete wireguard interface: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (m *Manager) Restart(ctx context.Context) error {
	err := exec.CommandContext(ctx, "ip", "link", "set", "down", m.interfaceName).Run()
	if err != nil {
		return fmt.Errorf("failed to stop wireguard interface: %w", err)
	}

	err = exec.CommandContext(ctx, "ip", "link", "set", "up", m.interfaceName).Run()
	if err != nil {
		return fmt.Errorf("failed to start wireguard interface: %w", err)
	}

	return nil
}
