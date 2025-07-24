# gotex/pkg - Go Test Analysis Library

[![Go Reference](https://pkg.go.dev/badge/github.com/sgrumley/gotex/pkg.svg)](https://pkg.go.dev/github.com/sgrumley/gotex/pkg)

This is the core library component of gotex, providing Go test analysis and execution capabilities as an importable package.

## Overview

The gotex/pkg library provides:
- **Project scanning**: Analyze Go projects to discover test structure
- **Test execution**: Run tests at various granularity levels (project, package, function, case)
- **Configuration management**: Flexible configuration with environment variable and file support
- **Structured logging**: Context-aware logging throughout the analysis process
- **Data models**: Rich data structures representing Go test hierarchies

## Installation

```bash
go get github.com/sgrumley/gotex/pkg
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/sgrumley/gotex/pkg/config"
    "github.com/sgrumley/gotex/pkg/scanner"
    "github.com/sgrumley/gotex/pkg/slogger"
)

func main() {
    ctx := context.Background()
    
    // Setup logging
    logger, err := slogger.New()
    if err != nil {
        log.Fatal(err)
    }
    ctx = slogger.AddToContext(ctx, logger)
    
    // Load configuration
    cfg, err := config.GetConfig(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Find project root
    root, err := scanner.FindGoProjectRoot()
    if err != nil {
        log.Fatal(err)
    }
    
    // Scan project for test structure
    project, err := scanner.Scan(ctx, cfg, root)
    if err != nil {
        log.Fatal(err)
    }
    
    // Navigate the discovered structure
    fmt.Printf("Project: %s\n", project.GetName())
    for _, pkg := range project.Packages {
        fmt.Printf("  Package: %s\n", pkg.GetName())
        for _, file := range pkg.Files {
            fmt.Printf("    File: %s\n", file.GetName())
            for _, fn := range file.Functions {
                fmt.Printf("      Function: %s (%d cases)\n", 
                    fn.GetName(), len(fn.Cases))
            }
        }
    }
}
```

## Packages

### config
Configuration management with support for YAML files and environment variables.

### models  
Core data structures representing the project hierarchy: Project → Package → File → Function → Case.

### scanner
Go source code analysis for extracting test structure from AST.

### runner
Test execution engine supporting various granularity levels.

### slogger
Structured logging with context support.

### path
File system path utilities including home directory expansion.

### ansi
Terminal color and formatting utilities.

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/sgrumley/gotex/pkg).

## License

See the LICENSE file in the repository root.