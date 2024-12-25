package components

import (
	"github.com/gdamore/tcell/v2"
)

func (t *TUI) setGlobalKeybinding(_ *tcell.EventKey) {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		// TODO: move all binds into functions
		case 'R':
			_ = RerunTest(t)
			return nil
		case 'q':
			t.app.Stop()
		case 'C':
			toggleConsole(t)
			return nil
		case 'c':
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
