package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testFunction struct {
	name string
	// list of cases
}

type testFunctions struct {
	*tview.List
}

func newTestFunctions(t *TUI) *testFunctions {
	funcs := &testFunctions{
		List: tview.NewList(),
	}

	funcs.SetTitle("Tests")
	funcs.SetBorder(true)
	funcs.setKeybinding(t)
	funcs.Populate()

	return funcs
}

func (f *testFunctions) setKeybinding(t *TUI) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			// TODO: run test
			f.AddItem("this should be added to resutls pane", "test completed", 'a', nil)
		}

		return event
	})
}

func (f *testFunctions) Populate() {
	f.AddItem("File A", "Details of test  A", 'a', nil).
		AddItem("File B", "Details of test  B", 'b', nil).
		AddItem("File C", "Details of test  C", 'c', nil)
}
