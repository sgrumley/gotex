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
	Menu                []menuItem
	CurrentIndex        int  // Index of the currently selected item in the main menu
	CurrentSubmenuIndex int  // Index of the currently selected item in the submenu
	InSubmenu           bool // Flag to indicate if the user is currently in a submenu
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
		// Render main menu item
		prefix := " "
		if i == m.CurrentIndex && !m.InSubmenu {
			prefix = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", prefix, item.Title))

		// Render submenu items if expanded
		if item.Expanded {
			for j, subitem := range item.Children {
				subprefix := "  "
				if j == m.CurrentSubmenuIndex && m.InSubmenu && i == m.CurrentIndex {
					subprefix = "> "
				}
				s.WriteString(fmt.Sprintf("%s%s\n", subprefix, subitem.Title))
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
			if m.InSubmenu {
				// Navigate within the submenu
				if m.CurrentSubmenuIndex < len(m.Menu[m.CurrentIndex].Children)-1 {
					m.CurrentSubmenuIndex++
				}
			} else if m.Menu[m.CurrentIndex].Expanded {
				// Enter the submenu
				m.InSubmenu = true
				m.CurrentSubmenuIndex = 0
			} else if m.CurrentIndex < len(m.Menu)-1 {
				// Navigate in the main menu
				m.CurrentIndex++
			}
		case "k", "up":
			if m.InSubmenu {
				// Navigate within the submenu
				if m.CurrentSubmenuIndex > 0 {
					m.CurrentSubmenuIndex--
				}
			} else if m.CurrentIndex > 0 {
				// Navigate in the main menu
				m.CurrentIndex--
			}
		case "enter":
			// Toggle expansion and navigation of submenu
			if m.InSubmenu {
				// Handle submenu item selection
				// (You can add actions here for when a submenu item is selected)
				m.InSubmenu = false
			} else if len(m.Menu[m.CurrentIndex].Children) > 0 {
				m.Menu[m.CurrentIndex].Expanded = !m.Menu[m.CurrentIndex].Expanded
			}
		case "esc":
			// Exit submenu
			if m.InSubmenu {
				m.InSubmenu = false
				m.CurrentSubmenuIndex = 0
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
