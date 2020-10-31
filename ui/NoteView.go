package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/logs"
	"dominiclavery/goplin/models"
	"fmt"
	"github.com/MichaelMure/go-term-markdown"
	"github.com/derailed/tview"
	"io/ioutil"
	"log"
)

func MakeNoteView(app *tview.Application, source data.Source) *tview.TextView {
	noteView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	noteView.SetBorder(true).SetTitle("Note")
	source.Note(func(note models.Note) {
		noteView.Clear()
		var buf []byte
		var err error
		if buf, err = ioutil.ReadAll(note.Body); err != nil {
			buf = []byte("Something went wrong, that file couldn't be opened")
			log.Println("Couldn't read a note", err)
		}
		result := markdown.Render(string(buf), 80, 1, markdown.WithImageDithering(markdown.DitheringWithChars))
		w := tview.ANSIWriter(noteView, "white", "black")

		if _, err = fmt.Fprintf(w, "%s", result); err != nil {
			logs.TeeLog("Error displaying the note", err)	
		}
		noteView.ScrollToBeginning()
	})
	return noteView
}
