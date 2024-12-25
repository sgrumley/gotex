package components

import (
	"fmt"
	"log/slog"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	// "github.com/sahilm/fuzzy"
)

type searchModal struct {
	modal *tview.Flex
	input *tview.InputField
	text  *tview.TextView
}

func newSearchModal(t *TUI) *searchModal {
	input := tview.NewInputField().
		SetLabel("Search test: ").
		SetFieldWidth(40)

	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		tests := t.state.data.flattened.Names
		if len(currentText) == 0 {
			return
		}

		entries = fuzzy.Find(currentText, tests)
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
			ref, exists := t.state.data.flattened.NodeMap[searchStr]
			if !exists {
				t.log.Error("error", slog.String("search term not in tree", searchStr))
				t.state.ui.search.text.SetText("search term does not exist in test tree: " + searchStr)
				return false
			}

			found := search(t.state.ui.testTree.TreeView, ref.GetName())
			if !found {
				t.log.Error("error", slog.String("search term not found", searchStr))
				t.state.ui.search.modal.SetBorderColor(tcell.ColorRed)
				return false
			}

			// t.state.ui.search.modal.SetBorderColor(t.theme.Border)
			t.log.Info("search",
				slog.String("search term", searchStr),
				slog.String("ref", ref.GetName()),
				slog.Bool("found", found),
			)

			t.state.ui.search.input.SetText("")
			t.state.ui.pages.SwitchToPage(homePage)
			return true
		default:
			return false
		}
	})
	// SetDoneFunc only executes if a field has not been autoselected, if this is hit it means that the user is searching for something that we know does not exist
	input.SetDoneFunc(func(key tcell.Key) {
		searchStr := t.state.ui.search.input.GetText()
		errMsg := fmt.Sprintf("[red]Search term \"%s\" does not exist in the test tree.[-] \nPlease try again or press \"s\" from the main page to resync the files ", searchStr)
		t.state.ui.search.text.SetText(errMsg)
	})

	textView := tview.NewTextView()
	// textView.SetTextColor(t.theme.Text)
	// textView.SetBackgroundColor(t.theme.Background)
	textView.SetDynamicColors(true)
	textView.SetText("Search for any test in the test tree")

	SetInputStyling(t, input)
	modal := NewModal(t, input, textView)
	modal.SetBorder(true)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// NOTE: looks like navigation here will require wrapping or reimplementing the inputfield type
		// or making a custom type

		// case tcell.KeyCtrlJ:
		// tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		// nav down
		// case tcell.KeyCtrlK:
		// tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		// nav down

		case tcell.KeyEsc:
			t.setGlobalKeybinding(event)
			t.state.ui.search.input.SetText("")
			t.state.ui.search.text.SetText("Search for any test in the test tree")
			t.state.ui.pages.SwitchToPage(homePage)
		}
		return event
	})
	return &searchModal{
		input: input,
		text:  textView,
		modal: modal,
	}
}

// NOTE: tview.Modal does not support embedding, this will use a flex to replicate the feature
func NewModal(t *TUI, input *tview.InputField, textView *tview.TextView) *tview.Flex {
	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	modal.SetBorder(true).
		SetTitle("Search")

	modalContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 3, 1, true).
		AddItem(input, 3, 1, true)

	modal.AddItem(modalContent, 0, 1, true)

	temp := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(modal, 0, 1, true).
		AddItem(nil, 0, 1, false)

	centeredModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(temp, 0, 3, true).
		AddItem(nil, 0, 1, false)

	return centeredModal
}
