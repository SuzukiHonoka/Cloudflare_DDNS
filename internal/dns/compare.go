package dns

import (
	"fmt"
	"net"
)

func EqualsTo(name string, want net.IP) (bool, error) {
	// Lookup the IP address of the domain
	ips, err := net.LookupIP(name)
	if err != nil {
		return false, err
	}

	// If there are more than one IP address, return an error
	if len(ips) > 1 {
		return false, fmt.Errorf("domain: %s has too many dns records", name)
	}

	// Compare the IP address
	ok := ips[0].Equal(want)
	return ok, nil
}
