package config

import (
	"log"
	"os"

	"github.com/oddmario/tunnel-manager/constants"
	"github.com/oddmario/tunnel-manager/tunnel"
	"github.com/oddmario/tunnel-manager/utils"
	"github.com/tidwall/gjson"
)

var config gjson.Result
var Config *ConfigObject = &ConfigObject{
	Mode:                             "",
	ApplyKernelTuningTweaks:          false,
	MainNetworkInterface:             "",
	Tunnels:                          []*tunnel.Tunnel{},
	DynamicIPUpdaterAPIIsEnabled:     false,
	DynamicIPUpdaterAPIListenAddress: "",
	DynamicIPUpdaterAPIListenPort:    0,
	PingTimeout:                      0,
	PingInterval:                     0,
	DynamicIPUpdateTimeout:           0,
	DynamicIPUpdateInterval:          0,
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
	Config.ApplyKernelTuningTweaks = config.Get("apply_kernel_tuning_tweaks").Bool()
	Config.MainNetworkInterface = config.Get("main_network_interface").String()

	Config.Tunnels = tunnel.TunsFromJson(config.Get("tunnels"))

	Config.DynamicIPUpdaterAPIIsEnabled = config.Get("dynamic_ip_updater_api.is_enabled").Bool()
	Config.DynamicIPUpdaterAPIListenAddress = config.Get("dynamic_ip_updater_api.listen_address").String()
	Config.DynamicIPUpdaterAPIListenPort = int(config.Get("dynamic_ip_updater_api.listen_port").Int())

	Config.PingTimeout = int(config.Get("timeouts.ping_timeout").Int())
	Config.PingInterval = int(config.Get("timeouts.ping_interval").Int())
	Config.DynamicIPUpdateTimeout = int(config.Get("timeouts.dynamic_ip_update_timeout").Int())
	Config.DynamicIPUpdateInterval = int(config.Get("timeouts.dynamic_ip_update_attempt_interval").Int())
}
