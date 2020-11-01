package data

import (
	"dominiclavery/goplin/logs"
	"dominiclavery/goplin/models"
	"errors"
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
	openBookId             int
}

func NewFilesystemSource(root string) *FilesystemSource {
	var notes []models.Note
	notebooks := models.Notebook{Id: 0, ParentId: -1}

	notebookCount := 0
	noteCount := 0
	walkFunction := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name()[0:1] == "." {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(root, path)
		parent, _ := parentByPath(relPath, &notebooks)
		if path == root {
			notebooks.Name = info.Name()
			notebooks.Path = path
		} else if info.IsDir() {
			notebook := models.Notebook{Name: info.Name(), Id: notebookCount, ParentId: parent.Id, Path: path}
			parent.Children = append(parent.Children, notebook)
			notebookCount++
		} else if strings.HasSuffix(info.Name(), ".md") {
			notes = append(notes, models.Note{Name: info.Name(), Id: noteCount, NotebookId: parent.Id, Path: path})
			noteCount++
		}
		return nil
	}
	if err := filepath.Walk(root, walkFunction); err != nil {
		logs.TeeLog("Could not read notebooks", err)
		notebooks.Name = "Error"
	}
	return &FilesystemSource{rootPath: root, notebooks: notebooks, notes: notes, highestNotebookId: notebookCount, highestNoteId: noteCount}
}

func (b *FilesystemSource) MakeBook(path string) error {
	absPath, _ := filepath.Abs(path)
	parent, err := parentByPath(path, &b.notebooks)
	if err != nil {
		return err
	}
	if err := os.Mkdir(absPath, os.ModePerm); err != nil {
		return err
	}
	_, dir := filepath.Split(path)
	parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: b.highestNotebookId, ParentId: parent.Id})
	b.highestNotebookId++
	b.notebooksUpdateHandler(b.notebooks)
	return nil
}

func (b *FilesystemSource) MakeNote(name string) error {
	notebook := notebookById(b.openBookId, &b.notebooks)
	notes := notesByNotebookId(b.notes, notebook.Id)
	for _, note := range notes {
		if note.Name == name+".md" {
			return errors.New("There is already a note named " + name)
		}
	}

	path := notebook.Path + "/" + name + ".md"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	b.highestNoteId++
	b.notes = append(b.notes, models.Note{Name: name + ".md", Id: b.highestNoteId, NotebookId: notebook.Id, Path: path})
	_ = file.Close()
	b.OpenBook(notebook.Id)
	return nil
}

func (b *FilesystemSource) OpenNote(id int) {
	if b.noteUpdateHandler != nil {
		note := noteById(&b.notes, id)
		var file *os.File
		var err error
		if file, err = os.Open(note.Path); err != nil {
			logs.TeeLog("Couldn't open the note", err)
			note.Body = strings.NewReader("Error!")
		} else {
			note.Body = file
			defer file.Close()
		}
		b.noteUpdateHandler(*note)
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
	notes := notesByNotebookId(b.notes, b.openBookId)
	if len(notes) > 0 {
		b.OpenNote(notes[0].Id)
	}
}

func (b *FilesystemSource) OpenBook(id int) {
	b.openBookId = id
	if b.notesUpdateHandler != nil {
		b.notesUpdateHandler(notesByNotebookId(b.notes, id))
	}
}
