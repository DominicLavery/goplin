package data

import (
	"dominiclavery/goplin/models"
	"errors"
	"path/filepath"
	"strings"
	"time"
)

type DummySource struct {
	NotebookReader
	NotebookWriter
}

type DummyReader struct {
	requestedBook int
	requestedNote int
	openNote      int
	openBook      int
	notebooks     Notebooks
	notes         Notes
}

type DummyWriter struct {
}

const dummyText = `

# Title Level 1
## Title Level 2
|test|test2|
|---|---|
|Details|more details|
`

func NewDummySource() *DummySource {
	writer := &DummyWriter{}
	reader := &DummyReader{requestedBook: 0, requestedNote: 0, openNote: -1, openBook: -1}
	fs := &DummySource{
		NotebookReader: reader,
		NotebookWriter: writer,
	}

	reader.notebooks = Notebooks{
		notebookRoot: models.Notebook{
			Id: RootId, ParentId: -1, Name: "rootBook", Children: []models.Notebook{
				{Id: 1, ParentId: 0, Name: "Book 1", Children: []models.Notebook{
					{Id: 2, ParentId: 1, Name: "Book 1.1"},
					{Id: 4, ParentId: 1, Name: "Book 1.2", Children: []models.Notebook{
						{Id: 5, ParentId: 4, Name: "Book 1.2.1"},
					}},
				}},
				{Id: 3, ParentId: 0, Name: "Book 2", Children: []models.Notebook{
					{Id: 6, ParentId: 3, Name: "Book 2.1"},
				}},
			},
		},
		highestNotebookId: 6,
	}
	reader.notes = Notes{
		notes: []models.Note{
			{Id: 0, NotebookId: 0, Name: "MyGreatNote1"},
			{Id: 1, NotebookId: 0, Name: "MyGreatNote2"},
			{Id: 2, NotebookId: 1, Name: "MyGreatNote3"},
			{Id: 3, NotebookId: 3, Name: "MyGreatNote4"},
			{Id: 4, NotebookId: 4, Name: "MyGreatNote5"},
			{Id: 5, NotebookId: 5, Name: "subbooking"},
		},
		highestNoteId: 5,
	}
	go func() {
		for {
			select {
			case id := <-OpenNoteChan:
				reader.requestedNote = id
			case id := <-OpenNotebooksChan:
				reader.requestedBook = id
				booksNotes := notesByNotebookId(id, reader.notes.notes)
				if len(booksNotes) > 0 {
					reader.requestedNote = booksNotes[0].Id
				}
			case path := <-makeBookChan:
				err := writer.makeBook(reader, path)
				makeBookErrorChan <- err
			case name := <-makeNoteChan:
				err := writer.makeNote(reader, name)
				makeNoteErrorChan <- err
			case <-time.After(250 * time.Millisecond): // TODO get rid of magic number for refresh interval
				if reader.openBook == -1 {
					reader.OpenBooks()
				}
				if reader.requestedBook != reader.openBook {
					reader.OpenBook(reader.requestedBook)
				}
				if reader.requestedNote != reader.openNote {
					reader.OpenNote(reader.requestedNote)
				}
			}
		}
	}()
	return fs
}

func (b *DummyReader) OpenBooks() {
	NotebooksChan <- b.notebooks.notebookRoot
}

func (b *DummyReader) OpenBook(id int) {
	b.notes.openBookId = id
	books := notesByNotebookId(id, b.notes.notes)
	NotesChan <- books
	b.openBook = id
}

func (b *DummyReader) OpenNote(id int) {
	note := noteById(&b.notes.notes, id)
	note.Body = strings.NewReader(note.Name + dummyText)
	NoteChan <- *note
	b.openNote = id
}

func (b *DummySource) MakeBook(path string) error {
	//Offload to the source goroutine
	makeBookChan <- path
	return <-makeBookErrorChan
}

func (b *DummySource) MakeNote(name string) error {
	//Offload to the source goroutine
	makeNoteChan <- name
	return <-makeNoteErrorChan
}

func (b *DummyReader) getNotebooks() *Notebooks {
	return &b.notebooks
}

func (b *DummyReader) getNotes() *Notes {
	return &b.notes
}
func (b *DummyReader) queueUpdate() {
	b.openBook = -1
	b.openNote = -1
}

func (b *DummyWriter) makeBook(reader NotebookReader, path string) error {
	notebooks := reader.getNotebooks()
	parent, err := parentByPath(path, &notebooks.notebookRoot)
	if err != nil {
		return err
	}
	_, dir := filepath.Split(path)
	for _, book := range parent.Children {
		if book.Name == dir {
			return errors.New("There is already a book at path " + path)
		}
	}
	parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: notebooks.highestNotebookId, ParentId: parent.Id, Path: path})
	notebooks.highestNotebookId++
	reader.queueUpdate()
	return nil
}

func (b *DummyWriter) makeNote(reader NotebookReader, name string) error {
	notes := reader.getNotes()
	notebooks := reader.getNotebooks()
	notebook := notebookById(notes.openBookId, &notebooks.notebookRoot)
	booksNotes := notesByNotebookId(notebook.Id, notes.notes)
	for _, note := range booksNotes {
		if note.Name == name+".md" {
			return errors.New("There is already a note named " + name)
		}
	}

	path := notebook.Path + "/" + name + ".md"
	notes.highestNoteId++
	note := models.Note{Name: name + ".md", Id: notes.highestNoteId, NotebookId: notebook.Id, Path: path}
	notes.notes = append(notes.notes, note)
	reader.queueUpdate()
	return nil
}
