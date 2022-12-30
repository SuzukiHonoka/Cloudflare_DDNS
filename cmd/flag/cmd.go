package flag

import "flag"

var ConfigPath string

func init() {
	ConfigPath = *flag.String("c", "./config.json", "configuration filepath")
	flag.Parse()
}
