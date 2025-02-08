package components

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	Name      string
	Project   tcell.Color
	Directory tcell.Color
	Package   tcell.Color
	File      tcell.Color
	Function  tcell.Color
	Case      tcell.Color
	tview.Theme
}

func getDefaultThemes() map[string]Theme {
	themes := map[string]Theme{
		// https://github.com/catppuccin/catppuccin/
		"catppuccin mocha": {
			Name: "catppuccin mocha",
			Theme: tview.Theme{
				PrimitiveBackgroundColor: HexToColor("#1E1E2E"), // Base
				BorderColor:              HexToColor("#b4befe"), // Lavender
				TitleColor:               HexToColor("#b4befe"), // Lavender
				PrimaryTextColor:         HexToColor("#cdd6f4"), // Text
				SecondaryTextColor:       HexToColor("#f5e0dc"), // Rosewater
			},
			Project:   HexToColor("#b4befe"), // Lavender
			Directory: HexToColor("#cdd6f4"), // Text
			Package:   HexToColor("#94e2d5"), // Teal
			File:      HexToColor("#89b4fa"), // Blue
			Function:  HexToColor("#cba6f7"), // Mauve
			Case:      HexToColor("#f5c2e7"), // Pink
		},
	}

	return themes
}

// TODO: Add more themes / custom themes
func SetTheme(name string) Theme {
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

func SetInputStyling(t *TUI, input *tview.InputField) {
	input.SetFieldTextColor(t.theme.File)
	autoCmpStyleMain := tcell.StyleDefault.Background(t.theme.PrimitiveBackgroundColor).Foreground(t.theme.Case)
	autoCmpStyleSelected := tcell.StyleDefault.Background(t.theme.Case).Foreground(t.theme.PrimitiveBackgroundColor)
	input.SetAutocompleteStyles(t.theme.PrimitiveBackgroundColor, autoCmpStyleMain, autoCmpStyleSelected)
}

func HexToColor(hex string) tcell.Color {
	color, err := strconv.ParseInt(hex[1:], 16, 32)
	if err != nil {
		return tcell.ColorDefault
	}
	return tcell.NewHexColor(int32(color))
}
