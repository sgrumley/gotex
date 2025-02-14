package components

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"

	"github.com/sgrumley/gotex/pkg/slogger"
)

type searchModal struct {
	modal  *tview.Flex
	input  *tview.InputField
	text   *tview.TextView
	active bool
}

func newSearchModal(ctx context.Context, t *TUI) *searchModal {
	input := tview.NewInputField().
		SetLabel("Search test: ").
		SetFieldWidth(40)

	input.SetAutocompleteFunc(fuzzyFindTest(t))
	input.SetAutocompletedFunc(selectTest(ctx, t))
	// SetDoneFunc only executes if a field has not been autoselected, if this is hit it means that the user is searching for something that we know does not exist
	input.SetDoneFunc(noTest(t))

	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetText("Search for any test in the test tree")

	modalContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 2, 1, true).
		AddItem(input, 2, 1, true)

	SetInputStyling(t, input)
	modal := NewModal("Search", modalContent)
	modal.SetBorder(true)

	sm := &searchModal{
		input:  input,
		text:   textView,
		modal:  modal,
		active: false,
	}

	sm.setKeybindings(t)
	return sm
}

func (s *searchModal) setKeybindings(t *TUI) {
	s.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			toggleSearch(t)
			return nil
		}
		return event
	})
}

func toggleSearch(t *TUI) {
	if t.state.ui.search.active {
		t.state.ui.pages.SwitchToPage(homePage)
		t.app.SetFocus(t.state.ui.testTree.TreeView)
		t.state.ui.search.active = false

		t.state.ui.search.text.SetText("Search for any test in the test tree")
		t.state.ui.search.input.SetText("")
		return
	}

	t.state.ui.pages.ShowPage(searchPage)
	t.app.SetFocus(t.state.ui.search.modal)
	t.state.ui.search.active = true
}

func fuzzyFindTest(t *TUI) func(currentText string) (entries []string) {
	return func(currentText string) (entries []string) {
		tests := t.state.data.flattened.Names
		if len(currentText) == 0 {
			return
		}

		entries = fuzzy.FindNormalizedFold(currentText, tests)
		if len(entries) < 1 {
			entries = nil
		} else if len(entries) > 10 {
			entries = entries[:10]
		}
		return
	}
}

func selectTest(ctx context.Context, t *TUI) func(text string, index, source int) bool {
	return func(text string, index, source int) bool {
		log, err := slogger.FromContext(ctx)
		if err != nil {
			log, _ = slogger.New()
		}
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
				log.Error("search failed", fmt.Errorf("search term not in tree %s", searchStr))
				t.state.ui.search.text.SetText("search term does not exist in test tree: " + searchStr)
				return false
			}

			found := search(t.state.ui.testTree.TreeView, ref.GetName())
			if !found {
				log.Error("search failed", fmt.Errorf("search term not found: %s", searchStr))
				t.state.ui.search.modal.SetBorderColor(tcell.ColorRed)
				return false
			}

			log.Debug("search",
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
	}
}

func noTest(t *TUI) func(key tcell.Key) {
	return func(key tcell.Key) {
		searchStr := t.state.ui.search.input.GetText()
		errMsg := fmt.Sprintf("[red]Search term \"%s\" does not exist in the test tree.[-] \nPlease try again or press \"s\" from the main page to resync the files ", searchStr)
		t.state.ui.search.text.SetText(errMsg)
	}
}
