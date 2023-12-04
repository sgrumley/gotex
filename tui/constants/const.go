package constants

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	P          *tea.Program
	WindowSize tea.WindowSizeMsg
)

var DocStyle = lipgloss.NewStyle().Margin(0, 2)

var (
	CatpFlavor = catppuccin.Mocha
	HelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(CatpFlavor.Green().Hex)).Render
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(CatpFlavor.Red().Hex)).Render
	AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(CatpFlavor.Peach().Hex)).Render
)

type keymap struct {
	Run      key.Binding
	Enter    key.Binding
	Help     key.Binding
	Forward  key.Binding
	Backward key.Binding
	Quit     key.Binding
}

var Keymap = keymap{
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "help"),
	),
	Forward: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "expand / right"),
	),
	Backward: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "collapse / left"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}
