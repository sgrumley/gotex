package testtree

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
