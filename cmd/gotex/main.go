package main

import (
	"fmt"
	"log/slog"
	"os"
	"sgrumley/gotex/internal/components"
	logger "sgrumley/gotex/pkg/logging"
)

func main() {
	os.Exit(run())
}

func run() int {
	log, err := logger.New(
		logger.WithLevel(slog.LevelDebug),
		logger.WithSource(true),
	)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return 1
	}

	app, err := components.New(log)
	if err != nil {
		fmt.Printf("failed to initialise project: %s", err.Error())
		return 1
	}
	err = app.Start()
	if err != nil {
		// to file
		log.Error("application crashed", slog.Any("error", err))
		// to stdout
		fmt.Printf("application crashed: %s", err.Error())
		return 1
	}

	return 0
}
