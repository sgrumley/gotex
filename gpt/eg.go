package gpt

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	Title       string
	Children    []menuItem
	Expanded    bool
	IsSubmenu   bool
	ParentIndex int // Index of the parent item, -1 for root items
}

type model struct {
	Menu         []menuItem
	CurrentIndex int // Index of the currently selected item
}

func Driver() {
	initialModel := model{
		Menu: []menuItem{
			{
				Title: "Item 1",
				Children: []menuItem{
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
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("Menu:\n\n")

	for i, item := range m.Menu {
		// Render menu item
		prefix := " "
		if i == m.CurrentIndex {
			prefix = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", prefix, item.Title))

		// Render children if expanded
		if item.Expanded {
			for _, child := range item.Children {
				s.WriteString(fmt.Sprintf("  - %s\n", child.Title))
			}
		}
	}

	return s.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			// Handle down key
			if m.CurrentIndex < len(m.Menu)-1 {
				m.CurrentIndex++
			}
		case "k", "up":
			// Handle up key
			if m.CurrentIndex > 0 {
				m.CurrentIndex--
			}
		case "enter", "o":
			// Toggle expansion of submenu
			if len(m.Menu[m.CurrentIndex].Children) > 0 {
				m.Menu[m.CurrentIndex].Expanded = !m.Menu[m.CurrentIndex].Expanded
			}
		case "q":
			os.Exit(1)
		}
	}
	return m, nil
}

func (m model) Init() tea.Cmd {
	// Initialize with the root menu
	return nil
}
