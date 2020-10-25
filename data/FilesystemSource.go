package data

import (
	"dominiclavery/goplin/models"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemSource struct {
	notebooks              models.Notebook
	notes                  []models.Note
	notebooksUpdateHandler func(models.Notebook)
	notesUpdateHandler     func([]models.Note)
	noteUpdateHandler      func(models.Note)
	highestNotebookId      int
	highestNoteId          int
	rootPath               string
	currentFile            *os.File
}

func NewFilesystemSource(root string) *FilesystemSource {
	var notes []models.Note
	notebooks := models.Notebook{Id: 0, ParentId: -1}

	notebookCount := 0
	noteCount := 0
	walkFunction := func(path string, info os.FileInfo, err error) error {
		if err != nil { //TODO handle me?
			return err
		}

		if info.IsDir() && info.Name()[0:1] == "." {
			return filepath.SkipDir
		} else if info.IsDir() {
			if path == root {
				notebooks.Name = info.Name()
			} else {
				relPath, _ := filepath.Rel(root, path)
				parent := parentByPath(relPath, &notebooks)
				parent.Children = append(parent.Children, models.Notebook{Name: info.Name(), Id: notebookCount, ParentId: parent.Id})
			}
			notebookCount++
		} else if strings.HasSuffix(info.Name(), ".md") {
			relPath, _ := filepath.Rel(root, path)
			parent := parentByPath(relPath, &notebooks)
			notes = append(notes, models.Note{Name: info.Name(), Id: noteCount, NotebookId: parent.Id, Path: path})
			noteCount++
		}
		return nil
	}
	filepath.Walk(root, walkFunction) //TODO handle error
	return &FilesystemSource{rootPath: root, notebooks: notebooks, notes: notes, highestNotebookId: notebookCount, highestNoteId: noteCount}
}

func (b *FilesystemSource) MakeBook(path string) {
	absPath, _ := filepath.Abs(path)
	parent := parentByPath(path, &b.notebooks)
	_ = os.Mkdir(absPath, os.ModePerm) //TODO
	_, dir := filepath.Split(path)
	parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: b.highestNotebookId, ParentId: parent.Id})
	b.highestNotebookId++
	b.notebooksUpdateHandler(b.notebooks)
}

func (b *FilesystemSource) OpenNote(id int) {
	if b.currentFile != nil {
		if err := b.currentFile.Close(); err != nil {
			log.Fatal("We couldn't close the read file!") // TODO more user friendly handling
		}
		b.currentFile = nil
	}
	if b.noteUpdateHandler != nil {
		for _, note := range b.notes {
			if note.Id == id {
				var file *os.File
				var err error
				if file, err = os.Open(note.Path); err != nil {
					note.Body = strings.NewReader("Couldn't open the file")
				} else {
					note.Body = file
					b.currentFile = file
				}
				b.noteUpdateHandler(note)
			}
		}
	}
}

func (b *FilesystemSource) Notebooks(notebookCallback func(models.Notebook)) {
	b.notebooksUpdateHandler = notebookCallback
	b.notebooksUpdateHandler(b.notebooks)
}

func (b *FilesystemSource) Notes(noteCallback func([]models.Note)) {
	b.notesUpdateHandler = noteCallback
	b.OpenBook(0)
}

func (b *FilesystemSource) Note(noteCallback func(models.Note)) {
	b.noteUpdateHandler = noteCallback
	b.OpenNote(0)
}

func (b *FilesystemSource) OpenBook(id int) {
	if b.notesUpdateHandler != nil {
		b.notesUpdateHandler(models.ByNotebookId(b.notes, id))
	}
}
