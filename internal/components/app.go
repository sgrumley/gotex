package components

import (
	"context"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/models"
	"github.com/sgrumley/gotex/pkg/scanner"
	"github.com/sgrumley/gotex/pkg/slogger"

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
	config   *ConfigModal
	lastKey  rune
}

type Data struct {
	project   *models.Project
	flattened *models.FlatProject // this should become useful once I update the search names to append the parent node

	lastTest models.Node
}

type consoleData struct {
	currentMessage string
	active         bool
	panel          *console
	flex           *tview.Flex
}

func newState(ctx context.Context, cfg config.Config, root string) (*state, error) {
	data, err := scanner.Scan(ctx, cfg, root)
	if err != nil {
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
}

func New(ctx context.Context, cfg config.Config, root string) (*TUI, error) {
	data, err := newState(ctx, cfg, root)
	if err != nil {
		return nil, err
	}
	return &TUI{
		app:   tview.NewApplication(),
		state: data,
	}, nil
}

func (t *TUI) Start(ctx context.Context) error {
	t.initPanels(ctx)
	if err := t.app.Run(); err != nil {
		t.app.Stop()

		return err
	}

	return nil
}

func (t *TUI) Stop() {
	t.app.Stop()
}

func (t *TUI) initPanels(ctx context.Context) {
	// TODO: there should be two configs -> theme and options
	// options should be found as part of finder
	// theme should be found here

	log, err := slogger.FromContext(ctx)
	if err != nil {
		log, _ := slogger.New()
		ctx = slogger.AddToContext(ctx, log)
	}

	SetAppStyling()
	t.theme = SetTheme("catppuccin mocha")
	tview.Styles = t.theme.Theme

	// pages
	pages := tview.NewPages()

	// panels
	help := newHelpPane(t)

	console := newConsolePane(t)
	t.state.ui.console.panel = console

	search := newSearchModal(ctx, t)
	t.state.ui.search = search

	config := newConfigModal(t)
	t.state.ui.config = config

	testTree := newTestTree(ctx, t)
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

	pages.AddPage(configPage, config.modal, true, false)
	pages.AddPage(searchPage, search.modal, true, false)

	t.app.SetRoot(pages, true)
	log.Info("app started successfully")
}
