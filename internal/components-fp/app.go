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
	lastTest  finder.Node
	navigate  *navigate
	resources resources
	result    *results
	testTree  *TestTree
}

func newState(log *slog.Logger) *state {
	data, err := finder.InitProject(log)
	if err != nil {
		log.Error("failed to initialise project", slog.Any("error", err))
	}
	return &state{
		resources: resources{
			data: data,
		},
	}
}

type TUI struct {
	app   *tview.Application
	state *state
	theme Theme
	log   *slog.Logger
}

func New(log *slog.Logger) *TUI {
	return &TUI{
		app:   tview.NewApplication(),
		state: newState(log),
		log:   log,
	}
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
	testTree := newTestTree(t)
	t.app.SetFocus(testTree)
	t.state.testTree = testTree

	results := newResultsPane(t)
	t.state.result = results

	// layouts
	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(testTree.TreeView, 0, 1, true).
		AddItem(results.TextView, 0, 6, false)
	SetFlexStyling(t, contentLayout)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)
	SetFlexStyling(t, layout)

	t.app.SetRoot(layout, true)
	t.log.Info("app started successfully")
}
