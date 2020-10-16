package ui

import (
	"dominiclavery/goplin/models"
	"github.com/gdamore/tcell"
	"github.com/derailed/tview"
)

// Enum for screens & focus
const (
	NotebookTree = iota
	NoteTree
	NoteView
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

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	renderFlex := func() {
		flex.Clear()
		flex.AddItem(tview.NewFlex().
			AddItem(notebookTree, 0, 1, focused == NotebookTree).
			AddItem(notesTree, 0, 1, focused == NoteTree).
			AddItem(noteView, 0, 2, focused == NoteView), // Twice as big
		0, 1, true)
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
