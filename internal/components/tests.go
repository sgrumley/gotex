package components

import (
	"fmt"
	"sgrumley/gotex/pkg/finder"
	"strings"

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

	SetListStyling(funcs.List)
	funcs.SetTitle("Tests")
	funcs.SetBorder(true)
	funcs.setKeybinding(t)
	funcs.Populate(t, true, "")
	funcs.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		t.state.panels.panel["cases"].Populate(t, false, "")
		// t.state.panels.panel["cases"].Populate(t, false, mainText)
	})
	funcs.SetSelectedFunc(func(index int, mainText, subText string, shortcut rune) {
		// call go test parsing function command (no appending /case)
	})

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

func (f *testFunctions) Populate(t *TUI, init bool, fileName string) {
	// clear panel so dupes aren't added
	f.Clear()

	// get selected files from files panel
	var selectedFile *finder.File

	if !init {
		selectedFileIndex := t.state.panels.panel["files"].GetList().GetCurrentItem()
		selectedFileName, _ := t.state.panels.panel["files"].GetList().GetItemText(selectedFileIndex)
		selectedFile = t.state.resources.data.Files[selectedFileName]

	} else {
		selectedFile = t.state.resources.currentFile
	}

	for _, test := range selectedFile.Functions {
		f.AddItem(test.Name, "", 0, nil)
	}

	// update title with list count
	currentTitle := f.GetTitle()
	if strings.Contains(currentTitle, "(") {
		titleSplit := strings.Split(currentTitle, "(")
		currentTitle = titleSplit[0]
	}

	newTitle := fmt.Sprintf("%s(%d)", currentTitle, f.GetItemCount())
	f.SetTitle(newTitle)

	// HACK: an initial value is required to choose which test->case is displayed in other panels
	// this may not sync correctly with no garunteed order to iterating a map

	// set state
	for _, function := range selectedFile.Functions {
		t.state.resources.currentTest = function
		break
	}
}

func (f *testFunctions) GetList() *tview.List {
	return f.List
}

func (f *testFunctions) SetList(l *tview.List) {
	f.List = l
}
