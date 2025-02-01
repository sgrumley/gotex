package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/logging"
	"github.com/sgrumley/gotex/pkgv2/scanner"
)

func main() {
	logr, err := logger.New(
		logger.WithLevel(slog.LevelDebug),
		logger.WithSource(false),
	)
	if err != nil {
		log.Fatal("failed setting logger")
		fmt.Println("error: ", err.Error())
	}
	root, err := scanner.FindGoProjectRoot()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.GetConfig(logr)
	if err != nil {
		log.Fatal("failed to load a config: %w", err)
	}

	p, err := scanner.Scan(slog.Default(), cfg, root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dir: ", p.RootDir)
	p.Tree.Print()
}
