package finder

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sgrumley/gotex/pkg/config"
	"sgrumley/gotex/pkg/runner"
	"strings"
)

type Node interface {
	GetName() string
	GetChildren() []Node
	HasChildren() bool
	RunTest() (string, error)
}

var _ Node = (*Project)(nil)

type Project struct {
	Config   config.Config
	RootDir  string
	Packages []*Package
	log      *slog.Logger
}

func (p *Project) GetName() string {
	paths := strings.Split(p.RootDir, "/")

	return paths[len(paths)-1]
}

func (p *Project) GetChildren() []Node {
	children := make([]Node, 0)
	for _, child := range p.Packages {
		children = append(children, child)
	}
	return children
}

func (p *Project) HasChildren() bool {
	if len(p.Packages) > 0 {
		return true
	}

	return false
}

func (p *Project) RunTest() (string, error) {
	path := filepath.Dir(p.RootDir)
	return runner.RunTest(runner.TEST_TYPE_PROJECT, "", path, p.Config)
}

func InitProject(log *slog.Logger) (*Project, error) {
	cfg, err := config.GetConfig(log)
	if err != nil {
		return nil, fmt.Errorf("failed to load a config: %w", err)
	}

	p := &Project{
		Config: cfg,
		log:    log,
	}
	projectRoot, err := FindGoProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %s\n", err)
	}

	pkgs, err := FindPackages()
	if err != nil {
		return nil, err
	}

	// PERF: this could be concurrent
	for i := range pkgs {
		pkgs[i].Parent = p
		for _, file := range pkgs[i].Files {
			log.Info("searching file: ",
				slog.String("file", file.Path),
			)
			file.Functions = make([]*Function, 0)
			file.FunctionMap = make(map[string]*Function)
			file.Parent = pkgs[i]

			err := SearchFile(file, log)
			if err != nil {
				return nil, fmt.Errorf("failed finding tests: %s\n", err)
			}
		}
	}
	p.RootDir = projectRoot
	p.Packages = pkgs

	log.Info("project starting data",
		slog.String("root dir", p.RootDir),
	)

	return p, nil
}
