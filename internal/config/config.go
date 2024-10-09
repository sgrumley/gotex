package config

// TODO: write a package to help read conf from ~/.config/go-tester/conf.yaml
// read from env var but default to .config
import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type EnvVar struct {
	ConfigFilePath string `envconfig:"GOTEX_CONFIG_FILE_PATH" default:"~/.config/gotex/config.yaml"`
	// ConfigFilePath string `envconfig:"GOTEX_CONFIG_FILE_PATH" default:"~/repo/go-tester/internal/config/default.yaml"`
	// ConfigFilePath string `envconfig:"GOTEX_CONFIG_FILE_PATH" default:"./internal/config/default.yaml"`
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

// Filepath param is temporary or moved to options pattern
func GetConfig(fp string) (Config, error) {
	configLocation := ""
	if fp != "" {
		configLocation = fp
	} else {
		e := EnvVar{}
		err := envconfig.Process("GOTEX_CONFIG_FILE_PATH", &e)
		if err != nil {
			// TODO: loading default is not an error
			fmt.Printf("failed to find config at %s, proceeding with default\n", e.ConfigFilePath)
		}
		// if err != nil {
		// 	return Config{}, fmt.Errorf("failed to find config, proceeding with default %w", err)
		// }
		// TODO: this needs to set default to the yaml in this folder
		configLocation  = e.ConfigFilePath
		fmt.Println(e.ConfigFilePath)
	}

	cfg, err := LoadYAML[Config](configLocation)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config at path: %s with error: %w", configLocation, err)
	}
	fmt.Printf("\nconfig: %+v\n\n", cfg)
	return *cfg, nil
}

func LoadYAML[T any](path string) (*T, error) {
	fp := replaceHomeDirChar(path)
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	return LoadConfig[T](bytes.NewReader(b))
}

func LoadConfig[T any](reader io.Reader) (*T, error) {
	var cfg T
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("the config yaml is invalid: %w", err)
	}

	return &cfg, nil
}

func replaceHomeDirChar(fp string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}
	// Replace `~` with the home directory path
	if fp[:2] == "~/" {
		fp = filepath.Join(homeDir, fp[2:])
	}
	return fp
}
