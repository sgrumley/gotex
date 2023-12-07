package tui

import (
	"sgrumley/test-tui/internal/tui/constants"
	secondpanelexample "sgrumley/test-tui/pkg/secondPanelExample"
	"sgrumley/test-tui/pkg/testtree"
)

type sessionState int

const (
	idleState sessionState = iota
	runState
	resultState
)

type model struct {
	keys        constants.Keymap
	state       sessionState
	testTree    testtree.Model
	secondPanel secondpanelexample.Model
	activeBox   int
}

func initModel() model {
	// cfg, err := config.ParseConfig()
	// theme := theme.GetTheme(cfg.Theme.AppTheme)

	return model{
		keys:        constants.Keymaps,
		secondPanel: secondpanelexample.InitialModel(),
		// this data needs to be retrieved from a function
		testTree: testtree.Model{
			Menu: []testtree.MenuItem{
				{
					Title: "Item 1",
					Children: []testtree.MenuItem{
						{Title: "Subitem 1.1", IsSubmenu: true, ParentIndex: 0},
						{Title: "Subitem 1.2", IsSubmenu: true, ParentIndex: 0},
					},
					ParentIndex: -1,
				},
				{Title: "Item 2", ParentIndex: -1},
				{Title: "Item 3", ParentIndex: -1},
			},
			CurrentIndex: 0,
		},
		activeBox: 0,
		state:     0,
	}
}
