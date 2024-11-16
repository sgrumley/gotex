package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ListTestFilesWithCWD() ([]*File, error) {
	rootDir, err := FindGoProjectRoot()
	if err != nil {
		return nil, err
	}

	files, err := listTestFilesWithPath(rootDir)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func listTestFilesWithPath(dirPath string) ([]*File, error) {
	var files []*File

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.Contains(path, "_test.go") {
			pathSplit := strings.Split(path, "/")
			files = append(files, &File{
				Name: pathSplit[len(pathSplit)-1],
				Path: path,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// TODO: add options so that config can determine if go.work is the root
func FindGoProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Walk backwards through directories until we find go.mod
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("go.mod file not found in current directory or any parent directories")
}
