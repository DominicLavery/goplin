package data

import (
	"github.com/google/uuid"
	"io"
	"sync"
)

var sources = make(map[string]Source)
var notebookTree = make(map[string]SourceDataTree)
var notebooks = make(map[uuid.UUID]*Notebook)
var notes = make(map[uuid.UUID]*Note)
var openBook uuid.UUID
var m sync.Mutex

func RegisterSource(sourceName string, source Source) {
	m.Lock()
	defer m.Unlock()
	root := source.openBooks()
	setRoot(sourceName, &root)
	sources[sourceName] = source
}

func setRoot(sourceName string, root *Notebook) {
	notebookTree[sourceName] = SourceDataTree{
		sourceId:     uuid.New(),
		sourceName:   sourceName,
		NotebookRoot: root,
	}
	addNotebook(root)
}

func addNotebook(book *Notebook) {
	book.Id = uuid.New()
	notebooks[book.Id] = book
	for _, note := range book.Notes {
		note.Id = uuid.New()
		notes[note.Id] = note
	}
	for _, child := range book.Children {
		addNotebook(child)
	}
}

func GetBooks() map[string]SourceDataTree {
	m.Lock()
	defer m.Unlock()
	return notebookTree
}

func GetBook(id uuid.UUID) *Notebook {
	m.Lock()
	defer m.Unlock()
	openBook = id
	return notebooks[id]
}

func GetNote(id uuid.UUID) io.Reader {
	m.Lock()
	defer m.Unlock()
	return sources["Local"].openNote(notes[id].Path) //TODO
}

func MakeBook(name string) error {
	parent := notebooks[openBook]
	book, err := sources["Local"].makeBook(parent.Path + "/" + name) // TODO
	if err != nil {
		return err
	}
	book.Id = uuid.New()
	parent.Children = append(parent.Children, book)
	notebooks[book.Id] = book
	return nil
}

func MakeNote(name string) error {
	parent := notebooks[openBook]
	note, err := sources["Local"].makeNote(parent.Path + "/" + name) //TODO
	if err != nil {
		return err
	}
	note.Id = uuid.New()
	parent.Notes = append(parent.Notes, note)
	notes[note.Id] = note
	return nil
}
