package components

import (
	"log/slog"
	"strings"

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
		SetFieldWidth(40)

	// TODO: fix theming
	// input.SetAutocompleteStyles()

	// TODO: integrate fzf
	// custom autocompleteFunc -> setup an input handler that takes the keys input (channel) and feeds to an output func (requires gui)
	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		tests := t.state.resources.flattened.Names
		if len(currentText) == 0 {
			return
		}

		for _, test := range tests {
			if strings.HasPrefix(strings.ToLower(test), strings.ToLower(currentText)) {
				entries = append(entries, test)
			}
		}
		if len(entries) <= 1 {
			entries = nil
		}
		return
	})
	input.SetAutocompletedFunc(func(text string, index, source int) bool {
		switch source {
		case 0:
			// navigate
			return false
		case 1:
			// tab key
			fallthrough
		case 2:
			// enter key
			fallthrough
		case 3:
			// mouse click
			searchStr := text
			ref, exists := t.state.resources.flattened.NodeMap[searchStr]
			if !exists {
				t.log.Error("error", slog.String("search term not in tree", searchStr))
				// TODO: this modal is not the one we intend to store
				t.state.search.modal.SetBorderColor(tcell.ColorRed)
				return false
			}

			found := search(t.state.testTree.TreeView, ref.GetName())
			if !found {
				t.log.Error("error", slog.String("search term not found", searchStr))
				t.state.search.modal.SetBorderColor(tcell.ColorRed)
				return false
			}

			t.state.search.modal.SetBorderColor(t.theme.Border)
			t.log.Info("search",
				slog.String("search term", searchStr),
				slog.String("ref", ref.GetName()),
				slog.Bool("found", found),
			)

			t.state.pages.SwitchToPage(homePage)
			return true
		default:
			return false
		}
	})

	SetInputStyling(t, input)
	// TODO: nuffily keymaps
	modal := NewModal(t, input)
	modal.SetBorder(true)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
			// nav down
		case tcell.KeyCtrlK:
			tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
			// nav down

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
	textView.SetText("Welcome to Input Modal, this is dummy text for now")

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
