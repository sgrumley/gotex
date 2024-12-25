package components

import (
	"sgrumley/gotex/pkg/ansi"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfigModal struct {
	modal   *tview.Flex
	content *tview.Table
}

func newConfigModal(t *TUI) *ConfigModal {
	cfg := &ConfigModal{}
	tab := tview.NewTable()
	cfg.content = tab
	cfg.Render(t, tab)

	// HACK: The table won't center well within it's flex, offset using 2,2,1 proportions
	modalContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 2, false).
		AddItem(tab, 0, 2, true).
		AddItem(nil, 0, 1, false)

	modal := NewModal("Config", modalContent)
	modal.SetBorder(true)
	cfg.modal = modal

	return cfg
}

func (m *ConfigModal) Render(t *TUI, tab *tview.Table) {
	h1 := tview.NewTableCell("Config")
	h1.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	tab.SetCell(0, 0, h1)

	h2 := tview.NewTableCell("Value")
	h2.SetAttributes(tcell.AttrBold | tcell.AttrUnderline)
	tab.SetCell(0, 1, h2)

	newRow(tab, "JSON:", t.state.data.project.Config.Json, 1)
	newRow(tab, "Timeout:", t.state.data.project.Config.Timeout, 2)
	newRow(tab, "Verbose:", t.state.data.project.Config.Verbose, 3)
	newRow(tab, "Cover:", t.state.data.project.Config.Cover, 4)
	newRow(tab, "Short:", t.state.data.project.Config.Short, 5)
	newRow(tab, "Fail Fast:", t.state.data.project.Config.FailFast, 6)
	newRow(tab, "Piped Command:", t.state.data.project.Config.PipeTo, 7)
}

func newRow(tab *tview.Table, key string, val interface{}, rowInd int) {
	cell1 := tview.NewTableCell(key)
	// cell1.SetAlign(tview.AlignCenter)
	// cell1.SetExpansion(1)
	tab.SetCell(rowInd, 0, cell1)

	cell2 := tview.NewTableCell(ansi.SimpleString(val))
	// cell2.SetExpansion(1)
	// cell2.SetAlign(tview.AlignCenter)
	tab.SetCell(rowInd, 1, cell2)
}

// Consider how to close
func (c *ConfigModal) setKeybindings(t *TUI) {
	c.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			// TODO: this should toggle
			t.state.ui.pages.SwitchToPage(homePage)
		}

		switch event.Rune() {
		case 'c':
			// TODO: this should toggle
			t.state.ui.pages.ShowPage(configPage)
			return nil
		}
		return event
	})
}
