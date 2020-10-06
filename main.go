package main

import (
	"dominiclavery/goplin/models"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var notebooks = []models.Notebook{
	{"1", "", "."},
	{"2", "1", "Book 1"},
	{"3", "2", "Book 1.1"},
	{"4", "1", "Book 2"},
	{"5", "2", "Book 1.2"},
	{"6", "5", "Book 1.2.1"},
	{"7", "4", "Book 2.1"},
}

var notes = []models.Note{
	{"1", "1", "MyGreatNote1", "Stuff is cool1. Here are more words \n# Test"},
	{"2", "1", "MyGreatNote2", "Stuff is cool2"},
	{"3", "1", "MyGreatNote3", "Stuff is cool3"},
	{"4", "1", "MyGreatNote4", "Stuff is cool4"},
	{"5", "2", "MyGreatNote5", "Stuff is cool5"},
	{"6", "5", "subbooking", "Stuff is cool6"},
}

func makeTreeNode(book models.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func add(target *tview.TreeNode, parentId string) {
	for _, book := range notebooks {
		if book.ParentId == parentId {
			node := makeTreeNode(book)
			target.AddChild(node)
		}
	}
}

func makeNotebookView(updateNotesView func([]models.Note)) *tview.TreeView {
	notebookRoot := makeTreeNode(notebooks[0]).
		SetColor(tcell.ColorRed)

	add(notebookRoot, "1")
	updateNotesView(models.ByNotebookId(notes, "1"))

	notebookTree := tview.NewTreeView().
		SetRoot(notebookRoot).
		SetCurrentNode(notebookRoot).
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			updateNotesView(models.ByNotebookId(notes, reference.(string)))
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == "" {
				return // Selecting the notebookRoot node does nothing.
			}
			children := node.GetChildren()
			if len(children) == 0 {
				add(node, reference.(string))
			} else {
				node.SetExpanded(!node.IsExpanded())
			}
		})
	notebookTree.SetBorder(true).SetTitle("Notebooks")
	return notebookTree
}

func makeNotesTree(updateNoteView func(models.Note)) (*tview.Table, func([]models.Note)) {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	updateNoteTree := func(notes []models.Note) {
		table.Clear()
		if len(notes) > 0 {
			updateNoteView(notes[0])
		}
		for i, note := range notes {
			table.SetCell(i, 0, tview.NewTableCell(note.Name)).SetBorder(true)
		}
	}
	table.SetSelectionChangedFunc(func(row int, column int) {
		for i, note := range notes {
			if i == row {
				updateNoteView(note)
			}
		}
	})
	table.SetTitle("Notes")
	return table, updateNoteTree
}

func makeNoteView(app *tview.Application) (*tview.TextView, func(models.Note)) {
	noteView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	noteView.SetBorder(true).SetTitle("Note")
	updateNoteView := func(note models.Note) {
		noteView.Clear()
		fmt.Fprintf(noteView, "%s", note.Body)
	}
	return noteView, updateNoteView
}

func main() {
	app := tview.NewApplication()

	focused := 0
	noteView, updateNoteView := makeNoteView(app)
	notesTree, updateNotesTree := makeNotesTree(updateNoteView)
	notebookTree := makeNotebookView(updateNotesTree)

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
	})
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
