package components

import "github.com/rivo/tview"

func newHelpPane(t *TUI) *tview.TextView {
	help := tview.NewTextView()
	help.SetLabel("/: search, q: quit, R: rerun last, r: run test, ?: more keys")
	SetTextViewStyling(t, help)
	return help
}
