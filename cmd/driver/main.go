package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/models"
	"github.com/sgrumley/gotex/pkg/scanner"
	"github.com/sgrumley/gotex/pkg/slogger"
)

// Driver serves as an example of how to use the scanning api
func main() {
	ctx := context.Background()
	log, err := slogger.New(
		slogger.WithLevel(slog.LevelDebug),
		slogger.WithSource(false),
	)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	ctx = slogger.AddToContext(ctx, log)

	root, err := scanner.FindGoProjectRoot()
	if err != nil {
		log.Fatal("failed to find root of project", err)
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {

		fmt.Printf("failed to load a config: %s\n", err.Error())
		log.Fatal("failed to load a config: %w", err)
	}

	p, err := scanner.Scan(ctx, cfg, root)
	if err != nil {
		log.Fatal("failed scanning project", err)
	}

	err = models.GenerateTree(p)
	if err != nil {
		log.Fatal("failed to generate tree", err)
	}

	p.Tree.Print()
}
