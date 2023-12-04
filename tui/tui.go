package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"sgrumley/test-tui/tui/constants"
)

func StartTea() error {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("couldn't open a file for logging: ", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	m, _ := InitProject()
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		fmt.Println("error running program: ", err)
		os.Exit(1)
	}
	return nil
}
