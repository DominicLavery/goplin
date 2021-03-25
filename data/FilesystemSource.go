package data

import (
	"dominiclavery/goplin/logs"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemSource struct {
	fs    *afero.Afero
	root  string
	books map[string]*Notebook
	notes map[string]*Note
}

func NewFilesystemSource(root string) *FilesystemSource {
	fs := afero.NewBasePathFs(afero.NewOsFs(), root)
	return newFilesystemSource("/", fs)
}
func NewInMemorySource() *FilesystemSource {
	fs := afero.NewMemMapFs()
	return newFilesystemSource("/", fs)
}

func newFilesystemSource(root string, fs afero.Fs) *FilesystemSource {
	afs := &afero.Afero{Fs: fs}
	return &FilesystemSource{
		fs:    afs,
		root:  root,
		books: make(map[string]*Notebook),
		notes: make(map[string]*Note),
	}
}

func (b *FilesystemSource) openBooks() Notebook {
	root := Notebook{}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "/" {
			root.Name = info.Name()
			root.Path = path
			b.books[path] = &root
			return nil
		}
		
		if info.IsDir() && info.Name()[0:1] == "." {
			return filepath.SkipDir
		}
		 if info.IsDir() {
			parent := b.books[filepath.Dir(path)]
			notebook := &Notebook{Name: info.Name(), Path: path}
			parent.Children = append(parent.Children, notebook)
			b.books[path] = notebook
		} else if strings.HasSuffix(info.Name(), ".md") {
			parent := b.books[filepath.Dir(path)]
			note := &Note{Name: strings.TrimSuffix(info.Name(), ".md"), Path: path}
			parent.Notes = append(parent.Notes, note)
			b.notes[path] = note
		}
		return nil
	}
	if err := b.fs.Walk(b.root, walker); err != nil {
		logs.TeeLog("Could not read notebooks", err)
		root.Name = "Error" //TODO
	}

	return root
}

func (b *FilesystemSource) openNote(path string) io.Reader {
	note := b.notes[path]
	if file, err := b.fs.Open(note.Path); err != nil {
		logs.TeeLog("Couldn't open the note", err)
		return strings.NewReader("Error!")
	} else {
		return file
	}
}

func (b *FilesystemSource) makeBook(path string) (*Notebook, error) {
	if err := b.fs.Mkdir(path, os.ModePerm); err != nil {
		return nil, err
	}
	book := &Notebook{Name: filepath.Base(path), Path: path}
	b.books[path] = book
	return book, nil
}

func (b *FilesystemSource) makeNote(path string) (*Note, error) {
	fullPath := path + ".md"
	_, err := b.fs.Create(fullPath)
	if err != nil {
		return nil, err
	}
	note := &Note{Name: filepath.Base(path),  Path: fullPath}
	b.notes[fullPath] = note
	return note, nil
}