package data

import (
	"github.com/google/uuid"
	"io"
	"sync"
)

var sources = make(map[string]Source)
var notebookTree = make(map[string]SourceDataTree)
var notebooks = make(map[uuid.UUID]*Notebook)
var bookToSource = make(map[uuid.UUID]string)
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
	addNotebook(sourceName, root)
}

func addNotebook(sourceName string, book *Notebook) {
	book.Id = uuid.New()
	bookToSource[book.Id] = sourceName
	notebooks[book.Id] = book
	for _, note := range book.Notes {
		addNote(note)
	}
	for _, child := range book.Children {
		addNotebook(sourceName, child)
	}
}

func addNote(note *Note) {
	note.Id = uuid.New()
	notes[note.Id] = note
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
	parent := notebooks[openBook]
	source := bookToSource[parent.Id]
	return sources[source].openNote(notes[id].Path)
}

func MakeBook(name string) error {
	parent := notebooks[openBook]
	source := bookToSource[parent.Id]
	book, err := sources[source].makeBook(parent.Path + "/" + name)
	if err != nil {
		return err
	}
	addNotebook(source, book)
	parent.Children = append(parent.Children, book)
	return nil
}

func MakeNote(name string) error {
	parent := notebooks[openBook]
	source := bookToSource[parent.Id]
	note, err := sources[source].makeNote(parent.Path + "/" + name)
	if err != nil {
		return err
	}
	addNote(note)
	parent.Notes = append(parent.Notes, note)
	return nil
}
