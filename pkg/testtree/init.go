package testtree

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Driver() {
	_ = Model{
		Menu: []MenuItem{
			{
				Title: "Item 1",
				Children: []MenuItem{
					{Title: "Subitem 1.1", IsSubmenu: true, ParentIndex: 0},
					{Title: "Subitem 1.2", IsSubmenu: true, ParentIndex: 0},
				},
				ParentIndex: -1,
			},
			{Title: "Item 2", ParentIndex: -1},
			{Title: "Item 3", ParentIndex: -1},
		},
		CurrentIndex: 0,
	}
	// p := tea.NewProgram(initialModel)
	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
}

func (m Model) Init() tea.Cmd {
	// Initialize with the root menu
	return nil
}
