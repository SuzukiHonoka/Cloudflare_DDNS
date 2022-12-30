package config

import (
	"cf_ddns/internal/api"
	"cf_ddns/internal/dns"
	"cf_ddns/internal/proxy"
	"encoding/json"
	"os"
)

type Config struct {
	API    *api.API           `json:"api"`
	DNS    dns.DNS            `json:"dns,omitempty"`
	Target *proxy.Target      `json:"target"`
	Creds  *proxy.Credentials `json:"credentials"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c *Config
	err = json.Unmarshal(b, &c)
	return c, err
}
