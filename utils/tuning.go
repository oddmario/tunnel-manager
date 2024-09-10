package utils

func SysTuning(mode, mainNetworkInterface string) {
	Cmd("modprobe ip_gre", true)

	if mode == "tunnel_host" {
		Cmd("sysctl -w net.ipv4.ip_forward=1", true)
		Cmd("sysctl -w net.ipv4.conf."+mainNetworkInterface+".proxy_arp=1", true)
		Cmd("sysctl -w net.ipv4.conf.all.rp_filter=0", true)
		Cmd("sysctl -w net.ipv4.conf.default.rp_filter=0", true)
		Cmd("sysctl -w net.ipv4.conf.all.accept_redirects=0", true)
		Cmd("sysctl -w net.ipv4.conf.default.accept_redirects=0", true)
		Cmd("sysctl -w net.ipv4.route.flush=1", true)
		Cmd("sysctl -w net.ipv6.route.flush=1", true)
		Cmd("sysctl -w net.ipv4.tcp_mtu_probing=1", true)

		Cmd("sysctl -w fs.file-max=2097152", true)
		Cmd("sysctl -w fs.inotify.max_user_instances=2097152", true)
		Cmd("sysctl -w fs.inotify.max_user_watches=2097152", true)
		Cmd("sysctl -w fs.nr_open=2097152", true)
		Cmd("sysctl -w fs.aio-max-nr=2097152", true)
		Cmd("sysctl -w net.ipv4.tcp_syncookies=1", true)
		Cmd("sysctl -w net.core.somaxconn=65535", true)
		Cmd("sysctl -w net.core.netdev_max_backlog=65535", true)
		Cmd("sysctl -w net.core.dev_weight=128", true)
		Cmd("sysctl -w net.ipv4.ip_local_port_range=\"1024 65535\"", true)
		Cmd("sysctl -w net.nf_conntrack_max=1000000", true)
		Cmd("sysctl -w net.netfilter.nf_conntrack_max=1000000", true)
		Cmd("sysctl -w net.ipv4.tcp_max_tw_buckets=1440000", true)
		Cmd("sysctl -w net.ipv4.tcp_congestion_control=bbr", true)
		Cmd("sysctl -w net.core.default_qdisc=fq_codel", true)

		Cmd("modprobe tcp_bbr", true)
		Cmd("tc qdisc replace dev "+mainNetworkInterface+" root fq_codel", true)
		Cmd("ip link set "+mainNetworkInterface+" txqueuelen 15000", true)
		Cmd("ethtool -K "+mainNetworkInterface+" gro off gso off tso off", true)

		Cmd("iptables-nft -F", true)
	}
}
