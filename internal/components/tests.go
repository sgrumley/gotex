package components

import (
	"fmt"
	"sgrumley/gotex/pkg/finder"

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
	funcs.Populate(t, true)

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

// TODO: wipe the panel before adding things
func (f *testFunctions) Populate(t *TUI, init bool) {
	// clear panel so dupes aren't added
	f.Clear()

	// get selected files from files panel
	var selectedFileIndex int
	var selectedFileName string
	var selectedFile *finder.File

	if !init {
		selectedFileIndex = t.state.panels.panel["files"].GetCurrentItem()
		selectedFileName, _ = t.state.panels.panel["files"].GetItemText(selectedFileIndex)
		// TODO: this set of data should be maps to avoid the loops -> make this change in api??
		for _, file := range t.state.resources.data.Files {
			if file.Name == selectedFileName {
				selectedFile = file
				break
			}
		}
	} else {
		selectedFile = t.state.resources.currentFile
	}

	for _, test := range selectedFile.Functions {
		f.AddItem(test.Name, "", 0, nil)
	}

	// update title with list count
	currentTitle := f.GetTitle()
	newTitle := fmt.Sprintf("%s (%d)", currentTitle, f.GetItemCount())
	f.SetTitle(newTitle)

	// set state
	t.state.resources.currentTest = selectedFile.Functions[0]
}
