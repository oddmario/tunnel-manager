package utils

import "github.com/vishvananda/netlink"

func DoesNetworkInterfaceExist(if_name string) bool {
	var foundInterface bool = false

	addrs, _ := netlink.AddrList(nil, 0)

	for _, if_ := range addrs {
		if len(if_.Label) > 0 && if_.Label == if_name {
			foundInterface = true

			break
		}
	}

	return foundInterface
}
