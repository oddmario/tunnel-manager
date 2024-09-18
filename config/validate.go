package config

import (
	"github.com/oddmario/tunnel-manager/logger"
	"github.com/oddmario/tunnel-manager/utils"
)

func validateConfig() {
	if Config.Mode != "tunnel_host" && Config.Mode != "backend_server" {
		logger.Fatal("Invalid operating mode. `mode` has to be either `tunnel_host` or `backend_server`")
	}

	if !utils.DoesNetworkInterfaceExist(Config.MainNetworkInterface) {
		logger.Fatal("Invalid main network interface specified")
	}
}
