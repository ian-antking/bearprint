package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Config struct {
    ServerHost string
		ServerPort string
}

func NewConfig() (Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return Config{}, fmt.Errorf("cannot find home directory: %w", err)
    }

    configPath := filepath.Join(home, ".bearprint", "config")
    iniFile, err := ini.Load(configPath)
    if err != nil {
        return Config{}, fmt.Errorf("failed to load config file: %w", err)
    }

    host := iniFile.Section("default").Key("server_host").String()
    port := iniFile.Section("default").Key("server_port").String()

    if host == "" || port == "" {
        return Config{}, fmt.Errorf("missing server_host or server_port in config")
    }

    cfg := Config{
        ServerHost: host,
        ServerPort: port,
    }

    return cfg, nil
}
