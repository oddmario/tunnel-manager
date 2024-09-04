# 🚇 GRE Manager

An easy to use GRE tunnels manager written in Go (Golang).

## ✨ Features
- Dynamic IP addresses support for the backend server(s)
- Full & split tunnels support
- Port forwarding configuration for split tunnels
- Easy multi-tunnels management, all in a single configuration file

## 📝 Notes
For the optimal experience, kindly have a look at the notes listed at https://github.com/oddmario/GRE-setup-guide/blob/f0681a21edbc7a99f0d2a798529529a807357b5d/README.md#notes (you mainly need to pay attention to notes 1 to 5. anything that follows the 5th note can be ignored)

## 📙 Glossary
- **Backend server:** It is the server that you are trying to hide/protect the IP address of.
- **GRE host:** It is the VPS (or server in general) that has the public IP address that you want to expose publicly instead of the IP of the destination server. (e.g. a BuyVM VPS)

## 🧐 Configuration documentation
- `mode`
  * Can be either `backend_server` or `gre_host`.

- `main_network_interface`
  * The name of the main network interface on the system (e.g. eth0)

- `dynamic_ip_updater_api`
  * `is_enabled`: Whether to enable the dynamic IP updater API or no. **This can be enabled only on the GRE host mode.** Note that you have to enable this if you have any tunnels with a **DYNAMIC** `backend_server_public_ip`. (this config parameter is ignored on the backend server mode)

  * `listen_address`: The IP address that will be used for the dynamic IP updater HTTP server. Make sure that it's binding to an IP address that the backend server(s) can access. (this config parameter is ignored on the backend server mode)

  * `listen_port`: The port that will be used for the dynamic IP updater HTTP server. (**This is not ignored on the backend server mode!** Note that if you are configuring the GRE-Manager copy of a backend server, you need to specify this to be the same `listen_port` configured on the GRE Host configuration file)

