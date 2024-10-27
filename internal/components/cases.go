package components

import (
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
	cases.Populate()

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

func (c *testCases) Populate() {
	c.AddItem("File A", "Details of test case A", 'a', nil).
		AddItem("File B", "Details of test case B", 'b', nil).
		AddItem("File C", "Details of test case C", 'c', nil)
}
