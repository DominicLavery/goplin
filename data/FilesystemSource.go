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
	"time"
)

type FilesystemSource struct {
	NotebookReader
	NotebookWriter
}

type FilesystemReader struct {
	fs            *afero.Afero
	requestedBook int
	requestedNote int
	openNote      int
	openBook      int
	notebooks     Notebooks
	notes         Notes
}

type FilesystemWriter struct {
	fs       *afero.Afero
}

func NewFilesystemSource(root string) *FilesystemSource {
	fs := afero.NewBasePathFs(afero.NewOsFs(), root)
	return newFilesystemSource("/", fs)
}

func newFilesystemSource(root string, fs afero.Fs) *FilesystemSource {
	afs := &afero.Afero{Fs: fs}
	reader := &FilesystemReader{
		fs:            afs,
		requestedBook: 0,
		requestedNote: 0,
		openNote:      -1,
		openBook:      -1,
	}
	writer := &FilesystemWriter{
		fs:       afs,
	}
	if err := afs.Walk(root, reader.walkFn); err != nil {
		logs.TeeLog("Could not read notebooks", err)
		reader.notebooks.notebookRoot.Name = "Error"
	}

	filesource := &FilesystemSource{
		NotebookReader: reader,
		NotebookWriter: writer,
	}

	go func() {
		for {
			select {
			case id := <-OpenNoteChan:
				reader.requestedNote = id
			case id := <-OpenNotebooksChan:
				reader.requestedBook = id
				booksNotes := notesByNotebookId(id, reader.notes.notes)
				if len(booksNotes) > 0 {
					reader.requestedNote = booksNotes[0].Id
				}
			case path := <-makeBookChan:
				err := writer.makeBook(reader, path)
				makeBookErrorChan <- err
			case name := <-makeNoteChan:
				err := writer.makeNote(reader, name)
				makeNoteErrorChan <- err
			case <-time.After(250 * time.Millisecond): // TODO get rid of magic number for refresh interval
				if reader.openBook == -1 {
					reader.OpenBooks()
				}
				if reader.requestedBook != reader.openBook {
					reader.OpenBook(reader.requestedBook)
				}
				if reader.requestedNote != reader.openNote {
					reader.OpenNote(reader.requestedNote)
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

	parent, _ := parentByPath(path, &b.notebooks.notebookRoot)
	if path == "/" {
		b.notebooks.notebookRoot.Name = info.Name()
		b.notebooks.notebookRoot.Id = RootId
		b.notebooks.notebookRoot.Path = path
		b.notebooks.notebookRoot.ParentId = -1
		b.notebooks.highestNotebookId++
	} else if info.IsDir() {
		notebook := models.Notebook{Name: info.Name(), Id: b.notebooks.highestNotebookId, ParentId: parent.Id, Path: path}
		parent.Children = append(parent.Children, notebook)
		b.notebooks.highestNotebookId++
	} else if strings.HasSuffix(info.Name(), ".md") {
		b.notes.notes = append(b.notes.notes, models.Note{Name: info.Name(), Id: b.notes.highestNoteId, NotebookId: parent.Id, Path: path})
		b.notes.highestNoteId++
	}
	return nil
}

func (b *FilesystemSource) MakeBook(path string) error {
	//Offload to the source goroutine
	makeBookChan <- path
	return <-makeBookErrorChan
}

func (b *FilesystemSource) MakeNote(name string) error {
	//Offload to the source goroutine
	makeNoteChan <- name
	return <-makeNoteErrorChan
}

func (b *FilesystemReader) getNotebooks() *Notebooks {
	return &b.notebooks
}

func (b *FilesystemReader) getNotes() *Notes {
	return &b.notes
}

func (b *FilesystemReader) getOpenBookId() int {
	return b.openBook
}

func (b *FilesystemReader) queueUpdate() {
	b.openBook = -1
	b.openNote = -1
}

func (b *FilesystemWriter) makeBook(reader NotebookReader, name string) error {
	notebooks := reader.getNotebooks()
	parent := notebookById(reader.getOpenBookId(), &notebooks.notebookRoot)
	path := parent.Path + "/" + name
	if err := b.fs.Mkdir(path, os.ModePerm); err != nil {
		return err
	}
	notebooks.highestNotebookId++
	parent.Children = append(parent.Children, models.Notebook{Name: name, Id: notebooks.highestNotebookId, ParentId: parent.Id, Path: path})
	reader.queueUpdate()
	return nil
}

func (b *FilesystemWriter) makeNote(reader NotebookReader, name string) error {
	notebooks := reader.getNotebooks()
	notes := reader.getNotes()
	notebook := notebookById(reader.getOpenBookId(), &notebooks.notebookRoot)
	booksNotes := notesByNotebookId(notebook.Id, notes.notes)
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
	reader.queueUpdate()
	return nil
}

func (b *FilesystemReader) OpenNote(id int) {
	if b.notes.openNote != nil {
		if closer, ok := b.notes.openNote.Body.(io.Closer); ok {
			_ = closer.Close()
		}
	}

	note := noteById(&b.notes.notes, id)
	var file afero.File
	var err error
	if file, err = b.fs.Open(note.Path); err != nil {
		logs.TeeLog("Couldn't open the note", err)
		note.Body = strings.NewReader("Error!")
	} else {
		note.Body = file
	}
	NoteChan <- *note
	b.openNote = id
}

func (b *FilesystemReader) OpenBook(id int) {
	books := notesByNotebookId(id, b.notes.notes)
	NotesChan <- books
	b.openBook = id
}

func (b *FilesystemReader) OpenBooks() {
	NotebooksChan <- b.notebooks.notebookRoot
}
