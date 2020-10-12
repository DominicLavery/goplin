package ui

import (
	"dominiclavery/goplin/models"
	"github.com/rivo/tview"
)

func MakeNotesTree(notes []models.Note, updateNoteView func(models.Note)) (*tview.Table, func([]models.Note)) {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	updateNoteTree := func(notes []models.Note) {
		table.Clear()
		if len(notes) > 0 {
			updateNoteView(notes[0])
			for i, note := range notes {
				table.SetCell(i, 0, tview.NewTableCell(note.Name)).SetBorder(true)
			}
		} else {
			table.SetCell(0, 0, tview.NewTableCell("No notes found")).SetBorder(true)
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

