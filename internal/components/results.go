package components

import "github.com/rivo/tview"

func newResultsPane() *tview.TextView {
	results := tview.NewTextView()
	results.SetBorder(true).SetTitle("Results")
	RenderResults(results)
	SetTextViewStyling(results)
	return results
}

func RenderResults(view *tview.TextView) {
	view.SetDynamicColors(true).
		SetText("See test results here")
}
