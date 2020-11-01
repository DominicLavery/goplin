package data

import (
	"dominiclavery/goplin/models"
	"errors"
	"strings"
)

type Source interface {
	Notebooks(notebookCallback func(models.Notebook))
	Notes(noteCallback func([]models.Note))
	Note(noteCallback func(models.Note))
	OpenBook(id int)
	OpenNote(id int)
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

func notesByNotebookId(notes []models.Note, notebookId int) []models.Note {
	filtered := make([]models.Note, 0)
	for _, note := range notes {
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