- `tunnels`: An array of the tunnel(s) that you would like to have.
  * `gre_host_main_public_ip`: The main/primary public IP address of the GRE host.

  * `gre_host_public_ip`: The public IP address of the GRE host that you would like to use instead of the backend server IP address. If your GRE host has only one public IP address, make **gre_host_main_public_ip** and **gre_host_public_ip** equal. When they are equal, you will use your single public IP address of the GRE host for the tunneling.

  * `backend_server_public_ip`: The public IP address of the backend server. Set this to `DYNAMIC` if the backend server has a dynamic IP address.

  * `tunnel_key`: The index/key of the GRE tunnel. This has to be unique for each configured tunnel. (e.g. 1, 2, 3, etc). **It also has to be the same configured value for the tunnel on both the GRE host configuration file & the backend one.**

  * `tunnel_interface_name`: The name of the GRE tunnel interface. This has to be unique for each configured tunnel. (e.g. gre1, gre2, gre3, etc)

  * `tunnel_rttables_id`: The ID of the routing table used by the GRE tunnel. This has to be unique for each configured tunnel. (e.g. 100, 200, 300, etc). **[This is ignored on the GRE host as it is used only by the backend server instance. So it doesn't matter what value you set for `tunnel_rttables_id` on the GRE host instance of GRE Manager]**

  * `tunnel_rttables_name`: The name of the routing table used by the GRE tunnel. This has to be unique for each configured tunnel. (e.g. TUN1, TUN2, TUN3, etc). **[This is ignored on the GRE host as it is used only by the backend server instance. So it doesn't matter what value you set for `tunnel_rttables_name` on the GRE host instance of GRE Manager]**

  * `tunnel_gateway_ip`: The gateway that will be used by GRE Manager to setup the GRE tunnel.

  * `gre_host_tunnel_ip`: The IP address of the GRE host inside the tunnel.

  * `backend_server_tunnel_ip`: The IP address of the backend server inside the tunnel.

  * `tunnel_type`: Can be either **split** for a split tunnel, or **full** for a full tunnel. A full tunnel forwards all the ports, meanwhile a split tunnel forwards certain ports that you can configure in `split_tunnel_ports`. **[This is ignored on the backend server as it is used only by the GRE host instance. So it doesn't matter what value you set for `tunnel_type` on a backend server instance of GRE Manager]**

  * `split_tunnel_ports`: An array containing the ports to forward for the purpose of split tunneling. This is ignored if `tunnel_type` is set to "full"
    * `proto`: Can be either TCP or UDP
    * `port`: The port(s) to forward. To use a port range, you can use the `start_port:end_port` format (e.g. `8000:8050`).

  * `route_all_traffic_through_tunnel`: Whether to route all the traffic on the backend server through the GRE tunnel. This is ignored on the GRE host mode and only applies to the backend server. **Note that this can be `true` only on ONE tunnel!** You can't have more than a tunnel with `route_all_traffic_through_tunnel` set as `true`.

  * `dynamic_ip_updater_key`: The secret key (and also the key that identifies each tunnel) used for dynamic IP updates. This key is used to communicate between the GRE Manager instance hosted on the GRE host, and the instance hosted on the backend server, for the purpose of updating the dynamic IP [in case a backend server is configured as "DYNAMIC"]. Make sure to keep `dynamic_ip_updater_key` a secret, **and make sure to set the same key on the configuration files of both the GRE host and the backend server. [This has to be unique for each configured tunnel.]**

## 🛠️ Installation as a service

**On both the GRE host and the backend server(s):**

1. Store your configuration file at `/etc/gremanager/config.json`

   You can copy the example configuration file and change it to serve your needs.
2. Place the binary file of GRE Manager at `/usr/local/bin` (e.g. `/usr/local/bin/gremanager`)
3. Make the binary file executable: `chmod u+x /usr/local/bin/gremanager`
4. Create a systemd service for GRE Manager. This can be done by creating `/etc/systemd/system/gremanager.service` to have this content:
```
[Unit]
Description=GREManager
After=network.target

[Service]
User=root
WorkingDirectory=/usr/local/bin
LimitNOFILE=2097152
TasksMax=infinity
ExecStart=/usr/local/bin/gremanager /etc/gremanager/config.json
Restart=on-failure
StartLimitInterval=180
StartLimitBurst=30
RestartSec=5s

[Install]
WantedBy=multi-user.target
```
5. Enable the GRE Manager service on startup & start it now:
```
systemctl enable --now gremanager.service
```

## 💡 Example configuration case scenarios

### Protect a backend server behind a BuyVM DDoS-protected IP

**GRE Host configuration**:
```json
{
    "mode": "gre_host",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "[buyvm non-ddos protected ip]",
            "gre_host_public_ip": "[buyvm ddos protected ip]",
            "backend_server_public_ip": "[backend public ip]",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "[buyvm non-ddos protected ip]",
            "gre_host_public_ip": "[buyvm ddos protected ip]",
            "backend_server_public_ip": "[backend public ip]",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        }
    ]
}
```

### Two basic full tunnels

**GRE Host configuration**:
```json
{
    "mode": "gre_host",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        },
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.3",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 2,
            "tunnel_interface_name": "gre2",
            "tunnel_rttables_id": 200,
            "tunnel_rttables_name": "TUN2",
            "tunnel_gateway_ip": "192.168.169.0",
            "gre_host_tunnel_ip": "192.168.169.1",
            "backend_server_tunnel_ip": "192.168.169.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers124@"
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        },
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.3",
            "backend_server_public_ip": "156.0.0.2",
            "tunnel_key": 2,
            "tunnel_interface_name": "gre2",
            "tunnel_rttables_id": 200,
            "tunnel_rttables_name": "TUN2",
            "tunnel_gateway_ip": "192.168.169.0",
            "gre_host_tunnel_ip": "192.168.169.1",
            "backend_server_tunnel_ip": "192.168.169.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers124@"
        }
    ]
}
```

### A basic full tunnel with a dynamic-IP backend server

**GRE Host configuration**:
```json
{
    "mode": "gre_host",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": true,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "DYNAMIC",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        }
    ]
}
```

**Backend server configuration**:
```json
{
    "mode": "backend_server",
    "main_network_interface": "eth0",
    "dynamic_ip_updater_api": {
        "is_enabled": false,
        "listen_address": "0.0.0.0",
        "listen_port": 30100
    },
    "tunnels": [
        {
            "gre_host_main_public_ip": "156.0.1.1",
            "gre_host_public_ip": "156.0.1.2",
            "backend_server_public_ip": "DYNAMIC",
            "tunnel_key": 1,
            "tunnel_interface_name": "gre1",
            "tunnel_rttables_id": 100,
            "tunnel_rttables_name": "TUN1",
            "tunnel_gateway_ip": "192.168.168.0",
            "gre_host_tunnel_ip": "192.168.168.1",
            "backend_server_tunnel_ip": "192.168.168.2",
            "tunnel_type": "full",
            "split_tunnel_ports": [
                {
                    "proto": "TCP",
                    "port": "80"
                }
            ],
            "route_all_traffic_through_tunnel": false,
            "dynamic_ip_updater_key": "wowsers123@"
        }
    ]
}
```