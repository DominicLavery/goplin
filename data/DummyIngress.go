package data

import (
	"dominiclavery/goplin/models"
	"errors"
	"path/filepath"
	"strings"
)

type DummySource struct {
	NotebookReader
	NotebookWriter
}

type DummyReader struct {
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
	fs := &DummySource{
		NotebookReader: &DummyReader{},
		NotebookWriter: &DummyWriter{},
	}

	notebooks = Notebooks{
		notebookRoot: models.Notebook{
			Id: 0, ParentId: -1, Name: "rootBook", Children: []models.Notebook{
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
	notes = Notes{
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
				fs.OpenNote(id)
			case id := <-OpenNotebooksChan:
				if id == 0 {
					fs.OpenBooks()
				} else {
					fs.OpenBook(id)
				}
			}
		}
	}()
	return fs
}

func (b *DummyReader) OpenBooks() {
	NotebooksChan <- notebooks.notebookRoot
	b.OpenBook(notebooks.notebookRoot.Id)
}

func (b *DummyReader) OpenBook(id int) {
	notes.openBookId = id
	books := notesByNotebookId(id)
	NotesChan <- books
	if len(books) > 0 {
		b.OpenNote(books[0].Id)
	}
}

func (b *DummyReader) OpenNote(id int) {
	note := noteById(&notes.notes, id)
	note.Body = strings.NewReader(note.Name + dummyText)
	NoteChan <- *note
}

func (b *DummyWriter) MakeBook(path string) error {
	parent, err := parentByPath(path, &notebooks.notebookRoot)
	if err != nil {
		return err
	}
	_, dir := filepath.Split(path)
	parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: notebooks.highestNotebookId, ParentId: parent.Id, Path: path})
	notebooks.highestNotebookId++
	NotebooksChan <- notebooks.notebookRoot
	return nil
}

func (b *DummyWriter) MakeNote(name string) error {

	notebook := notebookById(notes.openBookId, &notebooks.notebookRoot)
	booksNotes := notesByNotebookId(notebook.Id)
	for _, note := range booksNotes {
		if note.Name == name+".md" {
			return errors.New("There is already a note named " + name)
		}
	}

	path := notebook.Path + "/" + name + ".md"
	notes.highestNoteId++
	note := models.Note{Name: name + ".md", Id: notes.highestNoteId, NotebookId: notebook.Id, Path: path}
	notes.notes = append(notes.notes, note)
	NotesChan <- append(booksNotes, note)
	return nil
}
