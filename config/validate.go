package config

import (
	"log"

	"github.com/oddmario/tunnel-manager/utils"
)

func validateConfig() {
	if Config.Mode != "tunnel_host" && Config.Mode != "backend_server" {
		log.Fatal("[ERROR] Invalid operating mode. `mode` has to be either `tunnel_host` or `backend_server`")
	}

	if !utils.DoesNetworkInterfaceExist(Config.MainNetworkInterface) {
		log.Fatal("[ERROR] Invalid main network interface specified")
	}
}
