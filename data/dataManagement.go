package data

import (
	"github.com/google/uuid"
	"io"
	"path/filepath"
	"sync"
)

var sources = make(map[string]Source)
var notebookTree = make(map[string]SourceDataTree)
var notebooks = make(map[uuid.UUID]*Notebook)
var bookToSource = make(map[uuid.UUID]Source)
var notes = make(map[uuid.UUID]*Note)
var openBook uuid.UUID
var m sync.Mutex

func RegisterSource(sourceName string, source Source) {
	m.Lock()
	defer m.Unlock()
	root := source.openBooks()
	setRoot(sourceName, source, &root)
	sources[sourceName] = source
}

func setRoot(sourceName string, source Source, root *Notebook) {
	notebookTree[sourceName] = SourceDataTree{
		sourceId:     uuid.New(),
		sourceName:   sourceName,
		NotebookRoot: root,
	}
	addNotebook(source, root, notebookTree[sourceName].sourceId)
}

func addNotebook(source Source, book *Notebook, parentId uuid.UUID) {
	book.Id = uuid.New()
	bookToSource[book.Id] = source
	notebooks[book.Id] = book
	book.PId = parentId
	for _, note := range book.Notes {
		addNote(note)
	}
	for _, child := range book.Children {
		addNotebook(source, child, book.Id)
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
	return source.openNote(notes[id].Path)
}

func MakeBook(name string) error {
	parent := notebooks[openBook]
	source := bookToSource[parent.Id]
	book, err := source.makeBook(filepath.Join(parent.Path, name))
	if err != nil {
		return err
	}
	addNotebook(source, book, parent.Id)
	parent.Children = append(parent.Children, book)
	return nil
}

func MakeNote(name string) error {
	parent := notebooks[openBook]
	source := bookToSource[parent.Id]
	note, err := source.makeNote(filepath.Join(parent.Path, name))
	if err != nil {
		return err
	}
	addNote(note)
	parent.Notes = append(parent.Notes, note)
	return nil
}

func DeleteNotebook(id uuid.UUID) error {
	source := bookToSource[id]
	book := notebooks[id]
	err := source.deleteBook(book)
	parent := notebooks[book.PId]
	if err != nil {
		return err
	}
	delete(notebooks, id)
	for i, b := range parent.Children {
		if b.Id == id {
			parent.Children = remove(parent.Children, i)
		}
	}
	return nil
}
func remove(s []*Notebook, i int) []*Notebook {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
