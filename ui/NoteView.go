package ui

import (
	"dominiclavery/goplin/logs"
	"fmt"
	"github.com/MichaelMure/go-term-markdown"
	"github.com/derailed/tview"
	"io"
	"io/ioutil"
	"log"
)

type NoteView struct {
	*tview.TextView
}

func (nv *NoteView) SetNote(note io.Reader) {
	nv.Clear()
	var buf []byte
	var err error
	if buf, err = ioutil.ReadAll(note); err != nil {
		buf = []byte("Something went wrong, that file couldn't be opened")
		log.Println("Couldn't read a note", err)
	}
	result := markdown.Render(string(buf), 80, 1, markdown.WithImageDithering(markdown.DitheringWithChars))
	w := tview.ANSIWriter(nv, "white", "black")

	if _, err = fmt.Fprintf(w, "%s", result); err != nil {
		logs.TeeLog("Error displaying the note", err)
	}
	nv.ScrollToBeginning()
	if v, ok := note.(io.Closer); ok {
		_ = v.Close()
	}
}

func MakeNoteView() *NoteView {
	nv := NoteView{TextView: tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true),
	}
	nv.SetBorder(true).SetTitle("Note")
	return &nv
}
