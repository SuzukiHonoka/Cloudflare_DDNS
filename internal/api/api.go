package api

import (
	"cf_ddns/utils"
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

type API struct {
	URL  string `json:"url"`
	IPv6 bool   `json:"ipv6"`
}

func (x *API) GetIP() (string, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var d net.Dialer
		if x.IPv6 {
			// tcp6 only but seems cloudflare api doesn't support it yet
			return d.DialContext(ctx, "tcp6", addr)
		}
		return d.DialContext(ctx, "tcp4", addr)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(x.URL)
	if err != nil {
		return "", err
	}
	defer utils.ForceClose(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
