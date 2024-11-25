package components

import (
	"fmt"

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
			// TODO: run global instead of in testTree
		case '/':
			// call search modal
		case 'q':
			t.app.Stop()
		case 'C':
			if t.state.console.active {
				t.state.result.RenderResults("c is working")
				t.state.console.active = false
				t.state.console.flex.RemoveItem(t.state.console.panel)
			} else {
				t.state.console.active = true

				t.state.console.flex.AddItem(t.state.console.panel, 8, 1, false)
			}
		case 'c':
			// TODO: this probably needs the implementation of pages
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Current Config: %#v", t.state.resources.data.Config))
			// AddButtons([]string{"Quit", "Cancel"}).
			// SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// 	if buttonLabel == "Quit" {
			// 		app.Stop()
			// 	}
			// })
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
