package components

import (
	"github.com/gdamore/tcell/v2"
)

func (t *TUI) setGlobalKeybinding(event *tcell.EventKey) {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'q':
			t.app.Stop()
		}

		if event.Key() == tcell.KeyRune && event.Rune() == 'l' {
			// focusedFlex := t.app.GetFocus()
			// switch focused {
			// case textView1:
			// 	t.app.SetFocus(textView2) // Move to the next item
			// case textView2:
			// 	t.app.SetFocus(textView3)
			// case textView3:
			// 	t.app.SetFocus(textView1) // Loop back to the first item
			// }
			return nil
		}
		return event
	})

	// switch event.Rune() {
	// // case 'h':
	// // 	t.prevPanel()
	// // case 'l':
	// // 	t.nextPanel()
	// case 'q':
	// 	t.Stop()
	// case '/':
	// 	t.search()
	// }

	// switch event.Key() {
	// case tcell.KeyTab:
	// 	t.nextPanel()
	// case tcell.KeyBacktab:
	// 	t.prevPanel()
	// case tcell.KeyRight:
	// 	t.nextPanel()
	// case tcell.KeyLeft:
	// 	t.prevPanel()
	// }
	// switch event.Key() {
	// case tcell.KeyTab
	// }
}

// prompt fzf code in pop up window
func (t *TUI) search() {
	currentPanel := t.state.panels.panel[t.state.panels.currentPanel]
	if currentPanel.name() == "tasks" {
		return
	}
	// currentPanel.setFilterWord("")
	// currentPanel.updateEntries(t)

	// viewName := "filter"
	// searchInput := tview.NewInputField().SetLabel("Word")
	// searchInput.SetLabelWidth(6)
	// searchInput.SetTitle("filter")
	// searchInput.SetTitleAlign(tview.AlignLeft)
	// searchInput.SetBorder(true)

	// closeSearchInput := func() {
	// 	t.closeAndSwitchPanel(viewName, t.state.panels.panel[t.state.panels.currentPanel].name())
	// }

	// searchInput.SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyEnter {
	// 		closeSearchInput()
	// 	}
	// })

	// searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Key() == tcell.KeyEsc {
	// 		closeSearchInput()
	// 	}
	// 	return event
	// })

	// searchInput.SetChangedFunc(func(text string) {
	// 	currentPanel.setFilterWord(text)
	// 	currentPanel.updateEntries(t)
	// })

	// t.pages.AddAndSwitchToPage(viewName, t.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (t *TUI) nextPanel() {
	idx := (t.state.panels.currentPanel + 1) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}

func (t *TUI) prevPanel() {
	t.state.panels.currentPanel--

	if t.state.panels.currentPanel < 0 {
		t.state.panels.currentPanel = len(t.state.panels.panel) - 1
	}

	idx := (t.state.panels.currentPanel) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}

func (t *TUI) switchPanel(panelName string) {
	for i, panel := range t.state.panels.panel {
		if panel.name() == panelName {
			t.state.navigate.update(panelName)
			panel.focus(t)
			t.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}
