package components

import (
	"log/slog"
	"sgrumley/gotex/pkg/finder"

	"github.com/rivo/tview"
)

var (
	homePage   = "home"
	configPage = "config"
	searchPage = "search"
)

type state struct {
	ui   UI
	data Data
}

type UI struct {
	result   *results
	testTree *TestTree
	console  *consoleData
	pages    *tview.Pages
	search   *searchModal
}

type Data struct {
	project   *finder.Project
	flattened *finder.FlatProject // this should become useful once I update the search names to append the parent node

	lastTest finder.Node
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
		data: Data{
			project:   data,
			flattened: data.FlattenAllNodes(),
		},
		ui: UI{
			console: &consoleData{},
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
	tview.Styles = t.theme.Theme

	// pages
	pages := tview.NewPages()

	// panels
	help := newHelpPane(t)

	console := newConsolePane(t)
	t.state.ui.console.panel = console

	search := newSearchModal(t)
	t.state.ui.search = search

	testTree := newTestTree(t)
	t.app.SetFocus(testTree)
	t.state.ui.testTree = testTree

	results := newResultsPane(t)
	t.state.ui.result = results

	// layouts
	outputLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(results.TextView, 0, 8, false)

	t.state.ui.console.flex = outputLayout

	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(testTree.TreeView, 45, 1, true).
		AddItem(outputLayout, 0, 10, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)

	pages.AddPage(homePage, layout, true, true)

	// initialising pages state here so that newConfigModal has access
	t.state.ui.pages = pages

	configModal := newConfigModal(t)
	pages.AddPage(configPage, configModal, true, false)
	pages.AddPage(searchPage, search.modal, true, false)

	t.app.SetRoot(pages, true)
	t.log.Info("app started successfully")
}
