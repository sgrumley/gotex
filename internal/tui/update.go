package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) toggleBox() {
	m.activeBox = (m.activeBox + 1) % 2
	if m.activeBox == 0 {
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.testTree, cmd = m.testTree.Update(msg)
	cmds = append(cmds, cmd)

	m.secondPanel, cmd = m.secondPanel.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmds...)
}
