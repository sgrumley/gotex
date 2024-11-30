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

	SetModalStyling(t, cfgModal.Modal)
	// TODO: turn modal into a form for setting runtime config
	// use letters to jump to field??
	cfgModal.SetTitle("Current Config")
	cfgModal.Render(t)
	cfgModal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.state.pages.SwitchToPage(homePage)
		}
		return event
	})

	return cfgModal
}

func (m *ConfigModal) Render(t *TUI) {
	data := ansi.Data{
		Fields: []ansi.Field{
			ansi.CreateField("PipeTo", t.state.resources.data.Config.PipeTo),
			ansi.CreateField("Timeout", t.state.resources.data.Config.Timeout),
			ansi.CreateField("Json", t.state.resources.data.Config.Json),
			ansi.CreateField("Short", t.state.resources.data.Config.Short),
			ansi.CreateField("Verbose", t.state.resources.data.Config.Verbose),
			ansi.CreateField("FailFast", t.state.resources.data.Config.FailFast),
			ansi.CreateField("Cover", t.state.resources.data.Config.Cover),
		},
	}
	m.SetText(ansi.OutputKeyVal(data))
}
