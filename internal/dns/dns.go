package dns

import (
	"context"
	"net"
	"time"
)

type DNS string

func (s DNS) SetDefault() {
	if s == "" {
		return
	}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := &net.Dialer{
				Timeout: 3 * time.Second,
			}
			return d.DialContext(ctx, network, net.JoinHostPort(string(s), "53"))
		},
	}
}
