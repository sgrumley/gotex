package components

import (
	"sgrumley/gotex/pkg/ansi"
	"sgrumley/gotex/pkg/runner"

	"github.com/rivo/tview"
)

type console struct {
	*tview.TextView
}

func newConsolePane(t *TUI) *console {
	res := &console{
		TextView: tview.NewTextView(),
	}

	res.SetBorder(true).SetTitle("Console")
	res.SetTextAlign(tview.AlignLeft)
	res.RenderConsole(t, "[green]Run[-] a test to see the meta data")
	res.SetDynamicColors(true)
	SetTextViewStyling(t, res.TextView)
	res.SetWrap(true)

	return res
}

/*
Test Name: "success"
Command: "go test -run method/success"
Completed at: "ts.Now()"
Status: "[green]Pass[-]"
Logger Filepath: ~/.config/gotex/log.json
Test Location: ./test/go_test.go
*/

func (c *console) RenderConsole(t *TUI, msg string) {
	c.Clear()
	t.state.ui.console.currentMessage = msg
	msg = tview.TranslateANSI(msg)
	c.SetDynamicColors(true)
	c.SetText(msg)
}

func (c *console) UpdateMeta(t *TUI, meta *runner.Response) {
	// HACK: text is janky if not new lined
	if meta.Output != "" {
		meta.Output = "\n" + meta.Output
	}
	if meta.Error != "" {
		meta.Error = "\n" + meta.Error
	}
	if meta.ExternalOutput != "" {
		meta.ExternalOutput = "\n" + meta.ExternalOutput
	}
	if meta.ExternalError != "" {
		meta.ExternalError = "\n" + meta.ExternalError
	}
	data := ansi.Data{
		Fields: []ansi.Field{
			ansi.CreateField("Name", meta.TestName),
			ansi.CreateField("Command Executed", meta.CommandExecuted),
			ansi.CreateField("Execution Filepath", meta.TestDir),
			ansi.CreateField("Type (1=project,2=package,3=file,4=function,5=case)", meta.TestType),
			ansi.CreateField("Piped to external", meta.External),
			ansi.CreateField("Output", meta.Output),
			ansi.CreateField("Error", meta.Error),
			ansi.CreateField("External Output", meta.ExternalOutput),
			ansi.CreateField("External Error", meta.ExternalError),
		},
	}
	c.RenderConsole(t, ansi.OutputKeyVal(data))
}
