package main

import (
	"context"
	"log"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/models"
	"github.com/sgrumley/gotex/pkg/scanner"
)

func main() {
	// add to ctx
	// logr, err := logger.New(
	// 	logger.WithLevel(slog.LevelDebug),
	// 	logger.WithSource(false),
	// )
	// if err != nil {
	// 	log.Fatal("failed setting logger")
	// 	fmt.Println("error: ", err.Error())
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

	p, err := scanner.Scan(ctx, cfg, root)
	if err != nil {
		log.Fatal(err)
	}

	err = models.GenerateTree(p)
	if err != nil {
		log.Fatal(err)
	}

	p.Tree.Print()
}
