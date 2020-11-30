package data

import (
	"dominiclavery/goplin/logs"
	"dominiclavery/goplin/models"
	"errors"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemSource struct {
	NotebookReader
	NotebookWriter
}

type FilesystemReader struct {
	rootPath string
	fs       *afero.Afero
}

type FilesystemWriter struct {
	rootPath string
	fs       *afero.Afero
}

func NewFilesystemSource(root string) *FilesystemSource {
	return newFilesystemSource(root, afero.NewOsFs())
}

func newFilesystemSource(root string, fs afero.Fs) *FilesystemSource {
	afs := &afero.Afero{Fs: fs}
	reader := &FilesystemReader{
		rootPath: root,
		fs:       afs,
	}
	if err := afs.Walk(root, reader.walkFn); err != nil {
		logs.TeeLog("Could not read notebooks", err)
		notebooks.notebookRoot.Name = "Error"
	}

	filesource := &FilesystemSource{
		NotebookReader: reader,
		NotebookWriter: &FilesystemWriter{
			rootPath: root,
			fs:       afs,
		},
	}

	go func() {
		for {
			select {
			case id := <-OpenNoteChan:
				filesource.OpenNote(id)
			case id := <-OpenNotebooksChan:
				if id == 0 {
					filesource.OpenBooks()
				} else {
					filesource.OpenBook(id)
				}
			}
		}
	}()
	return filesource
}

func (b *FilesystemReader) walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() && info.Name()[0:1] == "." {
		return filepath.SkipDir
	}

	relPath, _ := filepath.Rel(b.rootPath, path)
	parent, _ := parentByPath(relPath, &notebooks.notebookRoot)
	if path == b.rootPath {
		notebooks.notebookRoot.Name = info.Name()
		notebooks.notebookRoot.Path = path
		notebooks.highestNotebookId++
	} else if info.IsDir() {
		notebook := models.Notebook{Name: info.Name(), Id: notebooks.highestNotebookId, ParentId: parent.Id, Path: path}
		parent.Children = append(parent.Children, notebook)
		notebooks.highestNotebookId++
	} else if strings.HasSuffix(info.Name(), ".md") {
		notes.notes = append(notes.notes, models.Note{Name: info.Name(), Id: notes.highestNoteId, NotebookId: parent.Id, Path: path})
		notes.highestNoteId++
	}
	return nil
}

func (b *FilesystemWriter) MakeBook(path string) error {
	absPath, _ := filepath.Abs(path)
	parent, err := parentByPath(path, &notebooks.notebookRoot)
	if err != nil {
		return err
	}
	if err := b.fs.Mkdir(absPath, os.ModePerm); err != nil {
		return err
	}
	_, dir := filepath.Split(path)
	notebooks.highestNotebookId++
	parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: notebooks.highestNotebookId, ParentId: parent.Id, Path: absPath}) // TODO make relPath
	NotebooksChan <- notebooks.notebookRoot
	return nil
}

func (b *FilesystemWriter) MakeNote(name string) error {
	notebook := notebookById(notes.openBookId, &notebooks.notebookRoot)
	booksNotes := notesByNotebookId(notebook.Id)
	for _, note := range booksNotes {
		if note.Name == name+".md" {
			return errors.New("There is already a note named " + name)
		}
	}

	path := notebook.Path + "/" + name + ".md"
	file, err := b.fs.Create(path)
	if err != nil {
		return err
	}
	notes.highestNoteId++
	note := models.Note{Name: name + ".md", Id: notes.highestNoteId, NotebookId: notebook.Id, Path: path}
	notes.notes = append(notes.notes, note)
	_ = file.Close()
	NotesChan <- append(booksNotes, note)
	return nil
}

func (b *FilesystemReader) OpenNote(id int) {
	if notes.openNote != nil {
		if closer, ok := notes.openNote.Body.(io.Closer); ok {
			_ = closer.Close()
		}
	}

	note := noteById(&notes.notes, id)
	var file afero.File
	var err error
	if file, err = b.fs.Open(note.Path); err != nil {
		logs.TeeLog("Couldn't open the note", err)
		note.Body = strings.NewReader("Error!")
	} else {
		note.Body = file
	}
	NoteChan <- *note
}

func (b *FilesystemReader) OpenBook(id int) {
	notes.openBookId = id
	books := notesByNotebookId(id)
	NotesChan <- books
	if len(books) > 0 {
		b.OpenNote(books[0].Id)
	}
}
func (b *FilesystemReader) OpenBooks() {
	NotebooksChan <- notebooks.notebookRoot
	b.OpenBook(notebooks.notebookRoot.Id)
}
