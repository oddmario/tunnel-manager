# 🚇 Tunnel Manager

An easy to use GRE, IPIP and WireGuard tunnel manager written in Go (Golang).

## ✨ Features
- Dynamic IP addresses support for the backend server(s)
- Full & split tunnels support
- Port forwarding configuration for split tunnels
- Easy multi-tunnel management, all in a single configuration file
- Support for GRE, IPIP and WireGuard

## 📙 Glossary
- **Backend server:** It is the server that you are trying to hide/protect the IP address of.
- **Tunnel host:** It is the VPS (or server in general) that has the public IP address that you want to expose publicly instead of the IP of the destination server. (e.g. a BuyVM VPS)

## 📝 Notes
For the optimal experience, kindly have a look at the notes listed at https://github.com/oddmario/GRE-setup-guide/blob/f0681a21edbc7a99f0d2a798529529a807357b5d/README.md#notes (you mainly need to pay attention to notes 1 to 5. anything that follows the 5th note can be ignored)

## 🧐 Configuration documentation
- `mode`
  * Can be either `backend_server` or `tunnel_host`.

- `apply_kernel_tuning_tweaks`
  * Should Tunnel Manager apply the recommended kernel tweaks for the purpose of tuning the network performance? Set this to "true" if yes.
  * The tweaks are applied by Tunnel Manager using the `sysctl` command.
  * The tweaks are mainly ones that increase the default networking & I/O limits of the Linux kernel.
  * The tweaks apply only on the tunnel host instance of Tunnel Manager.

- `main_network_interface`
  * The name of the main network interface on the system (e.g. eth0)

- `dynamic_ip_updater_api`
  * `is_enabled`: Whether to enable the dynamic IP updater API or no. **This can be enabled only on the tunnel host mode.** Note that you have to enable this if you have any tunnels with a **DYNAMIC** `backend_server_public_ip`. (this config parameter is ignored on the backend server mode)

  * `listen_address`: The IP address that will be used for the dynamic IP updater HTTP server. Make sure that it's binding to an IP address that the backend server(s) can access. (this config parameter is ignored on the backend server mode)

  * `listen_port`: The port that will be used for the dynamic IP updater HTTP server. (**This is not ignored on the backend server mode!** Note that if you are configuring the Tunnel-Manager copy of a backend server, you need to specify this to be the same `listen_port` configured on the tunnel host configuration file)

- `timeouts`: All the timeouts are in seconds.
  * `ping_timeout`: The maximum time allowed for a ping/ICMP request to finish.

  * `ping_interval`: The pause/sleep between each ping/ICMP request. (you can consider it a keepalive interval for dynamic IP changes detection)

  * `dynamic_ip_update_timeout`: The maximum time allowed for a "update_ip" HTTP request to finish.

  * `dynamic_ip_update_attempt_interval`: When a "update_ip" HTTP request fails, this is the pause/sleep before attempting to initiate another one.

