package tunnel

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oddmario/gre-manager/utils"
	"github.com/oddmario/gre-manager/workers"
)

func (t *Tunnel) sendIPToGREHost(dynamic_ip_updater_api_listen_port int) {
	for {
		external_ip, err := utils.GetExternalIP()

		if err != nil {
			fmt.Println("[WARN] Unable to send the public IP address of the backend to the GRE host. Retrying in 3 seconds...")

			time.Sleep(3 * time.Second)

			continue
		}

		req, _ := resty.New().R().
			SetHeader("X-Key", t.DynamicIPUpdaterKey).
			SetBody(map[string]interface{}{"new_ip": external_ip}).
			Post("http://" + t.GREHostMainPublicIP + ":" + utils.IToStr(dynamic_ip_updater_api_listen_port) + "/update_ip")
		if req.StatusCode() != 200 {
			fmt.Println("[WARN] Unable to send the public IP address of the backend to the GRE host. Retrying in 3 seconds...")

			time.Sleep(3 * time.Second)

			continue
		}

		t.BackendServerPublicIP = external_ip

		break
	}
}

func (t *Tunnel) Init(mode, main_network_interface string, dynamic_ip_updater_api_listen_port int, startNetworkStateMonitorThread bool) bool {
	if t.IsInitialised {
		fmt.Println("[WARN] Failed to initialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": The tunnel has already been initialised. Ignoring tunnel initialisation.")

		return false
	}

	if utils.DoesNetworkInterfaceExist(t.TunnelInterfaceName) {
		fmt.Println("[WARN] Failed to initialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": GRE tunnel interface already exists. Ignoring tunnel initialisation.")

		return false
	}

	if mode == "gre_host" {
		if t.BackendServerPublicIP == "DYNAMIC" {
			fmt.Println("[WARN] Failed to initialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": The initial backend IP is `DYNAMIC` but no IP has been received from the backend yet. Ignoring tunnel initialisation.")

			return false
		}

		utils.Cmd("ip tunnel add "+t.TunnelInterfaceName+" mode gre local "+t.GREHostMainPublicIP+" remote "+t.BackendServerPublicIP+" ttl 255 key "+utils.IToStr(t.TunnelKey), true)
		utils.Cmd("ip addr add "+t.GREHostTunnelIP+"/30 dev "+t.TunnelInterfaceName, true)
		utils.Cmd("ip link set "+t.TunnelInterfaceName+" up", true)

		utils.Cmd("iptables -A FORWARD -i gre+ -j ACCEPT", true)
		utils.Cmd("iptables -A FORWARD -d "+t.BackendServerTunnelIP+" -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT", true)
		utils.Cmd("iptables -A FORWARD -s "+t.BackendServerTunnelIP+" -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT", true)
		utils.Cmd("iptables -t nat -A POSTROUTING -s "+t.TunnelGatewayIP+"/30 ! -o gre+ -j SNAT --to-source "+t.GREHostPublicIP, true)

		if t.TunnelType == "full" {
			utils.Cmd("iptables -t nat -A PREROUTING -d "+t.GREHostPublicIP+" -j DNAT --to-destination "+t.BackendServerTunnelIP, true)
		} else {
			for _, port := range t.SplitTunnelPorts {
				p := utils.IToStr(port["port"].(int))
				proto := port["proto"].(string)

				utils.Cmd("iptables -t nat -A PREROUTING -d "+t.GREHostPublicIP+" -p "+proto+" -m "+proto+" --dport "+p+" -j DNAT --to-destination "+t.BackendServerTunnelIP, true)
			}
		}
	}

	if mode == "backend_server" {
		var backendIsDynamicIP bool = false

		routingTableExists, err := rttablesCheck(t.TunnelRoutingTablesID, t.TunnelRoutingTablesName)
		if err != nil {
			fmt.Println("[WARN] Failed to initialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": Routing table check failed -> " + err.Error() + ". Ignoring tunnel initialisation.")

			return false
		}

		if !routingTableExists {
			err := rttablesWrite(t.TunnelRoutingTablesID, t.TunnelRoutingTablesName)

			if err != nil {
				fmt.Println("[WARN] Failed to initialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": Routing table write failed -> " + err.Error() + ". Ignoring tunnel initialisation.")

				return false
			}
		}

		if t.BackendServerPublicIP == "DYNAMIC" {
			backendIsDynamicIP = true

			t.sendIPToGREHost(dynamic_ip_updater_api_listen_port)
		}

		utils.Cmd("ip tunnel add "+t.TunnelInterfaceName+" mode gre local "+t.BackendServerPublicIP+" remote "+t.GREHostMainPublicIP+" ttl 255 key "+utils.IToStr(t.TunnelKey), true)
		utils.Cmd("ip addr add "+t.BackendServerTunnelIP+"/30 dev "+t.TunnelInterfaceName, true)
		utils.Cmd("ip link set "+t.TunnelInterfaceName+" up", true)

		utils.Cmd("ip rule add from "+t.TunnelGatewayIP+"/30 table "+t.TunnelRoutingTablesName, true)
		utils.Cmd("ip route add default via "+t.GREHostTunnelIP+" table "+t.TunnelRoutingTablesName, true)

		if t.ShouldRouteAllTrafficThroughTunnel {
			gatewayIP, err := utils.Cmd("ip route show 0.0.0.0/0 dev "+main_network_interface+" | cut -d\\  -f3", true)
			if err == nil {
				gatewayIPString := strings.TrimSpace(utils.BytesToString(gatewayIP))

				utils.Cmd("echo 'nameserver 1.1.1.1' > /etc/resolv.conf", true)
				utils.Cmd("echo 'nameserver 1.0.0.1' >> /etc/resolv.conf", true)

				utils.Cmd("ip route add "+t.GREHostMainPublicIP+" via "+gatewayIPString+" dev "+main_network_interface+" onlink", true)
				utils.Cmd("ip route replace default via "+t.GREHostTunnelIP, true)
			}
		}

		if backendIsDynamicIP && startNetworkStateMonitorThread {
			go workers.MonitorInterfaceStateChange(main_network_interface, func(state string) {
				if state == "UP" {
					fmt.Println("[INFO] Tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": The main network interface has come up. Reinitialising the tunnel...")

					t.BackendServerPublicIP = "DYNAMIC"

					t.Init(mode, main_network_interface, dynamic_ip_updater_api_listen_port, false)
				} else {
					fmt.Println("[INFO] Tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": The main network interface has gone down. Deinitialising the tunnel till it comes up again...")

					t.Deinit(mode, main_network_interface, false)
				}
			})
		}
	}

	utils.Cmd("tc qdisc replace dev "+t.TunnelInterfaceName+" root fq_codel", true)
	utils.Cmd("ip link set "+t.TunnelInterfaceName+" txqueuelen 15000", true)
	utils.Cmd("ethtool -K "+t.TunnelInterfaceName+" gro off gso off tso off", true)

	t.IsInitialised = true

	fmt.Println("[DEBUG] The GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + " has been setup successfully.")

	return true
}

func (t *Tunnel) Deinit(mode, main_network_interface string, ignoreInitialisationStatus bool) bool {
	if !t.IsInitialised && !ignoreInitialisationStatus {
		fmt.Println("[WARN] Failed to deinitialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": The tunnel was not initialised. Ignoring tunnel deinitialisation.")

		return false
	}

	if !utils.DoesNetworkInterfaceExist(t.TunnelInterfaceName) {
		fmt.Println("[WARN] Failed to deinitialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": GRE tunnel interface does not exist. Ignoring tunnel deinitialisation.")

		return false
	}

	if mode == "gre_host" {
		utils.Cmd("iptables -D FORWARD -d "+t.BackendServerTunnelIP+" -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT", true)
		utils.Cmd("iptables -D FORWARD -s "+t.BackendServerTunnelIP+" -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT", true)
		utils.Cmd("iptables -t nat -D POSTROUTING -s "+t.TunnelGatewayIP+"/30 ! -o gre+ -j SNAT --to-source "+t.GREHostPublicIP, true)

		if t.TunnelType == "full" {
			utils.Cmd("iptables -t nat -D PREROUTING -d "+t.GREHostPublicIP+" -j DNAT --to-destination "+t.BackendServerTunnelIP, true)
		} else {
			for _, port := range t.SplitTunnelPorts {
				p := utils.IToStr(port["port"].(int))
				proto := port["proto"].(string)

				utils.Cmd("iptables -t nat -D PREROUTING -d "+t.GREHostPublicIP+" -p "+proto+" -m "+proto+" --dport "+p+" -j DNAT --to-destination "+t.BackendServerTunnelIP, true)
			}
		}

		utils.Cmd("ip addr del "+t.GREHostTunnelIP+"/30 dev "+t.TunnelInterfaceName, true)
		utils.Cmd("ip link set "+t.TunnelInterfaceName+" down", true)
		utils.Cmd("ip tunnel del "+t.TunnelInterfaceName, true)
	}

	if mode == "backend_server" {
		routingTableExists, _ := rttablesCheck(t.TunnelRoutingTablesID, t.TunnelRoutingTablesName)

		if routingTableExists {
			err := rttablesDel(t.TunnelRoutingTablesID, t.TunnelRoutingTablesName)

			if err != nil {
				fmt.Println("[WARN] Failed to deinitialise the GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + ": Routing table delete failed -> " + err.Error() + ". Continuing tunnel deinitialisation anyway...")
			}
		}

		if t.ShouldRouteAllTrafficThroughTunnel {
			gatewayIP, err := utils.Cmd("ip route show 0.0.0.0/0 dev "+main_network_interface+" | cut -d\\  -f3", true)
			if err == nil {
				gatewayIPString := strings.TrimSpace(utils.BytesToString(gatewayIP))

				utils.Cmd("ip route del default via "+t.GREHostTunnelIP, true)
				utils.Cmd("ip route del "+t.GREHostMainPublicIP+" via "+gatewayIPString+" dev "+main_network_interface+" onlink", true)
			}
		}

		utils.Cmd("ip route del default via "+t.GREHostTunnelIP+" table "+t.TunnelRoutingTablesName, true)
		utils.Cmd("ip rule del from "+t.TunnelGatewayIP+"/30 table "+t.TunnelRoutingTablesName, true)
		utils.Cmd("ip addr del "+t.BackendServerTunnelIP+"/30 dev "+t.TunnelInterfaceName, true)
		utils.Cmd("ip link set "+t.TunnelInterfaceName+" down", true)
		utils.Cmd("ip tunnel del "+t.TunnelInterfaceName, true)
	}

	t.IsInitialised = false

	fmt.Println("[DEBUG] The GRE tunnel " + t.GREHostMainPublicIP + " <-> " + t.BackendServerPublicIP + " has been removed successfully.")

	return true
}
