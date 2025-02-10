package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/scanner"
	"github.com/sgrumley/gotex/pkg/slogger"

	"github.com/sgrumley/gotex/internal/components"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	log, err := slogger.New(
		slogger.WithLevel(slog.LevelDebug),
		slogger.WithSource(false),
	)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return 1
	}
	ctx = slogger.AddToContext(ctx, log)

	root, err := scanner.FindGoProjectRoot()
	if err != nil {
		log.Fatal("No go project found, navigate to a repository with a go.mod file and try again", err)
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		fmt.Printf("failed to load a config: %s\n", err.Error())
		log.Fatal("failed to load a config", err)
	}

	app, err := components.New(ctx, cfg, root)
	if err != nil {
		fmt.Printf("failed to initialise project: %s", err.Error())
		return 1
	}
	err = app.Start(ctx)
	if err != nil {
		// to file
		log.Error("application crashed", err)
		// to stdout
		fmt.Printf("application crashed: %s", err.Error())
		return 1
	}

	return 0
}
