package components

import "github.com/rivo/tview"

type panel interface {
	name() string
	focus(*TUI)
	unfocus()
	setKeybinding(*TUI)
}

type panels struct {
	currentPanel string
	panel        map[string]*tview.List
}

type resources struct {
	files []*testFile
	tests []*testFunction
	cases []*testCase
}

type state struct {
	lastTest  string // TODO: probably easiest to try and capture the cmd
	panels    panels
	navigate  *navigate
	resources resources // TODO: should this be the types from finder?
}

func newState() *state {
	initPanels := make(map[string]*tview.List)
	return &state{
		panels: panels{
			currentPanel: "",
			panel:        initPanels,
		},
	}
}

type TUI struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
}

func New() *TUI {
	return &TUI{
		app:   tview.NewApplication(),
		state: newState(),
	}
}

func (t *TUI) Start() error {
	t.initPanels()
	if err := t.app.Run(); err != nil {
		t.app.Stop()

		return err
	}

	return nil
}

func (t *TUI) Stop() {
	t.app.Stop()
}

func (t *TUI) initPanels() {
	// styling can be moved else where and add colors
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	// Create the main list (left panel)
	files := newTestFiles(t)
	tests := newTestFunctions(t)
	cases := newTestCases(t)

	// initialise panel state
	t.state.panels.panel["files"] = files.List
	t.state.panels.panel["tests"] = tests.List
	t.state.panels.panel["cases"] = cases.List
	t.state.panels.currentPanel = "files"

	// Create the results panel (right panel)
	results := newResultsPane()

	// this is the navigations column made up of interactive panels
	navFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(files, 0, 1, true).
		AddItem(tests, 0, 1, false).
		AddItem(cases, 0, 1, false)

	help := tview.NewTextView()
	help.SetLabel("/: search, q: quit, R: rerun last, r: run test, ?: more keys")

	// this is the whole screen
	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(navFlex, 0, 1, true).
		AddItem(results, 0, 6, false)

	// content with helper bar
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)

	t.app.SetRoot(layout, true)
}
