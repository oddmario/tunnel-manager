package dynamicipupdater

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oddmario/tunnel-manager/config"
	"github.com/oddmario/tunnel-manager/tunnel"
	"github.com/oddmario/tunnel-manager/utils"
	"github.com/tidwall/gjson"
)

func updateIPHandler(c *gin.Context) {
	key := c.GetHeader("X-Key")

	if len(key) <= 0 {
		c.Status(http.StatusUnauthorized)

		return
	}

	var isKeyFound bool = false
	var tunnel *tunnel.Tunnel = nil

	for _, tun := range config.Config.Tunnels {
		if tun.DynamicIPUpdaterKey == key {
			isKeyFound = true
			tunnel = tun

			break
		}
	}

	if !isKeyFound {
		c.Status(http.StatusUnauthorized)

		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)

		return
	}

	jsonDataStr := utils.BytesToString(jsonData)

	if !gjson.Valid(jsonDataStr) {
		c.Status(http.StatusBadRequest)

		return
	}

	jsonDataParser := gjson.Parse(jsonDataStr)

	newIP := jsonDataParser.Get("new_ip").String()

	tunnel.BackendServerPublicIP = newIP

	if tunnel.IsInitialised {
		utils.Cmd("ip tunnel change "+tunnel.TunnelInterfaceName+" mode gre local "+tunnel.TunHostMainPublicIP+" remote "+tunnel.BackendServerPublicIP+" ttl 255 key "+utils.IToStr(tunnel.TunnelKey), true)
	} else {
		tunnel.Init(config.Config.Mode, config.Config.MainNetworkInterface, config.Config.DynamicIPUpdaterAPIListenPort, config.Config.DynamicIPUpdateInterval, config.Config.DynamicIPUpdateTimeout, config.Config.PingInterval, config.Config.PingTimeout)
	}

	c.Status(http.StatusOK)
}
