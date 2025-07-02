package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ian-antking/bear-print/bearprint-cli/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	originalUserHomeDir := config.OsUserHomeDir
	defer func() { config.OsUserHomeDir = originalUserHomeDir }()

	t.Run("flags provided", func(t *testing.T) {
		cfg, err := config.NewConfig("127.0.0.1", "8080")
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1", cfg.ServerHost)
		require.Equal(t, "8080", cfg.ServerPort)
	})

	t.Run("flags missing, config file present", func(t *testing.T) {
		tmpDir := t.TempDir()

		configDir := filepath.Join(tmpDir, ".bearprint")
		err := os.Mkdir(configDir, 0755)
		require.NoError(t, err)

		confFilePath := filepath.Join(configDir, "config")
		err = os.WriteFile(confFilePath, []byte(`
[default]
server_host=192.168.1.1
server_port=9090
`), 0644)
		require.NoError(t, err)

		config.OsUserHomeDir = func() (string, error) {
			return tmpDir, nil
		}

		cfg, err := config.NewConfig("", "")
		require.NoError(t, err)
		require.Equal(t, "192.168.1.1", cfg.ServerHost)
		require.Equal(t, "9090", cfg.ServerPort)
	})

	t.Run("missing flags and config file", func(t *testing.T) {
		tmpDir := t.TempDir()

		config.OsUserHomeDir = func() (string, error) {
			return tmpDir, nil
		}

		_, err := config.NewConfig("", "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load config file")
	})

	t.Run("error getting home dir", func(t *testing.T) {
		config.OsUserHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		_, err := config.NewConfig("", "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot find home directory")
	})
}
