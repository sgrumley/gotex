package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create the main list (left panel)
	packages := tview.NewList().
		AddItem("File A", "Details of package A", 'a', nil).
		AddItem("File B", "Details of package B", 'b', nil).
		AddItem("File C", "Details of package C", 'c', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	packages.SetBorder(true).SetTitle("Packages")

	files := tview.NewList().
		AddItem("File A", "Details of test file A", 'a', nil).
		AddItem("File B", "Details of test file B", 'b', nil).
		AddItem("File C", "Details of test file C", 'c', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	files.SetBorder(true).SetTitle("Files")

tests := tview.NewList().
		AddItem("File A", "Details of test  A", 'a', nil).
		AddItem("File B", "Details of test  B", 'b', nil).
		AddItem("File C", "Details of test  C", 'c', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	tests.SetBorder(true).SetTitle("Tests")

cases := tview.NewList().
		AddItem("File A", "Details of test case A", 'a', nil).
		AddItem("File B", "Details of test case B", 'b', nil).
		AddItem("File C", "Details of test case C", 'c', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	cases.SetBorder(true).SetTitle("Test Cases")



	// Create the detail panel (right panel)
	detail := tview.NewTextView().
		SetDynamicColors(true).
		SetText("See test results here")
	detail.SetBorder(true).SetTitle("Results")

	// Create the log/output panel (bottom panel)
	// log := tview.NewTextView().
	// 	SetDynamicColors(true).
	// 	SetText("This is where logs or output would appear.")
	// log.SetBorder(true).SetTitle("Log")

	// Handle list selection to update details
	tests.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		detail.SetText("Details of " + mainText)
		// log.SetText("Showing log for " + mainText)
	})

	// Create a flex layout for the main window
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(packages, 0, 2, true).
		AddItem(files, 0, 3, true).
		AddItem(tests, 0, 4, false).
		AddItem(cases, 0, 5, false)

	// Create another flex layout to add the log panel at the bottom
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(mainFlex, 0, 1, true).
		AddItem(detail, 0, 6, false)
		// AddItem(log, 7, 1, false)

	// Set the root and run the application
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
