package ui

import (
	"dominiclavery/goplin/models"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func MakeApp(notebooks []models.Notebook, notes []models.Note) *tview.Application {
	app := tview.NewApplication()

	focused := 0
	noteView, updateNoteView := MakeNoteView(app)
	notesTree, updateNotesTree := MakeNotesTree(notes, updateNoteView)
	notebookTree := MakeNotebookView(notebooks, notes, updateNotesTree)

	var displays = []tview.Primitive{
		notebookTree,
		notesTree,
		noteView,
	}

	flex := tview.NewFlex()
	renderFlex := func() {
		flex.Clear()
		for i, display := range displays {
			flex.AddItem(display, 0, 1, i == focused)
		}
		app.SetFocus(flex)
	}

	renderFlex()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			focused = (focused + 1) % len(displays) // Next element w/ wrapping
			renderFlex()
			return nil
		}
		return event
	}).SetRoot(flex, true)

	return app
}
