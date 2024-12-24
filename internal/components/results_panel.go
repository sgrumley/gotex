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
	res.RenderResults("[green]Run[-] a test to see results...")
	res.SetDynamicColors(true)
	SetTextViewStyling(t, res.TextView)
	res.SetWrap(true)
	res.SetChangedFunc(func() {
		// this is required to allow a test is running screen before results
		t.app.Draw()
	})

	return res
}

func (r *results) RenderResults(msg string) {
	r.Clear()
	msg = tview.TranslateANSI(msg)
	r.SetDynamicColors(true)
	r.SetText(msg)
}
