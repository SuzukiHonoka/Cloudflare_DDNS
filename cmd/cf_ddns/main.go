package main

import (
	"cf_ddns/internal/handler"
	"cf_ddns/pkg/config"
	"cf_ddns/utils"
	"flag"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("c", "./config.json", "configuration filepath")
	flag.Parse()

	// Load configuration
	conf, err := config.LoadConfig(*configPath)
	utils.PanicIfError(err)

	// Handle configuration
	err = handler.HandleConfig(conf)
	utils.PanicIfError(err)
}
