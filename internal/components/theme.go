package components

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetAppStyling() {
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'
}

func SetFlexStyling(flex *tview.Flex) {
	colors := getTheme()
	flex.SetBackgroundColor(colors["blue"])
}

func SetTextViewStyling(txt *tview.TextView) {
	colors := getTheme()
	txt.SetBackgroundColor(colors["background"])
}

func SetListStyling(list *tview.List) {
	colors := getTheme()
	// list.SetBackgroundColor(colors["background"])
	// list.SetMainTextColor(colors["text"])
	// list.SetSecondaryTextColor(colors["blue"])
	// list.SetSelectedTextColor(colors["surface0"])
	// list.SetSelectedBackgroundColor(colors["pink"])
	list.SetBackgroundColor(colors["background"])
	list.SetMainTextColor(colors["text"])
	list.SetSecondaryTextColor(colors["blue"])
	list.SetSelectedTextColor(colors["blue"])
	list.SetSelectedBackgroundColor(colors["pink"])
}

// TODO: would be nice to get this from a config file
func getTheme() map[string]tcell.Color {
	// catppuccin := map[string]tcell.Color{
	// 	"background": HexToColor("#1E1E2E"),
	// 	"surface0":   HexToColor("#302D41"),
	// 	"text":       HexToColor("#D9E0EE"),
	// 	"green":      HexToColor("#ABE9B3"),
	// 	"blue":       HexToColor("#89B4FA"),
	// 	"pink":       HexToColor("#F5C2E7"),
	// }
	// return catppuccin

	testing := map[string]tcell.Color{
		"background": HexToColor("#1E1E2E"),
		"surface0":   HexToColor("#302D41"),
		"text":       HexToColor("#D9E0EE"),
		"green":      HexToColor("#ABE9B3"),
		"blue":       HexToColor("#89B4FA"),
		"pink":       HexToColor("#F5C2E7"),
	}

	return testing
}

func HexToColor(hex string) tcell.Color {
	color, err := strconv.ParseInt(hex[1:], 16, 32)
	if err != nil {
		fmt.Println("Invalid hex color")
		return tcell.ColorDefault
	}
	return tcell.NewHexColor(int32(color))
}