- `tunnels`: An array of the tunnel(s) that you would like to have.
  * `driver`: The driver to use for the tunnel. Possible options are: gre, wireguard, ipip

  * `tunnel_host_main_public_ip`: The main/primary public IP address of the tunnel host.

  * `tunnel_host_public_ip`: The public IP address of the tunnel host that you would like to use instead of the backend server IP address. If your tunnel host has only one public IP address, make **tunnel_host_main_public_ip** and **tunnel_host_public_ip** equal. When they are equal, you will use your single public IP address of the tunnel host for the tunneling.

  * `backend_server_public_ip`: The public IP address of the backend server. Set this to `DYNAMIC` if the backend server has a dynamic IP address.

  * `tunnel_key`: The index/key of the tunnel. This has to be unique for each configured tunnel. (e.g. 1, 2, 3, etc). **It also has to be the same configured value for the tunnel on both the tunnel host configuration file & the backend one.**

  * `tunnel_interface_name`: The name of the tunnel interface. This has to be unique for each configured tunnel. (e.g. tun1, tun2, tun3, etc)

  * `tunnel_rttables_id`: The ID of the routing table used by the tunnel. This has to be unique for each configured tunnel. (e.g. 100, 200, 300, etc). **[This is ignored on the tunnel host as it is used only by the backend server instance. So it doesn't matter what value you set for `tunnel_rttables_id` on the tunnel host instance of Tunnel Manager]**

  * `tunnel_rttables_name`: The name of the routing table used by the tunnel. This has to be unique for each configured tunnel. (e.g. TUN1, TUN2, TUN3, etc). **[This is ignored on the tunnel host as it is used only by the backend server instance. So it doesn't matter what value you set for `tunnel_rttables_name` on the tunnel host instance of Tunnel Manager]**

  * `tunnel_gateway_ip`: The gateway that will be used by Tunnel Manager to setup the tunnel.

  * `tunnel_host_tunnel_ip`: The IP address of the tunnel host inside the tunnel.

  * `backend_server_tunnel_ip`: The IP address of the backend server inside the tunnel.

  * `tunnel_type`: Can be either **split** for a split tunnel, or **full** for a full tunnel. A full tunnel forwards all the ports, meanwhile a split tunnel forwards certain ports that you can configure in `split_tunnel_ports`. **[This is ignored on the backend server as it is used only by the tunnel host instance. So it doesn't matter what value you set for `tunnel_type` on a backend server instance of Tunnel Manager]**

  * `split_tunnel_ports`: An array containing the ports to forward for the purpose of split tunneling. This is ignored if `tunnel_type` is set to "full"
    * `proto`: Can be either TCP or UDP

    * `src_port`: The port(s) that the clients will connect to using the tunnel host public IP. To use a port range, you can use the `start_port:end_port` format (e.g. `8000:8050`).

    * `dest_port`: The port on the backend server. Keep this an empty string if you want it to be the same as `src_port`

    * More explanation on src & dest ports: Assume that 1.0.0.0 is the public IP address of the tunnel host, and 2.0.0.0 is the public IP address of the backend server. Creating a split tunnel port rule with "src_port" as 80 and "dest_port" as 9000 will make 1.0.0.0:80 forward to 2.0.0.0:9000. And if "dest_port" is empty, then 1.0.0.0:80 will simply forward to 2.0.0.0:80

    * If you want `dest_port` to be the same port as `src_port`, it is recommended to keep `dest_port` empty (i.e. don't set it to a value equal to that of "src_port"). This way you're allowing Tunnel Manager to handle it properly for you [especially if `src_port` is a port range]

  * `route_all_traffic_through_tunnel`: Whether to route all the traffic on the backend server through the tunnel. This is ignored on the tunnel host mode and only applies to the backend server. **Note that this can be `true` only on ONE tunnel!** You can't have more than a tunnel with `route_all_traffic_through_tunnel` set as `true`.

  * `dynamic_ip_updater_key`: The secret key (and also the key that identifies each tunnel) used for dynamic IP updates. This key is used to communicate between the Tunnel Manager instance hosted on the tunnel host, and the instance hosted on the backend server, for the purpose of updating the dynamic IP [in case a backend server is configured as "DYNAMIC"]. Make sure to keep `dynamic_ip_updater_key` a secret, **and make sure to set the same key on the configuration files of both the tunnel host and the backend server. [This has to be unique for each configured tunnel.]**

  * `wg_private_key_file_path`: The path to the WireGuard private key file. (works only when `driver` is set to `wireguard`).

  * `wg_server_tunnel_host_listen_port`: The port to make the WireGuard server start on [on the tunnel host]. (works only when `driver` is set to `wireguard`).

  * `wg_server_backend_server_listen_port`: The port to make the WireGuard server start on [on the backend server]. (works only when `driver` is set to `wireguard`).

  * `wg_tunnel_host_public_key`: The WireGuard public key of the tunnel host. (works only when `driver` is set to `wireguard`).

  * `wg_backend_server_public_key`: The WireGuard public key of the backend server. (works only when `driver` is set to `wireguard`).

## 🛠️ Installation as a service

**On both the tunnel host and the backend server(s):**

1. Store your configuration file at `/etc/tunmanager/config.json`

   You can copy the example configuration file and change it to serve your needs.
2. Place the binary file of Tunnel Manager at `/usr/local/bin` (e.g. `/usr/local/bin/tunmanager`)
3. Make the binary file executable: `chmod u+x /usr/local/bin/tunmanager`
4. Create a systemd service for Tunnel Manager. This can be done by creating `/etc/systemd/system/tunmanager.service` to have this content:
```
[Unit]
Description=TUNManager
After=network.target

[Service]
User=root
WorkingDirectory=/usr/local/bin
LimitNOFILE=2097152
TasksMax=infinity
ExecStart=/usr/local/bin/tunmanager /etc/tunmanager/config.json
Restart=on-failure
StartLimitInterval=180
StartLimitBurst=30
RestartSec=5s

[Install]
WantedBy=multi-user.target
```
5. Enable the Tunnel Manager service on startup & start it now:
```
systemctl enable --now tunmanager.service
```

## 💡 Example configuration case scenarios

### Protect a backend server behind a BuyVM DDoS-protected IP using a GRE tunnel

**Tunnel Host configuration**:
```json
{
    "mode": "tunnel_host",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "[buyvm non-ddos protected ip]",
            "tunnel_host_public_ip": "[buyvm ddos protected ip]",
            "backend_server_public_ip": "[backend public ip]",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "[buyvm non-ddos protected ip]",
            "tunnel_host_public_ip": "[buyvm ddos protected ip]",
            "backend_server_public_ip": "[backend public ip]",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```

### Two basic full GRE tunnels

**Tunnel Host configuration**:
```json
{
    "mode": "tunnel_host",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        },
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.3",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 2,
            "tunnel_interface_name": "tun2",
            "tunnel_rttables_id": 200,
            "tunnel_rttables_name": "TUN2",
            "tunnel_gateway_ip": "192.168.169.0",
            "tunnel_host_tunnel_ip": "192.168.169.1",
            "backend_server_tunnel_ip": "192.168.169.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers124@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        },
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.3",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 2,
            "tunnel_interface_name": "tun2",
            "tunnel_rttables_id": 200,
            "tunnel_rttables_name": "TUN2",
            "tunnel_gateway_ip": "192.168.169.0",
            "tunnel_host_tunnel_ip": "192.168.169.1",
            "backend_server_tunnel_ip": "192.168.169.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers124@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```

### A basic full GRE tunnel with a dynamic-IP backend server

**Tunnel Host configuration**:
```json
{
    "mode": "tunnel_host",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": true,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "DYNAMIC",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "apply_kernel_tuning_tweaks": false,
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "timeouts": {
        "ping_timeout": 5,
        "ping_interval": 10,
        "dynamic_ip_update_timeout": 30,
        "dynamic_ip_update_attempt_interval": 3
    },
    "tunnels": [
        {
            "driver": "gre",
            "tunnel_host_main_public_ip": "156.0.1.1",
            "tunnel_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "DYNAMIC",
            "tunnel_key": 1,
            "tunnel_interface_name": "tun1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "tunnel_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "src_port": "80",
                    "dest_port": ""
                },
                {
                    "proto": "TCP",
                    "src_port": "81",
                    "dest_port": "9000"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@",
            "wg_private_key_file_path": "/etc/tunmanager/wg-private",
            "wg_server_tunnel_host_listen_port": 51820,
            "wg_server_backend_server_listen_port": 51820,
            "wg_tunnel_host_public_key": "",
            "wg_backend_server_public_key": ""
        }
    ]
}
```