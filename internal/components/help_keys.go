package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"cases":     "/: search, q: quit, R: rerun last, r: run test",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
