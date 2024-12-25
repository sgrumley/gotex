package components

import (
	"sgrumley/gotex/pkg/ansi"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfigModal struct {
	*tview.Modal
}

func newConfigModal(t *TUI) *ConfigModal {
	cfgModal := &ConfigModal{
		tview.NewModal(),
	}

	cfgModal.SetBorder(true)
	cfgModal.setKeybindings(t)
	// TODO: turn modal into a form for setting runtime config
	// use letters to jump to field??
	cfgModal.SetTitle("Current Config")
	cfgModal.Render(t)

	return cfgModal
}

func (m *ConfigModal) setKeybindings(t *TUI) {
	m.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.state.ui.pages.SwitchToPage(homePage)
		}

		// TODO: needs updating with page system
		switch event.Rune() {
		case 'c':
			name, _ := t.state.ui.pages.GetFrontPage()
			if name == homePage {
				t.state.ui.pages.ShowPage(configPage)
				return nil
			} else {
				t.state.ui.pages.SwitchToPage(homePage)
			}
			// if name == configPage {
			// 	t.state.ui.pages.SwitchToPage(homePage)
			// }
		}

		return nil
	})
}

func (m *ConfigModal) Render(t *TUI) {
	data := ansi.Data{
		Fields: []ansi.Field{
			ansi.CreateField("PipeTo", t.state.data.project.Config.PipeTo),
			ansi.CreateField("Timeout", t.state.data.project.Config.Timeout),
			ansi.CreateField("Json", t.state.data.project.Config.Json),
			ansi.CreateField("Short", t.state.data.project.Config.Short),
			ansi.CreateField("Verbose", t.state.data.project.Config.Verbose),
			ansi.CreateField("FailFast", t.state.data.project.Config.FailFast),
			ansi.CreateField("Cover", t.state.data.project.Config.Cover),
		},
	}
	m.SetText(ansi.OutputKeyVal(data))
}
