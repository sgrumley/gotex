package main

import (
	"fmt"
	"log/slog"
	"os"
	"sgrumley/gotex/internal/components-fp"
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

	app := components.New(log)
	err = app.Start()
	if err != nil {
		log.Error("application crashed", slog.Any("error", err))
		fmt.Printf("application crashed: %s", err.Error())
		return 1
	}

	return 0
}
