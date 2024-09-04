package config

import (
	"log"

	"github.com/oddmario/gre-manager/utils"
)

func validateConfig() {
	if Config.Mode != "gre_host" && Config.Mode != "backend_server" {
		log.Fatal("[ERROR] Invalid operating mode. `mode` has to be either `gre_host` or `backend_server`")
	}

	if !utils.DoesNetworkInterfaceExist(Config.MainNetworkInterface) {
		log.Fatal("[ERROR] Invalid main network interface specified")
	}
}
