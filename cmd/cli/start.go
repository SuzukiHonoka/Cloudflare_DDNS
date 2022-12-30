package cli

import (
	"cf_ddns/cmd/entry"
	"cf_ddns/cmd/flag"
	"cf_ddns/internal/config"
	"cf_ddns/utils"
)

func Main() {
	conf, err := config.LoadConfig(flag.ConfigPath)
	utils.CheckErr(err)
	entry.Handle(conf)
}
