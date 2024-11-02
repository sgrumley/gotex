package components

import (
	"sgrumley/gotex/pkg/config"
	"sgrumley/gotex/pkg/finder"

	"github.com/rivo/tview"
)

type resources struct {
	data *finder.Project
}

type state struct {
	// lastTest string // TODO: probably easiest to try and capture the cmd
	navigate  *navigate
	resources resources
	result    *results
}

func newState() *state {
	project := finder.InitProject()

	return &state{
		resources: resources{
			data: project,
		},
	}
}

type TUI struct {
	app   *tview.Application
	state *state
	theme Theme
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
	// TODO: there should be two configs -> theme and options
	// options should be found as part of finder
	// theme should be found here
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	SetAppStyling()
	t.theme = SetTheme("catppuccin mocha")

	// panels
	testTree := newTestTree(t, cfg)
	help := newHelpPane(t)

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
}
