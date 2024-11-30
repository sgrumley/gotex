package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type searchModal struct {
	modal *tview.Flex
	input *tview.InputField
}

func newSearchModal(t *TUI) *searchModal {
	input := tview.NewInputField().
		SetLabel("Search test: ").
		SetFieldWidth(20)

	SetInputStyling(t, input)

	modal := NewModal(t, input)
	modal.SetBorder(true)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.state.pages.SwitchToPage(homePage)
		}
		return event
	})
	return &searchModal{
		input: input,
		modal: modal,
	}
}

// NOTE: tview.Modal does not support embedding, this will use a flex to replicate the feature
func NewModal(t *TUI, input *tview.InputField) *tview.Flex {
	// Create a modal-like container
	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	modal.SetBorder(true).
		SetTitle("Search")

	SetFlexStyling(t, modal)

	textView := tview.NewTextView()
	textView.SetTextColor(tcell.ColorWhite)
	textView.SetBackgroundColor(t.theme.Background)
	textView.SetText("Welcome to Input Modal")

	// Create a flex to arrange content inside the modal
	modalContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 3, 1, true).
		AddItem(input, 3, 1, true)

	// Add the content to the modal
	modal.AddItem(modalContent, 0, 1, true)

	temp := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(modal, 0, 1, true).
		AddItem(nil, 0, 1, false)

	temp.SetBackgroundColor(t.theme.Background)

	// Create a centered container for the modal
	centeredModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(temp, 0, 3, true).
		AddItem(nil, 0, 1, false)

	return centeredModal
}
