package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

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

func GetConfig() (Config, error) {
	filepath := GetConfigPath()
	if filepath == "" {
		return getDefaultCfg()
	}

	cfg, err := LoadYAML(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config at path: %s with error: %w", filepath, err)
	}
	fmt.Printf("\nconfig: %+v\n\n", cfg)
	return cfg, nil
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
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

func GetConfigPath() string {
	// check for user specified config path
	e := EnvVar{}
	err := envconfig.Process("GOTEX_CONFIG_FILE_PATH", &e)
	if err != nil {
		fmt.Printf("no config env found, proceeding with default path: %s\n", defaultFilePath)
		return defaultFilePath
	}

	// check that the user specified exists
	if FileExists(e.ConfigFilePath) {
		fmt.Printf("config file not found at specified location %s, proceeding with default config\n", e.ConfigFilePath)
		return e.ConfigFilePath
	}

	return ""
}

func LoadYAML(path string) (Config, error) {
	fp := ReplaceHomeDirChar(path)
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

func ReplaceHomeDirChar(fp string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}
	// Replace ~ with the home directory path
	if fp[:2] == "~/" {
		fp = filepath.Join(homeDir, fp[2:])
	}
	return fp
}
