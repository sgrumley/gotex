package config

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var embedCfg embed.FS

var (
	defaultCfg      *Config
	defaultFilePath = "~/.config/gotex/config.yaml"
)

type EnvVar struct {
	ConfigFilePath string `envconfig:"GOTEX_CONFIG_FILE_PATH"` //  default:"~/.config/gotex/config.yaml"`
}

type Config struct {
	PipeTo   string `yaml:"pipeto"`
	Timeout  string `yaml:"timeout"`
	Json     bool   `yaml:"json"`
	Short    bool   `yaml:"short"`
	Verbose  bool   `yaml:"verbose"`
	FailFast bool   `yaml:"failfast"`
	Cover    bool   `yaml:"cover"`
}

/*
- check env var for config file path
- else it should check ~/.config/gotex/config.yaml
- else use the default yaml file in this folder
*/

func GetConfig(log *slog.Logger) (Config, error) {
	filepath, err := GetConfigPath()
	if err != nil {
		cfg, err := getDefaultCfg()
		log.Warn("failed to find user specified config, using default",
			slog.String("cause", err.Error()),
			slog.Any("default config", cfg),
		)
		return cfg, err
	}

	cfg, err := LoadYAML(filepath)
	if err != nil {
		cfg, err := getDefaultCfg()
		log.Error("invalid config file",
			slog.Any("error", fmt.Errorf("failed to load config at path: %s with error: %w", filepath, err)),
			slog.Any("default config", cfg),
		)
		return cfg, err
	}

	log.Info("loaded user config from environment variable", slog.Any("config", cfg))
	return cfg, nil
}

func FileExists(filepath string) bool {
	fp, err := ReplaceHomeDirChar(filepath)
	if err != nil {
		return false
	}
	_, err = os.Stat(fp)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

// getDefaultCfg returns the config of default config specified in the local file default.yaml
// it uses embed to allow for a nicer way to document a example config as opposed to an incode struct
func getDefaultCfg() (Config, error) {
	if defaultCfg != nil {
		return *defaultCfg, nil
	}

	b, err := embedCfg.ReadFile("default.yaml")
	if err != nil {
		return Config{}, err
	}

	cfg, err := LoadConfig(b)
	if err != nil {
		return Config{}, err
	}
	defaultCfg = &cfg

	return *defaultCfg, nil
}

func GetConfigPath() (string, error) {
	// check for user specified config path
	e := EnvVar{}
	err := envconfig.Process("GOTEX_CONFIG_FILE_PATH", &e)
	if err != nil {
		return defaultFilePath, fmt.Errorf("no environment variable found")
	}

	// check that the user specified exists
	if !FileExists(e.ConfigFilePath) {
		return defaultFilePath, fmt.Errorf("failed to load config from environment variable, attempted path: %s", e.ConfigFilePath)
	}

	return e.ConfigFilePath, nil
}

func LoadYAML(path string) (Config, error) {
	fp, err := ReplaceHomeDirChar(path)
	if err != nil {
		return Config{}, err
	}
	b, err := os.ReadFile(fp)
	if err != nil {
		return Config{}, err
	}
	return LoadConfig(b)
}

func LoadConfig(b []byte) (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return cfg, nil
}

func ReplaceHomeDirChar(fp string) (string, error) {
	if !strings.Contains(fp, "~") {
		return fp, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting home directory: %w", err)
	}
	// Replace ~ with the home directory path
	if fp[:2] == "~/" {
		fp = filepath.Join(homeDir, fp[2:])
	}
	return fp, nil
}
