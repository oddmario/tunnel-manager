package main

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/oddmario/tunnel-manager/config"
	"github.com/oddmario/tunnel-manager/dynamicipupdater"
	"github.com/oddmario/tunnel-manager/logger"
	"github.com/oddmario/tunnel-manager/tunnel"
	"github.com/oddmario/tunnel-manager/utils"
	"github.com/oddmario/tunnel-manager/vars"
)

func main() {
	logger.Init()

	if runtime.GOOS != "linux" {
		logger.Fatal("Sorry! Tunnel Manager can only run on Linux systems.")
	}

	args := os.Args[1:]

	if len(args) >= 1 {
		vars.ConfigFilePath = args[0]
	} else {
		vars.ConfigFilePath, _ = filepath.Abs("./config.json")
	}

	if _, err := os.Stat(vars.ConfigFilePath); errors.Is(err, os.ErrNotExist) {
		logger.Fatal("The specified configuration file does not exist.")
	}

	logger.Info("Starting Tunnel Manager v" + vars.Version + "...")

	config.LoadConfig()
	tunnel.InitStorage()

	defer func() {
		tunnel.DestroyStorage(config.Config.Tunnels, config.Config.Mode, config.Config.MainNetworkInterface)
	}()

	var shouldEnableIPIPmod bool = false
	var shouldEnableGREmod bool = false
	var shouldEnableWGmod bool = false

	for _, tun := range config.Config.Tunnels {
		if tun.TunnelDriver == "gre" && !shouldEnableGREmod {
			shouldEnableGREmod = true
		}
		if tun.TunnelDriver == "ipip" && !shouldEnableIPIPmod {
			shouldEnableIPIPmod = true
		}
		if tun.TunnelDriver == "wireguard" && !shouldEnableWGmod {
			shouldEnableWGmod = true
		}
	}

	utils.SysTuning(shouldEnableIPIPmod, shouldEnableGREmod, shouldEnableWGmod, config.Config.Mode, config.Config.MainNetworkInterface, config.Config.ApplyKernelTuningTweaks)

	if config.Config.DynamicIPUpdaterAPIIsEnabled {
		if config.Config.Mode == "tunnel_host" {
			go dynamicipupdater.InitServer()
		} else {
			logger.Warn("The dynamic IP updater API is meant to be enabled only on the tunnel host. Ignoring `dynamic_ip_updater_api.is_enabled`...")
		}
	}

	for _, tun := range config.Config.Tunnels {
		tun.Init(config.Config.Mode, config.Config.MainNetworkInterface, config.Config.DynamicIPUpdaterAPIListenPort, config.Config.DynamicIPUpdateInterval, config.Config.DynamicIPUpdateTimeout, config.Config.PingInterval, config.Config.PingTimeout, config.Config.ApplyKernelTuningTweaks)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
