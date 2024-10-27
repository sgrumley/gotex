package components

import (
	"fmt"
	"sgrumley/gotex/pkg/finder"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testFile struct {
	name  string
	files []finder.File
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
	files.Populate(t)

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

func (f *testFiles) hoverEvent() {
	// f.SetSelectedFunc will give it a function to execute when 'enter' is pressed on an element this should be the same as 'r'
}

func (f *testFiles) Populate(t *TUI) {
	f.Clear()
	for _, file := range t.state.resources.data.Files {
		f.AddItem(file.Name, "", 0, nil)
	}

	currentTitle := f.GetTitle()
	newTitle := fmt.Sprintf("%s (%d)", currentTitle, f.GetItemCount())
	f.SetTitle(newTitle)

	// HACK: an initial value is required to choose which test->case is displayed in other panels
	// this may not sync correctly with no garunteed order to iterating a map

	// set state
	for _, file := range t.state.resources.data.Files {
		t.state.resources.currentFile = file
		break
	}
}
