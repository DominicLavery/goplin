package ui

import (
	"dominiclavery/goplin/models"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
)

func MakeNoteView(app *tview.Application) (*tview.TextView, func(models.Note)) {
	noteView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	noteView.SetBorder(true).SetTitle("Note")
	updateNoteView := func(note models.Note) {
		noteView.Clear()
		source, err := ioutil.ReadFile(note.Path)
		if err != nil {
			source = []byte("Something went wrong, that file couldn't be opened")
			log.Println("Error during file reading file:", note.Path, "\nError:", err)
		}

		result := markdown.Render(string(source), 80, 1, markdown.WithImageDithering(markdown.DitheringWithChars))
		w := tview.ANSIWriter(noteView)
		fmt.Fprintf(w, "%s", result)
	}
	return noteView, updateNoteView
}
