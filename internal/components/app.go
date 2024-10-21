package components

import "github.com/rivo/tview"

type panel interface {
	name() string
	focus(*TUI)
	unfocus()
	setKeybinding(*TUI)
}

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	packages []*testPackage
	files    []*testFile
	tests    []*testFunction
	cases    []*testCase
}

type state struct {
	lastTest  string // TODO: probably easiest to try and capture the cmd
	panels    panels
	navigate  *navigate
	resources resources
	stopChans map[string]chan int
}

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
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

// TODO:  methods on tui for each panel

func (t *TUI) initPanels() {
	// styling can be moved else where and add colors
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'

	// Create the main list (left panel)
	// packages := tview.NewList()
	// packages.SetBorder(true).SetTitle("Packages")
	// PopulatePackages(packages)
	pkgs := newTestPackages(t)

	// files := tview.NewList()
	// files.SetBorder(true).SetTitle("Files")
	// PopulateFiles(files)
	files := newTestFiles(t)

	tests := tview.NewList()
	tests.SetBorder(true).SetTitle("Tests")
	PopulateTests(tests)

	cases := tview.NewList()
	cases.SetBorder(true).SetTitle("Test Cases")
	PopulateCases(cases)

	// Create the results panel (right panel)
	results := tview.NewTextView()
	results.SetBorder(true).SetTitle("Results")
	RenderResults(results)

	// this is the navigations column made up of interactive panels
	navFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pkgs, 0, 1, true).
		AddItem(files, 0, 1, true).
		AddItem(tests, 0, 1, false).
		AddItem(cases, 0, 1, false)

	help := tview.NewTextView()
	// help.SetBorder(true).SetTitle("Keys").SetTitleAlign(tview.AlignLeft)
	help.SetLabel("/: search, q: quit, R: rerun last, r: run test, ?: more keys")

	// this is the whole screen
	contentLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(navFlex, 0, 1, true).
		AddItem(results, 0, 6, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 15, true).
		AddItem(help, 2, 1, false)

	t.app.SetRoot(layout, true)
}

// func (t *TUI) startMonitoring() {
// 	stop := make(chan int, 1)
// 	t.state.stopChans["task"] = stop
// 	t.state.stopChans["image"] = stop
// 	t.state.stopChans["volume"] = stop
// 	t.state.stopChans["network"] = stop
// 	t.state.stopChans["container"] = stop
// 	// go t.monitoringTask()
// 	go t.pkgPanel().monitoringPackages(t)
// 	// go t.networkPanel().monitoringNetworks(t)
// 	// go t.volumePanel().monitoringVolumes(t)
// 	// go t.containerPanel().monitoringContainers(t)
// }

func (t *TUI) pkgPanel() *testPackages {
	for _, panel := range t.state.panels.panel {
		if panel.name() == "Packages" {
			pnl, ok := panel.(*testPackages)
			if ok {
				return pnl
			}
		}
	}
	return nil
}

func (t *TUI) stopMonitoring() {
	t.state.stopChans["task"] <- 1
	t.state.stopChans["image"] <- 1
	t.state.stopChans["volume"] <- 1
	t.state.stopChans["network"] <- 1
	t.state.stopChans["container"] <- 1
}

// func (t *TUI) monitoringTask() {
// 	// common.Logger.Info("start monitoring task")
// LOOP:
// 	for {
// 		select {
// 		case task := <-t.taskPanel().tasks:
// 			go func() {
// 				if err := task.Func(task.Ctx); err != nil {
// 					task.Status = err.Error()
// 				} else {
// 					task.Status = success
// 				}
// 				t.updateTask()
// 			}()
// 		case <-t.state.stopChans["task"]:
// 			common.Logger.Info("stop monitoring task")
// 			break LOOP
// 		}
// 	}
// }
