package components

import (
	"fmt"
	"sgrumley/gotex/pkg/finder"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testCase struct {
	// is this just the same as case from project?
}

type testCases struct {
	*tview.List
}

func newTestCases(t *TUI) *testCases {
	cases := &testCases{
		List: tview.NewList(),
	}

	cases.SetTitle("Cases")
	cases.SetBorder(true)
	cases.setKeybinding(t)
	cases.Populate(t, true)
	// cases.SetChangedFunc(ChangeCase)

	return cases
}

func (c *testCases) setKeybinding(t *TUI) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			// TODO: run test
			c.AddItem("this should be added to resutls pane", "test completed", 'a', nil)
		}

		return event
	})
}

func ChangeCase(index int, mainText string, secondaryText string, shortcut rune) {
	// not sure if there is a case for this
}

func (c *testCases) Populate(t *TUI, init bool) {
	// clear panel so dupes aren't added
	c.Clear()

	// get selected files from files panel
	var selectedFunction *finder.Function

	if !init {
		// choice of file
		selectedFileIndex := t.state.panels.panel["files"].GetCurrentItem()
		selectedFileName, _ := t.state.panels.panel["files"].GetItemText(selectedFileIndex)

		// choice of function
		selectedFunctionIndex := t.state.panels.panel["tests"].GetCurrentItem()
		selectedFunctionName, _ := t.state.panels.panel["tests"].GetItemText(selectedFunctionIndex)

		selectedFunction = t.state.resources.data.Files[selectedFileName].Functions[selectedFunctionName]

	} else {
		selectedFunction = t.state.resources.currentTest
	}

	for _, cases := range selectedFunction.Cases {
		c.AddItem(cases.Name, "", 0, nil)
	}

	// set state
	for _, cases := range selectedFunction.Cases {
		t.state.resources.currentCase = cases
		break
	}

	// update title with list count
	currentTitle := c.GetTitle()
	newTitle := fmt.Sprintf("%s (%d)", currentTitle, c.GetItemCount())
	c.SetTitle(newTitle)
}
