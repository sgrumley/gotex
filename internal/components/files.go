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

func (f *testFiles) hoverEvent() {}

func (f *testFiles) Populate(t *TUI) {
	f.Clear()
	for _, file := range t.state.resources.data.Files {
		f.AddItem(file.Name, "", 0, nil)
	}

	currentTitle := f.GetTitle()
	newTitle := fmt.Sprintf("%s (%d)", currentTitle, f.GetItemCount())
	f.SetTitle(newTitle)
	// HACK: assuming that every time this function is called it will reset the selected item to index 0
	t.state.resources.currentFile = t.state.resources.data.Files[0]
}
