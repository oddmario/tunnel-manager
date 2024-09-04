package config

import (
	"log"
	"os"

	"github.com/oddmario/gre-manager/constants"
	"github.com/oddmario/gre-manager/tunnel"
	"github.com/oddmario/gre-manager/utils"
	"github.com/tidwall/gjson"
)

var config gjson.Result
var Config *ConfigObject = &ConfigObject{
	Mode:                             "",
	MainNetworkInterface:             "",
	Tunnels:                          []*tunnel.Tunnel{},
	DynamicIPUpdaterAPIIsEnabled:     false,
	DynamicIPUpdaterAPIListenAddress: "",
	DynamicIPUpdaterAPIListenPort:    0,
}

func LoadConfig() {
	cfg_content, _ := os.ReadFile(constants.ConfigFilePath)
	cfgContentString := utils.BytesToString(cfg_content)

	if !gjson.Valid(cfgContentString) {
		log.Fatal("[ERROR] Malformed configuration file")
	} else {
		config = gjson.Parse(cfgContentString)
		storeEssentialConfigValues()
	}

	validateConfig()
}

func storeEssentialConfigValues() {
	Config.Mode = config.Get("mode").String()
	Config.MainNetworkInterface = config.Get("main_network_interface").String()

	Config.Tunnels = tunnel.TunsFromJson(config.Get("tunnels"))

	Config.DynamicIPUpdaterAPIIsEnabled = config.Get("dynamic_ip_updater_api.is_enabled").Bool()
	Config.DynamicIPUpdaterAPIListenAddress = config.Get("dynamic_ip_updater_api.listen_address").String()
	Config.DynamicIPUpdaterAPIListenPort = int(config.Get("dynamic_ip_updater_api.listen_port").Int())
}
