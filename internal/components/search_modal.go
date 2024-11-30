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

	// TODO: integrate fzf
	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, test := range t.state.resources.flattened.Names {
			if strings.HasPrefix(strings.ToLower(test), strings.ToLower(currentText)) {
				entries = append(entries, test)
			}
		}
		if len(entries) <= 1 {
			entries = nil
		}
		return
	})
	input.SetDoneFunc(func(key tcell.Key) {
		searchStr := t.state.search.input.GetText()
		ref, exists := t.state.resources.flattened.NodeMap[searchStr]
		if !exists {
			t.log.Error("error", slog.String("search term not in tree", searchStr))
			// TODO: this modal is not the one we intend to store
			t.state.search.modal.SetBorderColor(tcell.ColorRed)
			return
		}

		found := search(t.state.testTree.TreeView, ref.GetName(), t)
		if !found {
			t.log.Error("error", slog.String("search term not found", searchStr))
			t.state.search.modal.SetBorderColor(tcell.ColorRed)
			return
		}

		t.state.search.modal.SetBorderColor(t.theme.Border)
		t.log.Info("search",
			slog.String("search term", searchStr),
			slog.String("ref", ref.GetName()),
			slog.Bool("found", found),
		)

		t.state.pages.SwitchToPage(homePage)
	})

	SetInputStyling(t, input)
	// TODO: nuffily keymaps
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
