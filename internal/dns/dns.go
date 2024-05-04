package dns

import (
	"context"
	"net"
	"time"
)

var DefaultTimeout = 3 * time.Second

type DNS string

func (s DNS) SetDefault() {
	if s == "" {
		return
	}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := &net.Dialer{
				Timeout: DefaultTimeout,
			}
			return d.DialContext(ctx, network, net.JoinHostPort(string(s), "53"))
		},
	}
}
