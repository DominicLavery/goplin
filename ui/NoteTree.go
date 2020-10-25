package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/models"
	"github.com/derailed/tview"
)

func MakeNotesTree(source data.Source) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	source.Notes(func(notes []models.Note) {
		table.Clear()
		if len(notes) > 0 {
			source.OpenNote(notes[0].Id)
			for i, note := range notes {
				table.SetCell(i, 0, tview.NewTableCell(note.Name).SetReference(note.Id)).SetBorder(true)
			}
		} else {
			table.SetCell(0, 0, tview.NewTableCell("No notes found")).SetBorder(true)
		}
	})
	table.SetSelectionChangedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		source.OpenNote(cell.GetReference().(int))
	})
	table.SetTitle("Notes")
	return table
}
