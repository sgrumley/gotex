package config

// TODO: write a package to help read conf from ~/.config/go-tester/conf.yaml
// read from env var but default to .config
import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type EnvVar struct {
	ConfigFilePath string `envconfig:"CONFIG_FILE_PATH" default:"~/.config/tester"`
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
func GetConfig(filepath string) (Config, error) {
	configLocation := ""
	if filepath != "" {
		configLocation = filepath
	} else {
		e := EnvVar{}
		err := envconfig.Process("", &e)
		if err != nil {
			return Config{}, fmt.Errorf("failed to find config, proceeding with default %w", err)
		}
		configLocation = e.ConfigFilePath
	}

	cfg, err := LoadYAML[Config](configLocation)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config at path: %s with error: %w", configLocation, err)
	}
	fmt.Printf("\nconfig: %+v\n\n", cfg)
	return *cfg, nil
}

func LoadYAML[T any](path string) (*T, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading file with error: %w", err)
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
