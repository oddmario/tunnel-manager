package workers

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

type monitorInterfaceStateChangeCallback func(string)

func MonitorInterfaceStateChange(if_name string, callback monitorInterfaceStateChangeCallback) {
	linkChan := make(chan netlink.LinkUpdate)
	if err := netlink.LinkSubscribe(linkChan, nil); err != nil {
		fmt.Printf("[WARN] Failed to subscribe to link updates: %v", err)

		return
	}

	for update := range linkChan {
		if update.Link != nil {
			if update.Link.Attrs().Name == if_name {
				if update.Link.Attrs().Flags&net.FlagUp != 0 {
					callback("UP")
				} else {
					callback("DOWN")
				}
			}
		}
	}
}
