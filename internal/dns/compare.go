package dns

import (
	"fmt"
	"net"
)

func EqualsTo(name, want string) (bool, error) {
	ips, err := net.LookupIP(name)
	if err != nil {
		return false, err
	}
	if len(ips) > 1 {
		return false, fmt.Errorf("domain: %s has too many dns records", name)
	}
	return ips[0].String() == want, nil
}
