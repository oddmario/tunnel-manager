package config

import "github.com/oddmario/tunnels-manager/tunnel"

type ConfigObject struct {
	Mode                             string
	MainNetworkInterface             string
	Tunnels                          []*tunnel.Tunnel
	DynamicIPUpdaterAPIIsEnabled     bool
	DynamicIPUpdaterAPIListenAddress string
	DynamicIPUpdaterAPIListenPort    int
	PingTimeout                      int
	PingInterval                     int
	DynamicIPUpdateTimeout           int
	DynamicIPUpdateInterval          int
}
