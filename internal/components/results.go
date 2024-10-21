package components

import "github.com/rivo/tview"

func RenderResults(view *tview.TextView) {
	view.SetDynamicColors(true).
		SetText("See test results here")
}
