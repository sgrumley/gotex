package components

import (
	"github.com/rivo/tview"

	"sgrumley/gotex/pkg/config"
	"sgrumley/gotex/pkg/finder"
)

type panel interface {
	Populate(t *TUI, init bool, name string)
	GetList() *tview.List
	SetList(*tview.List)
}

type panels struct {
	currentPanel string
	// panel        map[string]*tview.List
	panel map[string]panel
}

type resources struct {
	currentFile *finder.File
	currentTest *finder.Function
	currentCase *finder.Case
	data        *finder.Project
}

type state struct {
	// lastTest string // TODO: probably easiest to try and capture the cmd
	panels    panels
	navigate  *navigate
	resources resources
	result    *results
}

func newState() *state {
	initPanels := make(map[string]panel)
	project := finder.InitProject()

	return &state{
		resources: resources{
			data: project,
		},
		panels: panels{
			currentPanel: "",
			panel:        initPanels,
		},
	}
}

type TUI struct {
	app   *tview.Application
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
	SetAppStyling()
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	// Create the main list (left panel)
	files := newTestFiles(t)
	tests := newTestFunctions(t)
	cases := newTestCases(t, cfg)

	// initialise panel state
	t.state.panels.panel["files"] = files
	t.state.panels.panel["tests"] = tests
	t.state.panels.panel["cases"] = cases
	t.state.panels.currentPanel = "files"

	// Create the results panel (right panel)
	results := newResultsPane(t)
	t.state.result = results

	// this is the navigations column made up of interactive panels
	navFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(files, 0, 1, true).
		AddItem(tests, 0, 1, false).
		AddItem(cases, 0, 1, false)
		// SetBackgroundColor(tcell.ColorPink)
	// SetFlexStyling(navFlex)

	help := tview.NewTextView()
	help.SetLabel("/: search, q: quit, R: rerun last, r: run test, ?: more keys")

	// this is the whole screen
	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(navFlex, 0, 1, true).
		AddItem(results.TextView, 0, 6, false)
		// SetBackgroundColor(tcell.ColorPink)
	// SetFlexStyling(contentLayout)

	// content with helper bar
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)
		// SetBackgroundColor(tcell.ColorPink)
	// SetFlexStyling(layout)

	t.app.SetRoot(layout, true)
}
