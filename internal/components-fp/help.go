package components

import "github.com/rivo/tview"

func newHelpPane(t *TUI) *tview.TextView {
	help := tview.NewTextView()
	help.SetLabel("/: search, r: run test, R: rerun last test, s: sync project, q: quit")
	SetTextViewStyling(t, help)
	return help
}
