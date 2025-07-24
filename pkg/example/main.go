package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/scanner"
	"github.com/sgrumley/gotex/pkg/slogger"
)

// Example demonstrates basic usage of the gotex/pkg library.
func main() {
	ctx := context.Background()

	// Setup structured logging
	logger, err := slogger.New()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	ctx = slogger.AddToContext(ctx, logger)

	// Load configuration (from env vars, config file, or defaults)
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Printf("Configuration loaded: verbose=%t, coverage=%t\n", cfg.Verbose, cfg.Cover)

	// Find the Go project root (looks for go.mod)
	root, err := scanner.FindGoProjectRoot()
	if err != nil {
		log.Fatal("Failed to find Go project root:", err)
	}

	fmt.Printf("Found project root: %s\n", root)

	// Scan the project to discover test structure
	project, err := scanner.Scan(ctx, cfg, root)
	if err != nil {
		log.Fatal("Failed to scan project:", err)
	}

	// Display the discovered structure
	fmt.Printf("\nProject: %s\n", project.GetName())
	fmt.Printf("Total packages: %d\n", len(project.Packages))

	for _, pkg := range project.Packages {
		fmt.Printf("\n  📦 Package: %s\n", pkg.GetName())
		for _, file := range pkg.Files {
			fmt.Printf("    📄 File: %s\n", file.GetName())
			for _, fn := range file.Functions {
				fmt.Printf("      🧪 Function: %s (%d test cases)\n", fn.GetName(), len(fn.Cases))
				for _, testCase := range fn.Cases {
					fmt.Printf("        ✅ Case: %s\n", testCase.GetName())
				}
			}
		}
	}

	fmt.Printf("\nProject analysis complete! Found %d packages with test structure.\n", len(project.Packages))
}