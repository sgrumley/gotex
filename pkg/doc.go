// Package gotex provides Go test analysis and execution capabilities.
//
// This library enables programmatic analysis of Go projects to discover
// test structure, execute tests at various granularity levels, and process
// test results. It's designed to be used by applications that need to
// interact with Go test suites in sophisticated ways.
//
// # Core Components
//
// The library is organized into several focused packages:
//
//   - config: Configuration management for test execution parameters
//   - models: Data structures representing project hierarchy and test structure
//   - scanner: Go source code analysis and test discovery
//   - runner: Test execution engine with configurable output handling
//   - slogger: Structured logging with context support
//   - path: File system path utilities
//   - ansi: Terminal formatting and color utilities
//
// # Basic Usage
//
// The typical workflow involves scanning a project to discover its test
// structure, then using that structure to execute tests:
//
//	// Setup
//	ctx := context.Background()
//	logger, _ := slogger.New()
//	ctx = slogger.AddToContext(ctx, logger)
//	cfg, _ := config.GetConfig(ctx)
//	
//	// Scan project
//	root, _ := scanner.FindGoProjectRoot()
//	project, _ := scanner.Scan(ctx, cfg, root)
//	
//	// Execute tests
//	response, _ := project.RunTest(ctx)
//
// # Project Structure
//
// The library models Go projects as a hierarchy:
//
//	Project (root directory with go.mod)
//	├── Package (Go package directory)
//	│   ├── File (Go test file)
//	│   │   ├── Function (test function with t.Run calls)
//	│   │   │   └── Case (individual test case within table-driven test)
//
// Each level implements the models.Node interface, enabling uniform
// tree traversal and test execution operations.
//
// # Configuration
//
// Test execution behavior is controlled through the config.Config struct,
// which supports loading from:
//   - Environment variables (GOTEX_CONFIG_FILE_PATH)
//   - YAML configuration files (~/.config/gotex/config.yaml)
//   - Embedded default configuration
//
// # Integration
//
// This library is used by the gotex interactive terminal application
// but is designed to be usable in any Go program that needs to analyze
// or execute tests programmatically.
package gotex