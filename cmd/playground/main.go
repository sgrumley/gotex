package main

// NOTE: this is an ongoing WIP for improving the search feature
import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createInputWithSuggestionsUI() *tview.Flex {
	suggestions := []string{"apple", "application", "apricot", "Banana", "Cherry", "Date", "Elderberry"}
	currentSuggestionIndex := -1
	currentEntries := []string{}

	// Create input field
	inputField := tview.NewInputField()
	inputField.SetLabel("Fruits: ").
		SetFieldWidth(20)

	// Create list for suggestions
	suggestionList := tview.NewList()
	suggestionList.SetBorder(true).
		SetTitle("Suggestions")

	// Create flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(inputField, 3, 1, true)
	flex.AddItem(suggestionList, 10, 1, false)

	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			suggestionList.Clear()
			return
		}

		// Filter suggestions
		entries = []string{}
		for _, test := range suggestions {
			if strings.HasPrefix(strings.ToLower(test), strings.ToLower(currentText)) {
				entries = append(entries, test)
			}
		}

		// Update suggestion list
		suggestionList.Clear()
		for i, entry := range entries {
			index := i
			suggestionList.AddItem(entry, "", 0, func() {
				inputField.SetText(entries[index])
			})
		}

		currentEntries = entries
		currentSuggestionIndex = -1

		if len(entries) <= 1 {
			entries = nil
		}
		return entries
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle suggestion navigation when suggestions exist
		if len(currentEntries) > 0 {
			switch event.Key() {
			case tcell.KeyDown:
				fallthrough
			case tcell.KeyCtrlJ:
				// Move to next suggestion
				currentSuggestionIndex++
				if currentSuggestionIndex >= len(currentEntries) {
					currentSuggestionIndex = 0
				}
				inputField.SetText(currentEntries[currentSuggestionIndex])
				suggestionList.SetCurrentItem(currentSuggestionIndex)
				return nil

			case tcell.KeyUp:
				fallthrough
			case tcell.KeyCtrlK:
				// Move to previous suggestion
				currentSuggestionIndex--
				if currentSuggestionIndex < 0 {
					currentSuggestionIndex = len(currentEntries) - 1
				}
				inputField.SetText(currentEntries[currentSuggestionIndex])
				suggestionList.SetCurrentItem(currentSuggestionIndex)
				return nil

			case tcell.KeyTab:
				// Allow switching focus between input and suggestion list
				app.SetFocus(suggestionList)
				return nil
			}
		}
		return event
	})

	// Add handler to suggestion list to update input when item is selected
	suggestionList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// When Enter is pressed on suggestion list, update input
			_, suggestion := suggestionList.GetItemText(suggestionList.GetCurrentItem())
			inputField.SetText(suggestion)
			app.SetFocus(inputField)
			return nil
		case tcell.KeyTab:
			// Allow switching back to input field
			app.SetFocus(inputField)
			return nil
		}
		return event
	})

	return flex
}

var app *tview.Application

func main() {
	app = tview.NewApplication()
	flex := createInputWithSuggestionsUI()

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
