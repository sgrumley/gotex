package components

import "github.com/rivo/tview"

type testFunction struct {
	name string
	// list of cases
}

func PopulateTests(list *tview.List) {
	list.AddItem("File A", "Details of test  A", 'a', nil).
		AddItem("File B", "Details of test  B", 'b', nil).
		AddItem("File C", "Details of test  C", 'c', nil)
}
