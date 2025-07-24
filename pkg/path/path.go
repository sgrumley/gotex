// Package path provides file system path utilities for the gotex library.
// It includes functions for handling home directory expansion and path
// manipulation operations commonly needed in configuration and file handling.
//
// The package primarily focuses on:
//   - Home directory path expansion (~ character handling)
//   - Cross-platform path operations
//   - File system path validation and manipulation
//
// Example usage:
//
//	expandedPath, err := path.ReplaceHomeDirChar("~/config/gotex.yaml")
//	if err != nil {
//		log.Fatal("Path expansion failed:", err)
//	}
//	fmt.Printf("Full path: %s\n", expandedPath)
package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceHomeDirChar expands the tilde (~) character in file paths to the
// user's home directory. If the path doesn't start with ~, it returns the
// path unchanged. This function handles the common pattern of using ~ to
// represent the home directory in configuration paths.
//
// Examples:
//   - "~/config/app.yaml" → "/home/user/config/app.yaml"
//   - "/absolute/path" → "/absolute/path" (unchanged)
//   - "relative/path" → "relative/path" (unchanged)
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
