package components

import "github.com/rivo/tview"

type testCase struct {
	// is this just the same as case from project?
}

func PopulateCases(list *tview.List) {
	list.AddItem("File A", "Details of test case A", 'a', nil).
		AddItem("File B", "Details of test case B", 'b', nil).
		AddItem("File C", "Details of test case C", 'c', nil)
}
