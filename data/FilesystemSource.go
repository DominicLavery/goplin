package data

import (
	"dominiclavery/goplin/logs"
	"dominiclavery/goplin/models"
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
}

type FilesystemWriter struct {
	rootPath string
}

func NewFilesystemSource(root string) *FilesystemSource {
	reader := &FilesystemReader{
		rootPath: root,
	}
	if err := filepath.Walk(root, reader.walkFn); err != nil {
		logs.TeeLog("Could not read notebooks", err)
		notebooks.notebookRoot.Name = "Error"
	}

	fs := &FilesystemSource{
		NotebookReader: reader,
		NotebookWriter: &FilesystemWriter{
			rootPath: root,
		},
	}

	go func() {
		for {
			select {
			case id := <-OpenNoteChan:
				fs.OpenNote(id)
			case id := <-OpenNotebooksChan:
				if id == 0 {
					fs.OpenBooks()
				} else {
					fs.OpenBook(id)
				}
			}
		}
	}()
	return fs
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

func (b *FilesystemWriter) MakeBook(path string) (models.Notebook, error) {
	var notebook models.Notebook
	//absPath, _ := filepath.Abs(path)
	//parent, err := parentByPath(path, &b.notebooks.notebookRoot)
	//if err != nil {
	//	return notebook, err
	//}
	//if err := os.Mkdir(absPath, os.ModePerm); err != nil {
	//	return notebook, err
	//}
	//_, dir := filepath.Split(path)
	//parent.Children = append(parent.Children, models.Notebook{Name: dir, Id: b.notebooks.highestNotebookId, ParentId: parent.Id, Path: absPath})
	//b.notebooks.highestNotebookId++
	//b.notebooksUpdateHandler(b.notebooks.notebookRoot)
	return notebook, nil
}

func (b *FilesystemWriter) MakeNote(name string) (models.Note, error) {
	var note models.Note
	//notebook := notebookById(b.notes.openBookId, &b.notebooks.notebookRoot)
	//notes := notesByNotebookId(b.notes.notes, notebook.Id)
	//for _, note := range notes {
	//	if note.Name == name+".md" {
	//		return errors.New("There is already a note named " + name)
	//	}
	//}
	//
	//path := notebook.Path + "/" + name + ".md"
	//file, err := os.Create(path)
	//if err != nil {
	//	return err
	//}
	//b.notes.highestNoteId++
	//b.notes.notes = append(b.notes.notes, models.Note{Name: name + ".md", Id: b.notes.highestNoteId, NotebookId: notebook.Id, Path: path})
	//_ = file.Close()
	//b.OpenBook(notebook.Id)
	return note, nil
}

func (b *FilesystemReader) OpenNote(id int) {
	if notes.openNote != nil {
		if closer, ok := notes.openNote.Body.(io.Closer); ok {
			_ = closer.Close()
		}
	}

	note := noteById(&notes.notes, id)
	var file *os.File
	var err error
	if file, err = os.Open(note.Path); err != nil {
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
