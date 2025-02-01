package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/scanner"

	"github.com/sgrumley/gotex/internal/components"
)

func main() {
	os.Exit(run())
}

func run() int {
	// add to ctx
	// log, err := logger.New(
	// 	logger.WithLevel(slog.LevelDebug),
	// 	logger.WithSource(false),
	// )
	// if err != nil {
	// 	fmt.Println("error: ", err.Error())
	// 	return 1
	// }

	root, err := scanner.FindGoProjectRoot()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatal("failed to load a config: %w", err)
	}

	app, err := components.New(ctx, cfg, root)
	if err != nil {
		fmt.Printf("failed to initialise project: %s", err.Error())
		return 1
	}
	err = app.Start()
	if err != nil {
		// to file
		// log.Error("application crashed", slog.Any("error", err))
		// to stdout
		fmt.Printf("application crashed: %s", err.Error())
		return 1
	}

	return 0
}
