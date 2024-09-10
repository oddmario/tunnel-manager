package tunnel

type Tunnel struct {
	IsInitialised                      bool
	TunnelDriver                       string
	TunHostMainPublicIP                string
	TunHostPublicIP                    string
	BackendServerPublicIP              string
	TunnelKey                          int
	TunnelInterfaceName                string
	TunnelRoutingTablesID              int
	TunnelRoutingTablesName            string
	TunnelGatewayIP                    string
	TunHostTunnelIP                    string
	BackendServerTunnelIP              string
	TunnelType                         string
	SplitTunnelPorts                   []map[string]interface{}
	ShouldRouteAllTrafficThroughTunnel bool
	DynamicIPUpdaterKey                string
	WGPrivateKeyFilePath               string
	WGServerTunnelHostListenPort       int
	WGServerBackendServerListenPort    int
	WGTunnelHostPubKey                 string
	WGBackendServerPubKey              string
}
