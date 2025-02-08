package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sgrumley/gotex/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: these tests would be better in a container
func Test_SuccessGetConfig(t *testing.T) {
	tcs := []struct {
		name    string
		setup   func() // TODO: this needs to be extended to store the state of anything used and restore it after test
		assert  func(actual config.Config, err error)
		cleanup func() // TODO: this should clean temp files and reset env vars to previous state
	}{
		{
			name: "success_user_config",
			setup: func() {
				// os.Setenv("GOTEX_CONFIG_FILE_PATH", "./testdata/test_config.yaml")
			},
			assert: func(actual config.Config, err error) {
				expected := config.Config{
					PipeTo:   "tparse",
					Timeout:  "",
					Json:     true,
					Short:    false,
					Verbose:  true,
					FailFast: false,
					Cover:    true,
				}

				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		// {
		// 	//NOTE: not sure how to go about this test since it requires the file to be in os path
		// 	name: "success_default_config_path",
		// 	setup: func() {
		// 		os.Setenv("GOTEX_CONFIG_FILE_PATH", "")
		// 	},
		// 	assert: func(actual config.Config, err error) {
		// 		expected := config.Config{
		// 			PipeTo:   "",
		// 			Timeout:  "",
		// 			Json:     true,
		// 			Short:    false,
		// 			Verbose:  true,
		// 			FailFast: false,
		// 			Cover:    true,
		// 		}

		// 		assert.NoError(t, err)
		// 		assert.Equal(t, actual, expected)
		// 	},

		// },
		// {
		// 	name: "success_default_config",
		// 	setup: func() {
		// 		os.Setenv("GOTEX_CONFIG_FILE_PATH", "")
		// 	},
		// 	assert: func(actual config.Config, err error) {
		// 		expected := config.Config{
		// 			PipeTo:   "",
		// 			Timeout:  "",
		// 			Json:     false,
		// 			Short:    false,
		// 			Verbose:  false,
		// 			FailFast: false,
		// 			Cover:    false,
		// 		}

		// 		assert.NoError(t, err)
		// 		assert.Equal(t, actual, expected)
		// 	},
		// },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			cfg, err := config.GetConfig(context.Background())
			tc.assert(cfg, err)
		})
	}
}

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
			input:    "~/",
			expected: filepath.Join(mockHomeDir),
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := config.ReplaceHomeDirChar(tt.input)
			require.NoError(t, err)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
