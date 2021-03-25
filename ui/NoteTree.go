package ui

import (
	"dominiclavery/goplin/data"
	"github.com/derailed/tview"
	"github.com/google/uuid"
)

type NotesTree struct {
	*tview.Table
}

func (nt *NotesTree) SetNotes(notes []*data.Note) {
	nt.Clear()
	if len(notes) > 0 {
		for i, note := range notes {
			nt.SetCell(i, 0, tview.NewTableCell(note.Name).SetReference(note.Id)).SetBorder(true)
		}
	} else {
		nt.SetCell(0, 0, tview.NewTableCell("No notes found")).SetBorder(true)
	}
	nt.Select(0, 0)
}

func MakeNotesTree() *NotesTree {
	table := NotesTree{Table: tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)}

	table.SetSelectionChangedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		if cell.GetReference() != nil {
			note := data.GetNote(cell.GetReference().(uuid.UUID))
			noteView.SetNote(note)
		} else {
			noteView.Clear()
		}
	})
	table.SetTitle("Notes")
	return &table
}
