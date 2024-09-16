package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func GetExternalIP(tunnel_host_main_public_ip string, dynamic_ip_update_timeout, dynamic_ip_updater_api_listen_port int) (string, error) {
	req, _ := resty.New().SetTimeout(time.Duration(dynamic_ip_update_timeout) * time.Second).R().
		Get("http://" + tunnel_host_main_public_ip + ":" + IToStr(dynamic_ip_updater_api_listen_port) + "/get_pub_ip")

	if req.StatusCode() != 200 {
		return "", errors.New("unable to get public ip")
	}

	return strings.TrimSpace(req.String()), nil
}
