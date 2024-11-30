package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createInputModal(app *tview.Application, pages *tview.Pages) (*tview.Flex, *tview.InputField) {
	// Define custom colors
	backgroundColor := tcell.ColorPink
	backgroundColor1 := tcell.ColorDarkGreen
	borderColor := tcell.ColorWhite

	// Create a text input field with custom colors
	input := tview.NewInputField()
	input.SetLabel("Enter your name: ")
	input.SetFieldWidth(20)
	input.SetLabelColor(tcell.ColorWhite)
	input.SetFieldTextColor(tcell.ColorWhite)
	input.SetBackgroundColor(backgroundColor1)

	// Create a modal-like container with custom background
	modal := tview.NewFlex()
	modal.SetDirection(tview.FlexColumn)
	modal.SetBorder(true)
	modal.SetBorderColor(borderColor)
	modal.SetBackgroundColor(backgroundColor1)
	modal.SetTitleColor(tcell.ColorWhite)
	modal.SetTitle("Input Modal")

	// Create text view
	titleView := tview.NewTextView()
	titleView.SetTextColor(tcell.ColorWhite)
	titleView.SetBackgroundColor(backgroundColor1)
	titleView.SetText("Welcome to Input Modal")

	// Create submit button
	submitButton := tview.NewButton("Submit")
	submitButton.SetSelectedFunc(func() {
		name := input.GetText()
		if name != "" {
			pages.HidePage("inputModal")
			fmt.Printf("Hello, %s!\n", name)
		}
	})
	submitButton.SetLabelColor(tcell.ColorWhite)
	submitButton.SetBackgroundColor(backgroundColor)

	// Create a flex to arrange content inside the modal
	modalContent := tview.NewFlex()
	modalContent.SetDirection(tview.FlexRow)
	modalContent.SetBackgroundColor(backgroundColor)
	modalContent.AddItem(titleView, 3, 1, false)
	modalContent.AddItem(input, 3, 1, true)
	modalContent.AddItem(submitButton, 3, 1, false)

	// Add the content to the modal
	modal.AddItem(modalContent, 0, 1, true)

	// Create a centered container for the modal
	centeredModal := tview.NewFlex()
	centeredModal.SetBackgroundColor(backgroundColor)
	centeredModal.AddItem(nil, 0, 1, false)
	centeredModal.AddItem(
		func() *tview.Flex {
			rowFlex := tview.NewFlex()
			rowFlex.SetDirection(tview.FlexRow)
			rowFlex.SetBackgroundColor(backgroundColor)
			rowFlex.AddItem(nil, 0, 1, false)
			rowFlex.AddItem(modal, 0, 1, true)
			rowFlex.AddItem(nil, 0, 1, false)
			return rowFlex
		}(), 0, 3, true,
	)
	centeredModal.AddItem(nil, 0, 1, false)

	return centeredModal, input
}

func main() {
	// Create the main application
	app := tview.NewApplication()

	// Create pages
	pages := tview.NewPages()

	// Your fullscreen page
	fullscreen := tview.NewFlex()
	fullscreen.AddItem(tview.NewTextView().SetText("Main Screen"), 0, 1, false)

	// Create input modal
	inputModal, inputField := createInputModal(app, pages)

	// Add pages
	pages.AddPage("fullscreen", fullscreen, true, true)
	pages.AddPage("inputModal", inputModal, true, false)

	// Optional: Add a keybinding to show the modal
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlI {
			pages.ShowPage("inputModal")
			app.SetFocus(inputField)
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			pages.HidePage("inputModal")
			return nil
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
