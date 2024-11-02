package components

import (
	"github.com/rivo/tview"
)

type results struct {
	*tview.TextView
}

func newResultsPane(t *TUI) *results {
	res := &results{
		TextView: tview.NewTextView(),
	}

	res.SetBorder(true).SetTitle("Results")
	res.RenderResults("Run a test to see results...")
	SetTextViewStyling(t, res.TextView)

	return res
}

func (r *results) RenderResults(msg string) {
	r.TextView.Clear()
	r.TextView.SetDynamicColors(true).
		SetText(msg)
}
