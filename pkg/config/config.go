// Package config provides configuration management for the gotex library.
// It handles loading configuration from environment variables, YAML files,
// and default embedded configurations. Configuration determines test execution
// behavior including output formatting, timeouts, and test filtering options.
//
// The configuration system follows this precedence order:
//  1. Environment variable specified file path (GOTEX_CONFIG_FILE_PATH)
//  2. Default user config file (~/.config/gotex/config.yaml)
//  3. Embedded default configuration
//
// Example usage:
//
//	cfg, err := config.GetConfig(ctx)
//	if err != nil {
//		log.Fatal("Failed to load config:", err)
//	}
//	fmt.Printf("Verbose mode: %t\n", cfg.Verbose)
package config

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sgrumley/gotex/pkg/path"
	"github.com/sgrumley/gotex/pkg/slogger"
	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var embedCfg embed.FS

var (
	defaultCfg      *Config
	defaultFilePath = "~/.config/gotgithub.com/sgrumley/gotex/pkg/config.yaml"
)

// EnvVar represents environment variable configuration options.
type EnvVar struct {
	ConfigFilePath string `envconfig:"GOTEX_CONFIG_FILE_PATH"`
}

// Config represents the test execution configuration.
// It controls various aspects of test running including output format,
// execution timeouts, and filtering options.
type Config struct {
	PipeTo   string `yaml:"pipeto"`   // External command to pipe test output to (e.g., "tparse")
	Timeout  string `yaml:"timeout"`  // Test execution timeout duration
	Json     bool   `yaml:"json"`     // Enable JSON output format
	Short    bool   `yaml:"short"`    // Use short test output format
	Verbose  bool   `yaml:"verbose"`  // Enable verbose test output
	FailFast bool   `yaml:"failfast"` // Stop on first test failure
	Cover    bool   `yaml:"cover"`    // Enable test coverage reporting
}

// GetConfig loads and returns the application configuration.
//
// It follows this precedence order:
//   - Check env var for config file path
//   - Otherwise check ~/.config/gotex/config.yaml
//   - Otherwise use the default yaml file embedded in the package
//
// The function will attempt to use default configuration if user-specified
// configuration cannot be loaded, logging warnings appropriately.

func GetConfig(ctx context.Context) (Config, error) {
	log, err := slogger.FromContext(ctx)
	if err != nil {
		return Config{}, err
	}
	filepath, err := GetConfigPath()
	if err != nil {
		cfg, errDefault := getDefaultCfg()
		if errDefault != nil {
			return Config{}, errDefault
		}
		log.Warn("failed to find user specified config, using default",
			slog.String("cause", err.Error()),
			slog.Any("default config", cfg),
		)
		return cfg, nil
	}

	cfg, err := LoadYAML(filepath)
	if err != nil {
		cfg, errDefault := getDefaultCfg()
		if errDefault != nil {
			return Config{}, errDefault
		}
		log.Error("invalid config file",
			fmt.Errorf("failed to load config at path: %s with error: %w", filepath, err),
			slog.Any("default config", cfg),
		)
		return cfg, err
	}

	log.Info("loaded user config from environment variable", slog.Any("config", cfg))
	return cfg, nil
}

// FileExists checks if a file exists at the specified path.
// It handles home directory expansion (~) and returns false if the file
// doesn't exist or if there's an error accessing it.
func FileExists(filepath string) bool {
	fp, err := path.ReplaceHomeDirChar(filepath)
	if err != nil {
		return false
	}
	_, err = os.Stat(fp)
	return !errors.Is(err, fs.ErrNotExist)
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
	if err != nil || e.ConfigFilePath == "" {
		return defaultFilePath, fmt.Errorf("no environment variable found")
	}

	// check that the user specified exists
	if !FileExists(e.ConfigFilePath) {
		return defaultFilePath, fmt.Errorf("failed to load config from environment variable, attempted path: %s", e.ConfigFilePath)
	}

	return e.ConfigFilePath, nil
}

// LoadYAML loads configuration from a YAML file at the specified path.
// It handles home directory expansion and parses the YAML content.
func LoadYAML(filePath string) (Config, error) {
	fp, err := path.ReplaceHomeDirChar(filePath)
	if err != nil {
		return Config{}, err
	}
	b, err := os.ReadFile(fp)
	if err != nil {
		return Config{}, err
	}
	return LoadConfig(b)
}

// LoadConfig parses configuration from YAML byte data.
func LoadConfig(b []byte) (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return cfg, nil
}

// deprecated
func (c *Config) String() string {
	return fmt.Sprintf(`
		Config:
			json: %t
			timeout: %s
			short: %t
			verbose: %t
			fail fast: %t
			cover: %t
			[green]pipe to[-]: %s
		`, c.Json, c.Timeout, c.Short, c.Verbose, c.FailFast, c.Cover, c.PipeTo)
}
