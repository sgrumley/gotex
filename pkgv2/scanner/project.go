package scanner

import (
	"fmt"
	"log/slog"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/models"
)

func Scan(log *slog.Logger, cfg config.Config, root string) (*models.Project, error) {
	cfg, err := config.GetConfig(log)
	if err != nil {
		return nil, fmt.Errorf("failed to load a config: %w", err)
	}

	pkgs, err := FindPackages(root)
	if err != nil {
		return nil, err
	}

	p := &models.Project{
		Config:   cfg,
		Log:      log,
		RootDir:  root,
		Packages: pkgs,
	}

	err = PopulateFromPackages(p, pkgs)
	if err != nil {
		return nil, fmt.Errorf("Error building tree: %v\n", err)
	}

	return p, nil
}
