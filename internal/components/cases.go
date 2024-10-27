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

func (c *testCases) Populate(t *TUI, init bool) {
	// clear panel so dupes aren't added
	c.Clear()

	// get selected files from files panel
	var selectedFunction *finder.Function

	if !init {
		// selectedFunctionIndex := t.state.panels.panel["tests"].GetCurrentItem()
		// selectedFunctionName, _ := t.state.panels.panel["tests"].GetItemText(selectedFunctionIndex)

		// TODO: this set of data should be maps to avoid the loops -> make this change in api??
		// this has to be a map to avoid traversing all files -> all function

		// for _, function := range t.state.resources.data.Files {
		// 	if function.Name == selectedFunctionName {
		// 		selectedFunction = cases
		// 		break
		// 	}
		// }
	} else {
		selectedFunction = t.state.resources.currentTest
	}

	for _, cases := range selectedFunction.Cases {
		c.AddItem(cases.Name, "", 0, nil)
	}

	// update title with list count
	currentTitle := c.GetTitle()
	newTitle := fmt.Sprintf("%s (%d)", currentTitle, c.GetItemCount())
	c.SetTitle(newTitle)

	// set state
	t.state.resources.currentCase = &selectedFunction.Cases[0]
}
