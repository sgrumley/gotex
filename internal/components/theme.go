package components

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	Name       string
	Background tcell.Color
	Surface    tcell.Color
	Border     tcell.Color
	Text       tcell.Color
	Project    tcell.Color
	Package    tcell.Color
	File       tcell.Color
	Function   tcell.Color
	Case       tcell.Color
}

func getDefaultThemes() map[string]Theme {
	themes := map[string]Theme{
		// https://github.com/catppuccin/catppuccin/
		"catppuccin mocha": {
			Name:       "catppuccin mocha",
			Background: HexToColor("#1E1E2E"), // Base
			Border:     HexToColor("#b4befe"), // Lavender
			Text:       HexToColor("#cdd6f4"),
			Project:    HexToColor("#b4befe"),
			Package:    HexToColor("#94e2d5"),
			File:       HexToColor("#89b4fa"),
			Function:   HexToColor("#cba6f7"),
			Case:       HexToColor("#f5c2e7"),
			Surface:    HexToColor("#303446"),
		},
	}
	return themes
}

// TODO: have a set of default themes to choose from
// allow custom from config
func SetTheme(name string) Theme {
	// TODO: check for custom themes or return a native theme
	themes := getDefaultThemes()
	if name == "" {
		return themes["catppuccin mocha"]
	}

	return themes[name]
}

func SetAppStyling() {
	// TODO: can this happen on the panels only??
	// tview.Borders.TopLeft = '╭'
	// tview.Borders.TopRight = '╮'
	// tview.Borders.BottomLeft = '╰'
	// tview.Borders.BottomRight = '╯'
}

// NOTE: node specific styling done in populate since it is applied as the tree is rendered
func SetTreeStyling(t *TUI, tree *tview.TreeView) {
	tree.SetBorder(true)
	tree.SetTitleColor(t.theme.Border)
	tree.SetBorderColor(t.theme.Border)
	tree.SetBackgroundColor(t.theme.Background)
	// tree.SetBorderStyle() // TODO: rounded corner option
}

func SetTextViewStyling(t *TUI, txt *tview.TextView) {
	txt.SetBackgroundColor(t.theme.Background)
	txt.SetBorderColor(t.theme.Border)
	txt.SetTextColor(t.theme.Text)
	txt.SetTitleColor(t.theme.Border)
}

func SetFlexStyling(t *TUI, flex *tview.Flex) {
	flex.SetBackgroundColor(t.theme.Background)
	flex.SetBorderColor(t.theme.Border)
	flex.SetTitleColor(t.theme.Border)
}

func SetBoxStyling(t *TUI, box *tview.Box) {
	box.SetBackgroundColor(t.theme.Background)
	box.SetBorder(true)
	box.SetBorderColor(t.theme.Border)
	box.SetTitleColor(t.theme.Border)
}

func SetInputStyling(t *TUI, input *tview.InputField) {
	input.SetFieldBackgroundColor(t.theme.Background)
	input.SetLabelColor(t.theme.Border)
	input.SetFieldTextColor(t.theme.File)
	input.SetBackgroundColor(t.theme.Background)
}

func SetModalStyling(t *TUI, modal *tview.Modal) {
	modal.SetBackgroundColor(t.theme.Background)
	modal.SetTitleColor(t.theme.Border)
	modal.SetBorder(true)
	modal.SetBorderColor(t.theme.Case)
	// NOTE: not sure why only the modal needs this level of granularity
	modal.SetBorderStyle(
		tcell.StyleDefault.Foreground(t.theme.Border).
			Background(t.theme.Background),
	)
}

func HexToColor(hex string) tcell.Color {
	color, err := strconv.ParseInt(hex[1:], 16, 32)
	if err != nil {
		return tcell.ColorDefault
	}
	return tcell.NewHexColor(int32(color))
}
