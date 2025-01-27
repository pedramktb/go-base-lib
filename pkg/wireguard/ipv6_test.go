package wireguard

import (
	"net"
	"testing"
)

func Test_ipV6String(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want string
	}{
		{
			name: "ipv4",
			ip:   net.ParseIP("10.0.0.0"),
			want: "::ffff:a00:0",
		},
		{
			name: "ipv4-compatibliity",
			ip:   net.ParseIP("::ffff:a00:0"),
			want: "::ffff:a00:0",
		},
		{
			name: "ipv4#2",
			ip:   net.ParseIP("255.127.63.31").To16(),
			want: "::ffff:ff7f:3f1f",
		},
		{
			name: "ipv4-compatibliity#2",
			ip:   net.ParseIP("::ffff:ff7f:3f1f").To4(),
			want: "255.127.63.31",
		},
		{
			name: "ipv4#3",
			ip:   net.ParseIP("255.255.255.255").To16(),
			want: "::ffff:ffff:ffff",
		},
		{
			name: "ipv6",
			ip:   net.ParseIP("2001:db8::68"),
			want: "2001:db8::68",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipV6String(tt.ip); got != tt.want {
				t.Errorf("forceIPv6String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_netV6String(t *testing.T) {
	tests := []struct {
		name  string
		ipnet net.IPNet
		want  string
	}{
		{
			name: "ipv4",
			ipnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.0"),
				Mask: net.CIDRMask(104, 128),
			},
			want: "::ffff:a00:0/104",
		},
		{
			name: "ipv4-compatibility",
			ipnet: net.IPNet{
				IP:   net.ParseIP("::ffff:a00:0"),
				Mask: net.CIDRMask(24, 32),
			},
			want: "10.0.0.0/24",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := netV6String(tt.ipnet); got != tt.want {
				t.Errorf("netV6String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addrV6String(t *testing.T) {
	tests := []struct {
		name string
		addr net.Addr
		want string
	}{
		{
			name: "ipv4",
			addr: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.0"),
				Port: 53,
			},
			want: "[::ffff:a00:0]:53",
		},
		{
			name: "ipv4-compatibility",
			addr: &net.TCPAddr{
				IP:   net.ParseIP("::ffff:a00:0").To4(),
				Port: 53,
			},
			want: "[::ffff:a00:0]:53",
		},
		{
			name: "ipv6",
			addr: &net.TCPAddr{
				IP:   net.ParseIP("2001:db8::68"),
				Port: 53,
			},
			want: "[2001:db8::68]:53",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addrV6String(tt.addr); got != tt.want {
				t.Errorf("addrV6String() = %v, want %v", got, tt.want)
			}
		})
	}
}
