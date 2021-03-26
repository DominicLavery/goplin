package data

import (
	"github.com/google/uuid"
	"io"
)

type Source interface {
	openBooks() Notebook
	openNote(path string) io.Reader
	makeBook(path string) (*Notebook, error)
	makeNote(path string) (*Note, error)
}

type SourceDataTree struct {
	sourceId     uuid.UUID
	sourceName   string
	NotebookRoot *Notebook
}

type Notebook struct {
	Id       uuid.UUID
	Name     string
	Path     string
	Children []*Notebook
	Notes    []*Note
}

type Note struct {
	Id   uuid.UUID
	Name string
	Path string
}
