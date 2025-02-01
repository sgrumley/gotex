package finder

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

func FindPackages() ([]*models.Package, error) {
	// TODO: make this read from cfg
	blacklist := map[string]struct{}{
		"node_modules": {},
		".git":         {},
		"vendor":       {},
	}
	rootDir, err := FindGoProjectRoot()
	if err != nil {
		return nil, err
	}

	var packages []*models.Package

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if _, blocked := blacklist[info.Name()]; blocked {
				return filepath.SkipDir
			}
			pkg, err := build.ImportDir(path, 0)
			if err == nil && len(pkg.GoFiles) > 0 {

				filenames := findTestFiles(path)
				if len(filenames) == 0 {
					return nil
				}

				paths := strings.Split(pkg.Dir, "/")
				pkgName := fmt.Sprintf("%s/%s", paths[len(paths)-2], paths[len(paths)-1])

				tmp := &Package{
					Name: pkgName,
					Path: pkg.Dir,
				}

				files := make([]*File, 0)
				for i := range filenames {
					files = append(files, &File{
						Name:   filenames[i],
						Path:   fmt.Sprintf("%s/%s", pkg.Dir, filenames[i]),
						Parent: tmp,
					})
				}
				tmp.Files = files

				packages = append(packages, tmp)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return packages, nil
}
