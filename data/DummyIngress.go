package data

import (
	"dominiclavery/goplin/models"
	"errors"
	"fmt"
	"strings"
)

type DummySource struct {
	notebooks              models.Notebook
	notes                  []models.Note
	notebooksUpdateHandler func(models.Notebook)
	notesUpdateHandler     func([]models.Note)
	noteUpdateHandler      func(models.Note)
	highestNotebookId      int
	highestNoteId          int
	openBookId             int
}

const dummyText = `

# Title Level 1
## Title Level 2
|test|test2|
|---|---|
|Details|more details|
`

func NewDummySource() *DummySource {
	return &DummySource{
		notebooks: models.Notebook{
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
		notes: []models.Note{
			{Id: 0, NotebookId: 0, Name: "MyGreatNote1"},
			{Id: 1, NotebookId: 0, Name: "MyGreatNote2"},
			{Id: 2, NotebookId: 1, Name: "MyGreatNote3"},
			{Id: 3, NotebookId: 3, Name: "MyGreatNote4"},
			{Id: 4, NotebookId: 4, Name: "MyGreatNote5"},
			{Id: 5, NotebookId: 5, Name: "subbooking"},
		},
		highestNotebookId: 6,
		highestNoteId:     5,
	}
}

func (b *DummySource) Notebooks(notebookCallback func(models.Notebook)) {
	b.notebooksUpdateHandler = notebookCallback
	b.notebooksUpdateHandler(b.notebooks)
}

func (b *DummySource) Notes(noteCallback func([]models.Note)) {
	b.notesUpdateHandler = noteCallback
	b.OpenBook(0)
}

func (b *DummySource) Note(noteCallback func(models.Note)) {
	b.noteUpdateHandler = noteCallback
	b.OpenNote(0)
}

func (b *DummySource) OpenBook(id int) {
	if b.notesUpdateHandler != nil {
		b.notesUpdateHandler(notesByNotebookId(b.notes, id))
	}
	b.openBookId = id
}

func (b *DummySource) OpenNote(id int) {
	if b.noteUpdateHandler != nil {
		for _, note := range b.notes {
			if note.Id == id {
				note.Body = strings.NewReader(note.Name + dummyText)
				b.noteUpdateHandler(note)
			}
		}
	}
}

func (b *DummySource) MakeBook(path string) error {
	b.highestNotebookId++
	var parent *models.Notebook
	var notebook models.Notebook

	if !strings.Contains(path, "/") {
		parent = &b.notebooks
		notebook = models.Notebook{Id: b.highestNotebookId, ParentId: parent.Id, Name: path}
	} else {
		var err error
		if parent, err = parentByPath(path, &b.notebooks); err != nil {
			return err
		}
		pathParts := strings.Split(path, "/")
		notebook = models.Notebook{Id: b.highestNotebookId, ParentId: parent.Id, Name: pathParts[len(pathParts)-1]}
	}

	if notebookByName(notebook.Name, &parent.Children) != nil {
		return errors.New("A notebook by this name already exists")
	}

	parent.Children = append(parent.Children, notebook)
	b.notebooksUpdateHandler(b.notebooks)
	return nil
}

func (b *DummySource) MakeNote(name string) error {
	notes := notesByNotebookId(b.notes, b.openBookId)
	for _, note := range notes {
		if note.Name == name {
			return errors.New(fmt.Sprintf("There is already a book called %s in this notebook", name))
		}
	}
	b.highestNoteId++
	newNote := models.Note{Id: b.highestNoteId, NotebookId: b.openBookId, Name: name}
	b.notes = append(b.notes, newNote)
	b.notesUpdateHandler(notesByNotebookId(b.notes, b.openBookId))
	return nil
}
