package components

import "github.com/rivo/tview"

// TODO: extend this to have a modal as well
func newHelpPane(t *TUI) *tview.TextView {
	help := tview.NewTextView()
	help.SetTextColor(t.theme.SecondaryTextColor)
	help.SetLabel("/: search, r: run test, R: rerun last test, A: run all tests, s: sync project, up: k, down:j, h: collapse, l: expand, C: console, c: config, ctrl-u: scroll up, ctrl-d: scroll down, q: quit")

	return help
}
