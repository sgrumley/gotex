package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *TUI) setGlobalKeybinding(event *tcell.EventKey) {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'R':
			// rerun last test
		case '/':
			// call modal
		case 'q':
			t.app.Stop()
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
			"cases": "/: search, q: quit, R: rerun last, r: run test",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
