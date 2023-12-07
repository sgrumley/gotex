package main

import (
	"log"

	"sgrumley/test-tui/internal/tui"
)

func main() {
	// menu, projectName := data.LoadDummyData()
	// fmt.Println(projectName)
	// data.PrintMenu(menu, 0)
	err := tui.StartTea()
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
