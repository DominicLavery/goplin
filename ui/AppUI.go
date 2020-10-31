package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/logs"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

// Enum for screens & focus
const (
	NotebookTree = iota
	NoteTree
	NoteView
)

func MakeApp(source data.Source) *tview.Application {
	focusedView := 0
	cmdMode := false
	displayConsole := true

	app := tview.NewApplication()
	noteView := MakeNoteView(app, source)
	notesTree := MakeNotesTree(source)
	notebookTree := MakeNotebookView(source)
	cmdLine := MakeCmdLine(source)
	consoleView := logs.SetApp(app)

	var displays = []tview.Primitive{
		notebookTree,
		notesTree,
		noteView,
	}

	cmdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	viewFlex := tview.NewFlex()

	renderFlex := func() {
		flex.Clear()
		viewFlex.Clear()
		cmdFlex.Clear()

		flex.AddItem(viewFlex.
			AddItem(notebookTree, 0, 1, focusedView == NotebookTree).
			AddItem(notesTree, 0, 1, focusedView == NoteTree).
			AddItem(noteView, 0, 2, focusedView == NoteView), // Twice as big
			0, 1, !cmdMode)

		cmdFlexSize := 1
		if displayConsole {
			cmdFlexSize = 0
			cmdFlex.AddItem(consoleView, 0, 1, false)
		}
		cmdFlex.AddItem(cmdLine, 1, 1, cmdMode)
		flex.AddItem(cmdFlex, cmdFlexSize, 1, cmdMode)
		app.SetFocus(flex)
	}
	cmdLine.SetFinishedFunc(func(key tcell.Key) {
		cmdMode = false
		renderFlex()
	})
	renderFlex()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			focusedView = (focusedView + 1) % len(displays) // Next element w/ wrapping
			renderFlex()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			cmdMode = false
			cmdLine.SetText("")
			renderFlex()
			return nil
		}
		if event.Rune() == ':' {
			cmdMode = true
			cmdLine.SetText(":")
			renderFlex()
			return nil
		}
		if event.Rune() == 'c' {
			displayConsole = !displayConsole
			renderFlex()
			return nil
		}
		return event
	}).SetRoot(flex, true)

	return app
}
