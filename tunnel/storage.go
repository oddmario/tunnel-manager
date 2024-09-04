package tunnel

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/oddmario/gre-manager/utils"
	"github.com/tidwall/gjson"
)

func getStoragePath() string {
	path, _ := filepath.Abs("./")
	storagePath := filepath.Join(path, ".gre-manager")

	return storagePath
}

func doesStorageExist() bool {
	path := getStoragePath()

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func InitStorage() bool {
	path := getStoragePath()

	if !doesStorageExist() {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return false
		}
	}

	storageCfgPath := filepath.Join(getStoragePath(), "config.json")

	if _, err := os.Stat(storageCfgPath); !errors.Is(err, os.ErrNotExist) {
		cfg_content, _ := os.ReadFile(storageCfgPath)
		cfgContentString := utils.BytesToString(cfg_content)

		if gjson.Valid(cfgContentString) {
			cfg := gjson.Parse(cfgContentString)
			tuns := TunsFromJson(cfg.Get("tunnels"))

			for _, tun := range tuns {
				tun.Deinit(cfg.Get("mode").String(), cfg.Get("main_network_interface").String(), true)
			}
		}

		os.Remove(storageCfgPath)
	}

	origCfgPath, _ := filepath.Abs("./config.json")

	utils.CopyFile(origCfgPath, storageCfgPath)

	return true
}

func DestroyStorage(tuns []*Tunnel, mode, main_network_interface string) bool {
	if !doesStorageExist() {
		return false
	}

	storageCfgPath := filepath.Join(getStoragePath(), "config.json")

	for _, tun := range tuns {
		if tun.IsInitialised {
			tun.Deinit(mode, main_network_interface, false)
		}
	}

	os.Remove(storageCfgPath)

	return true
}