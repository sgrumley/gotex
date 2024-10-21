package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testFile struct {
	name string
	// list of functions
}

type testFiles struct {
	*tview.List
}

func newTestFiles(t *TUI) *testFiles {
	files := &testFiles{
		List: tview.NewList(),
	}

	files.SetTitle("Files")
	files.SetBorder(true)
	files.setKeybinding(t)
	files.Populate()

	return files
}

func (f *testFiles) setKeybinding(t *TUI) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			// TODO: run test
			f.AddItem("key enter registered", "Details of test case A", 'a', nil).
				AddItem("File B", "Details of test case B", 'b', nil).
				AddItem("File C", "Details of test case C", 'c', nil)
		case tcell.KeyCtrlR:
			// TODO: other events availible to files
		}

		// example using key instead of event
		switch event.Rune() {
		case 'd':
			f.AddItem("key press registered", "Details of test case A", 'a', nil)
			// case 'c':
		}

		return event
	})
}

func (f *testFiles) Populate() {
	f.AddItem("File A", "Details of test case A", 'a', nil).
		AddItem("File B", "Details of test case B", 'b', nil).
		AddItem("File C", "Details of test case C", 'c', nil)
}
