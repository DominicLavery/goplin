package data

import (
	"dominiclavery/goplin/models"
	"errors"
	"strings"
	"sync"
)

var notebooks = Notebooks{
	notebookRoot: models.Notebook{Id: 0, ParentId: -1},
	mu:           sync.Mutex{},
}

var notes = Notes{
	notes: nil,
	mu:    sync.Mutex{},
}

var NotebooksChan = make(chan models.Notebook)
var NotesChan = make(chan []models.Note)
var NoteChan = make(chan models.Note)

var OpenNoteChan = make(chan int)
var OpenNotebooksChan = make(chan int)

type Source interface {
	NotebookReader
	NotebookWriter
}

type NotebookReader interface {
	OpenBook(id int)
	OpenNote(id int)
	OpenBooks()
}

type NotebookWriter interface {
	MakeBook(path string) error
	MakeNote(name string) error
}

func parentByPath(path string, notebooks *models.Notebook) (*models.Notebook, error) {
	pathParts := strings.Split(path, "/")
	pathParts = pathParts[:len(pathParts)-1] // remove the name
	parent := notebooks
	for _, part := range pathParts {
		parent = notebookByName(part, &parent.Children)
		if parent == nil {
			return nil, errors.New(part + "not found")
		}
	}
	return parent, nil
}

func notebookByName(name string, notebooks *[]models.Notebook) *models.Notebook {
	for i, book := range *notebooks {
		if name == book.Name {
			return &(*notebooks)[i] // Book is a copy, get the pointer. //TODO Has to be something better for this
		}
	}
	return nil
}

func notebookById(id int, notebooks *models.Notebook) *models.Notebook {
	var found *models.Notebook
	if id == notebooks.Id {
		return notebooks
	}
	for _, book := range notebooks.Children {
		found = notebookById(id, &book)
		if found != nil {
			return found
		}
	}

	return found
}

func notesByNotebookId(notebookId int) []models.Note {
	filtered := make([]models.Note, 0)
	for _, note := range notes.notes {
		if note.NotebookId == notebookId {
			filtered = append(filtered, note)
		}
	}
	return filtered
}

func noteById(notes *[]models.Note, id int) *models.Note {
	for _, note := range *notes {
		if note.Id == id {
			return &note
		}
	}
	return nil
}
