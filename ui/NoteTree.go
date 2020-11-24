package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/models"
	"github.com/derailed/tview"
)

type NotesTree struct {
	*tview.Table
	source data.NotebookReader
}

func (nt *NotesTree) SetNotes(notes []models.Note) {
	nt.Clear()
	if len(notes) > 0 {
		for i, note := range notes {
			nt.SetCell(i, 0, tview.NewTableCell(note.Name).SetReference(note.Id)).SetBorder(true)
		}
	} else {
		nt.SetCell(0, 0, tview.NewTableCell("No notes found")).SetBorder(true)
	}
}

func MakeNotesTree(source data.NotebookReader) *NotesTree {
	table := NotesTree{Table: tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false), source: source}

	table.SetSelectionChangedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		data.OpenNoteChan <- cell.GetReference().(int)
	})
	table.SetTitle("Notes")
	return &table
}
