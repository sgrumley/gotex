package testtree

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
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
