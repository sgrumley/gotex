package components

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	Name       string
	Background tcell.Color
	Border     tcell.Color
	Text       tcell.Color
	Project    tcell.Color
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
			Text:       HexToColor("#b4befe"), // Lavender
			Project:    HexToColor("#f5e0dc"),
			File:       HexToColor("#f5e0dc"),
			Function:   HexToColor("#f5e0dc"),
			Case:       HexToColor("#f5e0dc"),
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
	// tree.SetBorderStyle() TODO: rounded corner option
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
}

func HexToColor(hex string) tcell.Color {
	color, err := strconv.ParseInt(hex[1:], 16, 32)
	if err != nil {
		fmt.Println("Invalid hex color")
		return tcell.ColorDefault
	}
	return tcell.NewHexColor(int32(color))
}

var catppuccin = map[string]tcell.Color{
	"background": HexToColor("#1E1E2E"),
	"surface0":   HexToColor("#302D41"),
	"text":       HexToColor("#D9E0EE"),
	"green":      HexToColor("#ABE9B3"),
	"blue":       HexToColor("#89B4FA"),
	"pink":       HexToColor("#F5C2E7"),
}
