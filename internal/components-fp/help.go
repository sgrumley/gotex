package components

import "github.com/rivo/tview"

func newHelpPane(t *TUI) *tview.TextView {
	help := tview.NewTextView()
	help.SetLabel("/: search, r: run test, R: rerun last test, A: run all tests, s: sync project, cntrl-u: scroll up, cntrl-d: scroll down, q: quit")
	SetTextViewStyling(t, help)
	return help
}
