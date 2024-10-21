package components

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type testPackage struct {
	name string
	// list of files
}

type testPackages struct {
	*tview.List
}

func newTestPackages(t *TUI) *testPackages {
	pkgs := &testPackages{
		List: tview.NewList(),
	}

	pkgs.SetTitle("Packages")
	pkgs.SetBorder(true)
	pkgs.setKeybinding(t)
	pkgs.Populate()

	return pkgs
}

func (p *testPackages) focus(t *TUI) {
	t.app.SetFocus(p)
}

func (p *testPackages) unfocus() {
	// unsure
}

func (p *testPackages) name() string {
	return p.name()
}

func (p *testPackages) monitoringPackages(t *TUI) {
	// common.Logger.Info("start monitoring images")
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			// p.updateEntries(t)
			fmt.Println("TODO")
		case <-t.state.stopChans["image"]:
			ticker.Stop()
			break LOOP
		}
	}
	// common.Logger.Info("stop monitoring images")
}

func (f *testPackages) setKeybinding(t *TUI) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			// TODO: run test
			f.AddItem("key enter registered", "Details of test case A", 'a', nil).
				AddItem("File B", "Details of test case B", 'b', nil).
				AddItem("File C", "Details of test case C", 'c', nil)
		case tcell.KeyCtrlR:
			// TODO: other events availible to files
		}

		// example using key instead of event
		switch event.Rune() {
		case 'd':
			f.AddItem("key press registered", "Details of test case A", 'a', nil)
			// case 'c':
		}

		return event
	})
}

func (p *testPackages) Populate() {
	p.AddItem("File A", "Details of package A", 'a', nil).
		AddItem("File B", "Details of package B", 'b', nil).
		AddItem("File C", "Details of package C", 'c', nil)
}
