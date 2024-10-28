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
			t.nextPanel()
		case 'h':
			t.prevPanel()
		}

		return event
	})
}

// prompt fzf code in pop up window
func (t *TUI) search() {
	// get a popup modal for searching via fzf
}

// NOTE: navigating panels might require some middle ware for rendering
func (t *TUI) nextPanel() {
	switch t.state.panels.currentPanel {
	case "files":
		t.app.SetFocus(t.state.panels.panel["tests"].GetList())
		t.state.panels.currentPanel = "tests"

	case "tests":
		t.app.SetFocus(t.state.panels.panel["cases"].GetList())
		t.state.panels.currentPanel = "cases"

	case "cases":
		t.app.SetFocus(t.state.panels.panel["files"].GetList())
		t.state.panels.currentPanel = "files"
	}
}

func (t *TUI) prevPanel() {
	switch t.state.panels.currentPanel {
	case "files":
		t.app.SetFocus(t.state.panels.panel["cases"].GetList())
		t.state.panels.currentPanel = "cases"

	case "tests":
		t.app.SetFocus(t.state.panels.panel["files"].GetList())
		t.state.panels.currentPanel = "files"

	case "cases":
		t.app.SetFocus(t.state.panels.panel["tests"].GetList())
		t.state.panels.currentPanel = "tests"

	}
}

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
