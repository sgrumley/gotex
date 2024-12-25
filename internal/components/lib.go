package components

import "github.com/rivo/tview"

// NewModal returns a tview.Flex with the provided content centered
func NewModal(name string, modalContent *tview.Flex) *tview.Flex {
	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	modal.SetBorder(true).
		SetTitle(name)

	modal.AddItem(modalContent, 0, 1, true)

	temp := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(modal, 0, 1, true).
		AddItem(nil, 0, 1, false)

	centeredModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(temp, 0, 1, true).
		AddItem(nil, 0, 1, false)

	return centeredModal
}
