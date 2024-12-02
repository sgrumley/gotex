package components

import (
	"log/slog"

	"github.com/gdamore/tcell/v2"
)

func (t *TUI) setGlobalKeybinding(_ *tcell.EventKey) {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		// TODO: move all binds into functions
		// FIX: need a way to show the user that the test has been rerun/ is rerunning
		// maybe a job for the meta console?
		case 'R':
			// rerun last test
			// kept as a global to allow user in other interfaces e.g. config
			t.state.ui.result.RenderResults("Rerunning test")
			t.log.Error("this should not have run")

			node := t.state.data.lastTest
			if node == nil {
				t.state.ui.result.RenderResults("failed to run last test. Make sure you run a test before rerunning")
				t.log.Error("attempted test rerun, but no test has previously been run")
				return event
			}

			output, err := node.RunTest()
			if err != nil {
				t.log.Error("failed to re run valid test", slog.Any("error", err))
				t.state.ui.console.panel.UpdateMeta(t, output)
				t.state.ui.result.RenderResults(err.Error())
				return event
			}
			t.state.ui.result.RenderResults(output.Result)
			t.state.ui.console.panel.UpdateMeta(t, output)
			return nil

		case 'q':
			t.app.Stop()
		case 'C':
			if t.state.ui.console.active {
				t.state.ui.console.active = false
				t.state.ui.console.flex.RemoveItem(t.state.ui.console.panel)
				return nil
			} else {
				t.state.ui.console.active = true
				t.state.ui.console.flex.AddItem(t.state.ui.console.panel, 0, 10, false)
				return nil
			}
		case 'c':
			// SwitchToPage will hide all other pages
			// t.state.pages.SwitchToPage(configPage)
			t.state.ui.pages.ShowPage(configPage)
			return nil
		}
		return event
	})
}

// TODO: get this working again
// navigate will hold help text for different panels when needed
// type navigate struct {
// 	*tview.TextView
// 	keybindings map[string]string
// }

// func newNavigate() *navigate {
// 	return &navigate{
// 		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
// 		keybindings: map[string]string{
// 			"testTree": "/: search, r: run test, R: rerun last test, s: sync project, q: quit",
// 		},
// 	}
// }

// func (n *navigate) update(panel string) {
// 	n.SetText(n.keybindings[panel])
// }
