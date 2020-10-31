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
}

func parentByPath(path string, notebooks *models.Notebook) (*models.Notebook, error) {
	pathParts := strings.Split(path, "/")
	pathParts = pathParts[:len(pathParts)-1] // remove the name
	parent := notebooks
	for _, part := range pathParts {
		parent = byName(part, &parent.Children)
		if parent == nil {
			return nil, errors.New(part + "not found")
		}
	}
	return parent, nil
}

func byName(name string, notebooks *[]models.Notebook) *models.Notebook {
	var found *models.Notebook
	for i, book := range *notebooks {
		if name == book.Name {
			found = &(*notebooks)[i] // Book is a copy, get the pointer. //TODO Has to be something better for this
			break
		}
	}
	return found
}
