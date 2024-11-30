package components

import (
	"fmt"

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
	res.RenderConsole(t, ConsoleTemplate())
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

// TODO: key val data is probably better displayed as column or This row should be split into 2 columns
func ConsoleTemplate() string {
	return fmt.Sprintf(
		"THIS IS DUMMY DATA\nTest Name: %s\nCommand: %s\nStatus: %s\nCompleted at: %s\nLogger Filepath: %s\nTest Location: %s\n",
		"success/valid_json",
		"go test -run method/success",
		"[green]Pass[-]",
		"10:14",
		"~/.config/gotex/log.json",
		"./test/go_test.go",
	)
}

func (r *console) RenderConsole(t *TUI, msg string) {
	r.Clear()
	t.state.console.currentMessage = msg
	msg = tview.TranslateANSI(msg)
	r.SetDynamicColors(true)
	r.SetText(msg)
}
