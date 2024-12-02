package components

import "github.com/rivo/tview"

// TODO: extend this to have a modal as well
func newHelpPane(t *TUI) *tview.TextView {
	help := tview.NewTextView()
	help.SetLabel("/: search, r: run test, R: rerun last test, A: run all tests, s: sync project, up: k, down:j, h: collapse, l: expand, c: console, C: config, cntrl-u: scroll up, cntrl-d: scroll down, q: quit")
	SetTextViewStyling(t, help)
	return help
}
