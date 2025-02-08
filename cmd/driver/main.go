package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sgrumley/gotex/pkgv2/config"
	"github.com/sgrumley/gotex/pkgv2/models"
	"github.com/sgrumley/gotex/pkgv2/scanner"
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

	fmt.Println("dir: ", p.RootDir)
	p.Tree.Print()
}
