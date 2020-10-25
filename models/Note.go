package models

import "io"

type Note struct {
	Id         int
	NotebookId int
	Name       string
	Body       io.Reader
	Path       string
}

func ByNotebookId(notes []Note, notebookId int) []Note {
	filtered := make([]Note, 0)
	for _, note := range notes {
		if note.NotebookId == notebookId {
			filtered = append(filtered, note)
		}
	}
	return filtered
}
