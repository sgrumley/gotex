package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *TUI) setGlobalKeybinding(event *tcell.EventKey) {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		// TODO: move j and k to testTree
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'R':
			// rerun last test
		case 'q':
			t.app.Stop()
		case 'C':
			if t.state.console.active {
				t.state.console.active = false
				t.state.console.flex.RemoveItem(t.state.console.panel)
				return nil
			} else {
				t.state.console.active = true
				t.state.console.flex.AddItem(t.state.console.panel, 8, 1, false)
				return nil
			}
		case 'c':
			// SwitchToPage will hide all other pages
			// t.state.pages.SwitchToPage(configPage)
			t.state.pages.ShowPage(configPage)
			return nil
		// TODO: move below to testTree
		case 'l':
			// expand
		case 'h':
			// collapse
		}

		return event
	})
}

// prompt fzf code in pop up window
// func (t *TUI) search() {
// get a popup modal for searching via fzf
// }

// navigate will hold help text for different panels when needed
type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"testTree": "/: search, r: run test, R: rerun last test, s: sync project, q: quit",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
