package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var OsUserHomeDir = os.UserHomeDir

type Config struct {
	ServerHost string
	ServerPort string
}

func NewConfig(hostFlag, portFlag string) (Config, error) {
	cfg := Config{
		ServerHost: hostFlag,
		ServerPort: portFlag,
	}

	if cfg.ServerHost == "" || cfg.ServerPort == "" {
		home, err := OsUserHomeDir()
		if err != nil {
			return Config{}, fmt.Errorf("cannot find home directory: %w", err)
		}

		configPath := filepath.Join(home, ".bearprint", "config")
		iniFile, err := ini.Load(configPath)
		if err != nil {
			if cfg.ServerHost != "" && cfg.ServerPort != "" {
				return cfg, nil
			}
			return Config{}, fmt.Errorf("failed to load config file at %s: %w", configPath, err)
		}

		if cfg.ServerHost == "" {
			cfg.ServerHost = iniFile.Section("default").Key("server_host").String()
		}
		if cfg.ServerPort == "" {
			cfg.ServerPort = iniFile.Section("default").Key("server_port").String()
		}
	}

	if cfg.ServerHost == "" || cfg.ServerPort == "" {
		return Config{}, fmt.Errorf("missing configuration: please provide flags (-host, -port) or a config file")
	}

	return cfg, nil
}
