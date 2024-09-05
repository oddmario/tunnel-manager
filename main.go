package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/oddmario/gre-manager/config"
	"github.com/oddmario/gre-manager/constants"
	"github.com/oddmario/gre-manager/dynamicipupdater"
	"github.com/oddmario/gre-manager/tunnel"
	"github.com/oddmario/gre-manager/utils"
)

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("Sorry! GRE Manager can only run on Linux systems.")
	}

	args := os.Args[1:]

	if len(args) >= 1 {
		constants.ConfigFilePath = args[0]
	} else {
		constants.ConfigFilePath, _ = filepath.Abs("./config.json")
	}

	if _, err := os.Stat(constants.ConfigFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("The specified configuration file does not exist.")
	}

	fmt.Println("[INFO] Starting GRE Manager v" + constants.Version + "...")

	config.LoadConfig()
	tunnel.InitStorage()

	defer func() {
		tunnel.DestroyStorage(config.Config.Tunnels, config.Config.Mode, config.Config.MainNetworkInterface)
	}()

	utils.SysTuning(config.Config.Mode, config.Config.MainNetworkInterface)

	if config.Config.DynamicIPUpdaterAPIIsEnabled {
		if config.Config.Mode == "gre_host" {
			go dynamicipupdater.InitServer()
		} else {
			fmt.Println("[WARN] The dynamic IP updater API is meant to be enabled only on the GRE host. Ignoring `dynamic_ip_updater_api.is_enabled`...")
		}
	}

	for _, tun := range config.Config.Tunnels {
		tun.Init(config.Config.Mode, config.Config.MainNetworkInterface, config.Config.DynamicIPUpdaterAPIListenPort)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
