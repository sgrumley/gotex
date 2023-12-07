package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	leftBox := m.testTree.View()
	rightBox := m.secondPanel.View()

	// switch m.state {
	// case idleState:
	// 	rightBox = m.help.View()

	// case runState:
	// 	rightBox = m.LiveTestSummary.View()

	// case resultState:
	// 	rightBox = m.ResultTestSummary.View()
	// }

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}
