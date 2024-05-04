package api

import (
	"cf_ddns/utils"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var DefaultTimeout = 5 * time.Second

type API struct {
	URL  string `json:"url"`
	IPv6 bool   `json:"ipv6"`
	UA   string `json:"ua"`
}

func (x *API) GetIP() (net.IP, error) {
	// Clone default transport
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var d net.Dialer
		if x.IPv6 {
			// tcp6 only but seems cloudflare api doesn't support it yet
			return d.DialContext(ctx, "tcp6", addr)
		}
		return d.DialContext(ctx, "tcp4", addr)
	}

	// Custom http client with transport and timeout
	client := &http.Client{
		Timeout:   DefaultTimeout,
		Transport: transport,
	}

	// Building request
	req, err := http.NewRequest(http.MethodGet, x.URL, nil)
	if err != nil {
		return nil, err
	}

	// Custom User-Agent
	if x.UA != "" {
		req.Header.Set("User-Agent", x.UA)
	}

	// Actually do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Always cloe the body
	defer utils.ForceClose(resp.Body)

	// Read response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert body to str
	ipStr := string(b)

	// Parse IP since API may return IPv4-mapped IPv6 ("::ffff:192.0.2.1")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		err = fmt.Errorf("IP string: %s is not valid", ipStr)
	}
	return ip, err
}
