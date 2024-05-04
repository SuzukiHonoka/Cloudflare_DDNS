package handler

import (
	"cf_ddns/internal/dns"
	"cf_ddns/internal/wrapper"
	"cf_ddns/pkg/config"
	"log"
)

func HandleConfig(conf *config.Config) error {
	// Set default dns
	conf.DNS.SetDefault()

	// Get IP from API
	ip, err := conf.API.GetIP()
	if err != nil {
		return err
	}

	// Print the result
	log.Printf("Public IP: %s", ip)

	// Check if record match
	ok, err := dns.EqualsTo(conf.Target.String(), ip)
	if err != nil {
		return err
	}
	if ok {
		log.Println("Skip sync operation, record matches")
		return nil
	}

	// Create API wrapper instance
	api, err := wrapper.New(conf.Credential)
	if err != nil {
		return err
	}

	// Initialize API wrapper
	if err = api.Init(conf.Target); err != nil {
		// Second shoot
		// maybe the specified record does not exist or else, try to create one anyway
		log.Println("creating record")
		_, err = api.CreateRecord(ip)
		if err != nil {
			return err
		}

		// Issue seems to be resolved, skip update
		return nil
	}

	// Update record
	if err = api.UpdateRecord(ip); err != nil {
		return err
	}

	log.Println("Record updated")
	return nil
}
