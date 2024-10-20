package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"sgrumley/gotex/pkg/config"
)

func Test_GetConfig(t *testing.T) {}

func Test_getDefaultConfig(t *testing.T) {
	// ensure no filepath available
}

func Test_GetConfigPath(t *testing.T) {}

func Test_LoadConfig(t *testing.T) {}

// TEST: complete
func Test_LoadYAML(t *testing.T) {}

func Test_ReplaceHomeDirChar(t *testing.T) {
	mockHomeDir := "/mock/home"
	os.Setenv("HOME", mockHomeDir)

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "~/documents/file.txt",
			expected: filepath.Join(mockHomeDir, "documents/file.txt"),
		},
		{
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			input: "~/", 
			expected: filepath.Join(mockHomeDir),
		},
		{
			input: "", 
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := config.ReplaceHomeDirChar(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
