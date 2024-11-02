package components

// import (
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// type testPackage struct {
// 	name string

// 	// list of files
// }

// type testPackages struct {
// 	*tview.List
// }

// func newTestPackages(t *TUI) *testPackages {
// 	pkgs := &testPackages{
// 		List: tview.NewList(),
// 	}

// 	pkgs.SetTitle("Packages")
// 	pkgs.SetBorder(true)
// 	pkgs.setKeybinding(t)
// 	pkgs.Populate()

// 	return pkgs
// }

// func (f *testPackages) setKeybinding(t *TUI) {
// 	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		t.setGlobalKeybinding(event)
// 		switch event.Key() {
// 		case tcell.KeyEnter:
// 			// TODO: run test
// 			f.AddItem("this should be added to resutls pane", "test completed", 'a', nil)
// 		case tcell.KeyCtrlR:
// 			// TODO: other events availible to files
// 		}

// 		// example using key instead of event
// 		switch event.Rune() {
// 		case 'R':
// 			f.AddItem("last test run", "test reran", 'a', nil)
// 			// case 'c':
// 		}

// 		return event
// 	})
// }

// func (p *testPackages) Populate() {
// 	p.AddItem("File A", "Details of package A", 'a', nil).
// 		AddItem("File B", "Details of package B", 'b', nil).
// 		AddItem("File C", "Details of package C", 'c', nil)
// }
