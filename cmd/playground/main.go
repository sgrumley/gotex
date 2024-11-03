package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a new TextView
	textView := tview.NewTextView().
		SetDynamicColors(true). // Enable dynamic colors to interpret ANSI codes
		SetText("[green]PASS[-] Test Passed!")
		// SetText("[92mPASS[0m Test Passed!")

	// Set up layout and add the TextView
	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
