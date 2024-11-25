package components

import (
	"log/slog"
	"sgrumley/gotex/pkg/finder"

	"github.com/rivo/tview"
)

type resources struct {
	data *finder.Project
}

type state struct {
	lastTest finder.Node
	// navigate  *navigate
	resources resources
	result    *results
	testTree  *TestTree
	console   consoleData
}

type consoleData struct {
	currentMessage string
	active         bool
	panel          *console
	flex           *tview.Flex
}

func newState(log *slog.Logger) (*state, error) {
	data, err := finder.InitProject(log)
	if err != nil {
		log.Error("failed to initialise project", slog.Any("error", err))
		return nil, err
	}

	return &state{
		resources: resources{
			data: data,
		},
		console: consoleData{
			active: false, // TODO: off by default??
		},
	}, nil
}

type TUI struct {
	app   *tview.Application
	state *state
	theme Theme
	log   *slog.Logger
}

func New(log *slog.Logger) (*TUI, error) {
	data, err := newState(log)
	if err != nil {
		return nil, err
	}
	return &TUI{
		app:   tview.NewApplication(),
		state: data,
		log:   log,
	}, nil
}

func (t *TUI) Start() error {
	t.initPanels()
	if err := t.app.Run(); err != nil {
		t.log.Error("app stopping", slog.Any("error", err))
		t.app.Stop()

		return err
	}

	return nil
}

func (t *TUI) Stop() {
	t.app.Stop()
}

func (t *TUI) initPanels() {
	// TODO: there should be two configs -> theme and options
	// options should be found as part of finder
	// theme should be found here

	SetAppStyling()
	t.theme = SetTheme("catppuccin mocha")

	// panels
	help := newHelpPane(t)

	// TODO: this should have a struct of all printed fields in the state and this is kept up to date and the console should render on open
	console := newConsolePane(t)
	t.state.console.panel = console

	testTree := newTestTree(t)
	t.app.SetFocus(testTree)
	t.state.testTree = testTree

	results := newResultsPane(t)
	t.state.result = results

	// layouts
	outputLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(results.TextView, 0, 8, false)
	// NOTE: if on by default needs a bit of work
	// AddItem(console.TextView, 8, 1, false)
	t.state.console.flex = outputLayout

	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(testTree.TreeView, 45, 1, true).
		AddItem(outputLayout, 0, 6, false)
	SetFlexStyling(t, contentLayout)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)
	SetFlexStyling(t, layout)

	t.app.SetRoot(layout, true)
	t.log.Info("app started successfully")
}
