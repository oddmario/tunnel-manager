package tunnel

type Tunnel struct {
	IsInitialised                      bool
	GREHostMainPublicIP                string
	GREHostPublicIP                    string
	BackendServerPublicIP              string
	TunnelKey                          int
	TunnelInterfaceName                string
	TunnelRoutingTablesID              int
	TunnelRoutingTablesName            string
	TunnelGatewayIP                    string
	GREHostTunnelIP                    string
	BackendServerTunnelIP              string
	TunnelType                         string
	SplitTunnelPorts                   []map[string]interface{}
	ShouldRouteAllTrafficThroughTunnel bool
	DynamicIPUpdaterKey                string
}
