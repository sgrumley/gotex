package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
