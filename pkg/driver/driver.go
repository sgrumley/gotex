package driver

import (
	"fmt"
	"log/slog"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/finder"
	"github.com/sgrumley/gotex/pkg/models"
)

func InitProject(log *slog.Logger) (*models.Project, error) {
	cfg, err := config.GetConfig(log)
	if err != nil {
		return nil, fmt.Errorf("failed to load a config: %w", err)
	}

	root, err := finder.FindGoProjectRoot()
	if err != nil {
		return nil, err
	}

	p, err := models.NewProject(log, cfg, root)
	if err != nil {
		return nil, err
	}

	return p, nil
}
