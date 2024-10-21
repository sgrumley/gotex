package main

import (
	"fmt"
	"os"
	"sgrumley/gotex/internal/components"
)

func main() {
	os.Exit(run())
}

func run() int {
	app := components.New()
	err := app.Start()
	if err != nil {
		fmt.Printf("application crashed: %s", err.Error())
		return 1
	}

	return 0
}
