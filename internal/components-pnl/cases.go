package components

import (
	"fmt"
	"path/filepath"
	"sgrumley/gotex/pkg/config"
	"sgrumley/gotex/pkg/finder"
	"sgrumley/gotex/pkg/runner"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testCase struct {
	// is this just the same as case from project?
}

type testCases struct {
	*tview.List
}

func newTestCases(t *TUI, cfg config.Config) *testCases {
	cases := &testCases{
		List: tview.NewList(),
	}

	SetListStyling(cases.List)
	cases.SetTitle("Cases")
	cases.SetBorder(true)
	cases.setKeybinding(t)
	cases.Populate(t, true, "")
	cases.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// not sure if there is a case for this
	})
	cases.SetSelectedFunc(func(index int, mainText, subText string, shortcut rune) {
		selectedFileIndex := t.state.panels.panel["files"].GetList().GetCurrentItem()
		selectedFileName, _ := t.state.panels.panel["files"].GetList().GetItemText(selectedFileIndex)

		// choice of function
		selectedFunctionIndex := t.state.panels.panel["tests"].GetList().GetCurrentItem()
		selectedFunctionName, _ := t.state.panels.panel["tests"].GetList().GetItemText(selectedFunctionIndex)

		path := filepath.Dir(t.state.resources.data.FileMap[selectedFileName].Path)
		currentFunction := t.state.resources.data.FileMap[selectedFileName].FunctionMap[selectedFunctionName]
		currentCase := currentFunction.CaseMap[mainText]

		// running test requires testname/casename
		caseName := fmt.Sprintf("%s/%s", currentFunction.Name, currentCase.Name)
		testResults, err := runner.RunTest(caseName, path, cfg)
		if err != nil {
			t.state.result.RenderResults(fmt.Sprintf("Error running test: %s", err.Error()))
			return
		}

		t.state.result.RenderResults(fmt.Sprintf("Test results:\n%s", testResults))
	})
	return cases
}

func (c *testCases) setKeybinding(t *TUI) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		// switch event.Key() {
		// case tcell.KeyEnter:
		// 	// TODO: run test
		// 	c.AddItem("this should be added to resutls pane", "test completed", 'a', nil)
		// }

		return event
	})
}

func (c *testCases) Populate(t *TUI, init bool, functionName string) {
	// clear panel so dupes aren't added
	c.Clear()

	// get selected files from files panel
	var selectedFunction *finder.Function

	if !init {
		// choice of file
		selectedFileIndex := t.state.panels.panel["files"].GetList().GetCurrentItem()
		selectedFileName, _ := t.state.panels.panel["files"].GetList().GetItemText(selectedFileIndex)

		// choice of function
		selectedFunctionIndex := t.state.panels.panel["tests"].GetList().GetCurrentItem()
		selectedFunctionName, _ := t.state.panels.panel["tests"].GetList().GetItemText(selectedFunctionIndex)

		selectedFunction = t.state.resources.data.FileMap[selectedFileName].FunctionMap[selectedFunctionName]

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
	if strings.Contains(currentTitle, "(") {
		titleSplit := strings.Split(currentTitle, "(")
		currentTitle = titleSplit[0]
	}

	newTitle := fmt.Sprintf("%s(%d)", currentTitle, c.GetItemCount())
	c.SetTitle(newTitle)
}

func (c *testCases) GetList() *tview.List {
	return c.List
}

func (c *testCases) SetList(l *tview.List) {
	c.List = l
}
