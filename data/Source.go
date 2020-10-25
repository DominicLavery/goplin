package data

import (
	"dominiclavery/goplin/models"
	"strings"
)

type Source interface {
	Notebooks(notebookCallback func(models.Notebook))
	Notes(noteCallback func([]models.Note))
	Note(noteCallback func(models.Note))
	OpenBook(id int)
	OpenNote(id int)
	MakeBook(path string)
}

func parentByPath(path string, notebooks *models.Notebook) *models.Notebook {
	pathParts := strings.Split(path, "/")
	pathParts = pathParts[:len(pathParts)-1] // remove the name
	parent := notebooks
	for _, part := range pathParts {
		parent = byName(part, &parent.Children)
		if parent == nil {
			panic(part + " not found") //todo not found
		}
	}
	return parent
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
