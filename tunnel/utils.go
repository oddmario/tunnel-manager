package tunnel

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/oddmario/gre-manager/utils"
	"github.com/tidwall/gjson"
)

func TunsFromJson(j gjson.Result) []*Tunnel {
	tuns := []*Tunnel{}

	j.ForEach(func(key, value gjson.Result) bool {
		if value.Get("tunnel_type").String() != "split" && value.Get("tunnel_type").String() != "full" {
			fmt.Println("[WARN] Failed to initialise the GRE tunnel " + value.Get("gre_host_main_public_ip").String() + " <-> " + value.Get("backend_server_public_ip").String() + ": The tunnel type has to be either `split` or `full`. Ignoring tunnel initialisation.")

			return true
		}

		var splitTunnelPorts []map[string]interface{} = []map[string]interface{}{}

		if value.Get("tunnel_type").String() == "split" {
			value.Get("split_tunnel_ports").ForEach(func(key_port, value_port gjson.Result) bool {
				proto := strings.ToLower(value_port.Get("proto").String())

				if proto != "tcp" && proto != "udp" {
					fmt.Println("[WARN] Failed to configure split tunnel port for " + value.Get("gre_host_main_public_ip").String() + " <-> " + value.Get("backend_server_public_ip").String() + ": Invalid split tunnel port protocol specified. Only TCP & UDP are allowed. Ignoring port rule.")

					return true
				}

				splitTunnelPorts = append(splitTunnelPorts, map[string]interface{}{
					"proto": proto,
					"port":  int(value_port.Get("port").Int()),
				})

				return true
			})
		}

		tuns = append(tuns, &Tunnel{
			IsInitialised:                      false,
			GREHostMainPublicIP:                value.Get("gre_host_main_public_ip").String(),
			GREHostPublicIP:                    value.Get("gre_host_public_ip").String(),
			BackendServerPublicIP:              value.Get("backend_server_public_ip").String(),
			TunnelKey:                          int(value.Get("tunnel_key").Int()),
			TunnelInterfaceName:                value.Get("tunnel_interface_name").String(),
			TunnelRoutingTablesID:              int(value.Get("tunnel_rttables_id").Int()),
			TunnelRoutingTablesName:            value.Get("tunnel_rttables_name").String(),
			TunnelGatewayIP:                    value.Get("tunnel_gateway_ip").String(),
			GREHostTunnelIP:                    value.Get("gre_host_tunnel_ip").String(),
			BackendServerTunnelIP:              value.Get("backend_server_tunnel_ip").String(),
			TunnelType:                         value.Get("tunnel_type").String(),
			SplitTunnelPorts:                   splitTunnelPorts,
			ShouldRouteAllTrafficThroughTunnel: value.Get("route_all_traffic_through_tunnel").Bool(),
			DynamicIPUpdaterKey:                value.Get("dynamic_ip_updater_key").String(),
		})

		return true
	})

	return tuns
}

func rttablesCheck(table_id int, table_name string) (bool, error) {
	var routingTableFound bool = false

	os.MkdirAll("/etc/iproute2/", os.ModePerm)

	file_create, err := os.OpenFile("/etc/iproute2/rt_tables", os.O_CREATE|os.O_EXCL, 0644)
	if err == nil {
		file_create.Close()
	} else {
		if !os.IsExist(err) {
			return routingTableFound, errors.New("failed to create /etc/iproute2/rt_tables")
		}
	}

	file, err := os.Open("/etc/iproute2/rt_tables")
	if err != nil {
		return routingTableFound, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) <= 0 || strings.HasPrefix(line, "#") {
			continue
		}

		split := strings.Split(line, " ")

		if len(split) != 2 {
			continue
		}

		id := split[0]
		name := split[1]

		if utils.StrToI(id) == table_id && name == table_name {
			routingTableFound = true
		}

		if utils.StrToI(id) == table_id && name != table_name {
			return routingTableFound, errors.New("another routing table with that id already exists. please try using another id in your config file")
		}
	}

	if err := scanner.Err(); err != nil {
		return routingTableFound, err
	}

	return routingTableFound, nil
}

func rttablesWrite(table_id int, table_name string) error {
	file_w, err := os.OpenFile("/etc/iproute2/rt_tables", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file_w.Close()

	writer := bufio.NewWriter(file_w)
	writer.WriteString(utils.IToStr(table_id) + " " + table_name + "\n")
	writer.Flush()

	return nil
}

func rttablesDel(table_id int, table_name string) error {
	file, err := os.Open("/etc/iproute2/rt_tables")
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, utils.IToStr(table_id)+" "+table_name) {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file_w, err := os.OpenFile("/etc/iproute2/rt_tables", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file_w.Close()

	writer := bufio.NewWriter(file_w)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return nil
}