package data

import (
	"github.com/spf13/afero"
	"strings"
	"testing"
)

func TestNewFilesystemSourceNotebooksLoading(t *testing.T) {
	//arrange
	tables := []struct {
		name          string
		fileStructure fileStructure
		expectedBooks *Notebook
	}{
		{
			name:          "flat notebook structure",
			fileStructure: flatFileStructure,
			expectedBooks: &Notebook{
				Name: "root",
				Path: "/",
				Children: []*Notebook{
					{
						Name:  "book1",
						Path:  "/book1",
						Notes: []*Note{{Name: "note1", Path: "/book1/note1.md"}},
					},
					{
						Name:  "book2",
						Path:  "/book2",
						Notes: []*Note{{Name: "note2", Path: "/book2/note2.md"}},
					},
					{
						Name:  "book3",
						Path:  "/book3",
						Notes: []*Note{{Name: "note3", Path: "/book3/note3.md"}},
					},
				},
			},
		},
		{
			name:          "Nested notebook structure",
			fileStructure: nestedFileStructure,
			expectedBooks: &Notebook{
				Name: "root",
				Path: "/",
				Children: []*Notebook{{
					Name:  "book1",
					Path:  "/book1",
					Notes: []*Note{{Name: "note1", Path: "/book1/note1.md"}},
					Children: []*Notebook{{
						Name:  "book2",
						Path:  "/book1/book2",
						Notes: []*Note{{Name: "note2", Path: "/book1/book2/note2.md"}},
						Children: []*Notebook{{
							Name:  "book3",
							Path:  "/book1/book2/book3",
							Notes: []*Note{{Name: "note3", Path: "/book1/book2/book3/note3.md"}},
						}},
					}},
				}},
			},
		},
		{
			name: "ignore dirs starting with a .",
			fileStructure: fileStructure{
				name: "root",
				children: []fileStructure{
					{
						name: "book1",
					},
					{
						name: ".ignored",
					},
				},
			},
			expectedBooks: &Notebook{
				Name: "root",
				Path: "/",
				Children: []*Notebook{
					{
						Name: "book1",
						Path: "/book1",
					},
				},
			},
		},
	}

	for _, test := range tables {
		t.Run(test.name, func(t *testing.T) {

			afs := afero.NewMemMapFs()
			makeDirStructure(test.fileStructure, "/", afs)

			//act
			fs := newFilesystemSource("/", &afero.Afero{Fs: afero.NewBasePathFs(afs, "/root")})
			actual := fs.openBooks()
			//assert
			assertBookIsLike(test.expectedBooks, &actual, t)
		})
	}
}

func TestFilesystemWriter_MakeBook(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)
	fs := newFilesystemSource("/", &afero.Afero{Fs: afero.NewBasePathFs(afs, "/root")})

	expectedBook := &Notebook{
		Name:     "new",
		Path:     "/new",
		Children: nil,
	}

	//act
	book, _ := fs.makeBook(expectedBook.Path)

	//assert
	assertBookIsLike(book, expectedBook, t)
	_, err := afs.Stat("/root/" + expectedBook.Path)
	if err != nil {
		t.Errorf("Expected a directory to be created at %s", expectedBook.Path)
	}
}

func TestFilesystemWriter_MakeNote(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)
	fs := newFilesystemSource("/", &afero.Afero{Fs: afero.NewBasePathFs(afs, "/root")})

	expectedNote := Note{Name: "new", Path: "/new.md"}

	//act
	actual, _ := fs.makeNote("/new")
	//assert
	assertNotesAreLike(&expectedNote, actual, t)
	_, err := afs.Stat("/root/" + expectedNote.Path)
	if err != nil {
		t.Errorf("Expected a file to be created at %s", expectedNote.Path)
	}
}

func assertNotesAreLike(expected *Note, actual *Note, t *testing.T) {
	if expected.Name != actual.Name {
		t.Errorf("Expected note to have name [%s]. Had [%s]", expected.Name, actual.Name)
	}
	if expected.Path != actual.Path {
		t.Errorf("Expected note to have path [%s]. Had [%s]", expected.Path, actual.Path)
	}
}

func assertBookIsLike(expected *Notebook, actual *Notebook, t *testing.T) {
	if expected.Name != actual.Name {
		t.Errorf("Expected book to have name [%s]. Had [%s]", expected.Name, actual.Name)
	}
	if expected.Id != actual.Id {
		t.Errorf("Expected book to have ID [%s]. Had [%s]", expected.Id, actual.Id)
	}
	if expected.Path != actual.Path {
		t.Errorf("Expected book to have Path [%s]. Had [%s]", expected.Path, actual.Path)
	}
	if expected.Children != nil {
		if actual.Children == nil {
			t.Errorf("Expected book to have [%d] childrem. Had 0", len(expected.Children))
		}
		if len(expected.Children) != len(actual.Children) {
			t.Errorf("Expected book to have [%d] childrem. Had [%d]", len(expected.Children), len(actual.Children))
		}
		for i, book := range expected.Children {
			assertBookIsLike(book, actual.Children[i], t)
		}
	}
	if expected.Notes != nil {
		for i, note := range expected.Notes {
			assertNotesAreLike(note, actual.Notes[i], t)
		}
	}
}

type fileStructure struct {
	name     string
	children []fileStructure
}

func makeDirStructure(notebook fileStructure, path string, fs afero.Fs) {
	createdPath := path + notebook.name
	if strings.HasSuffix(notebook.name, ".md") {
		_, _ = fs.Create(createdPath)
	} else {
		_ = fs.Mkdir(createdPath, 0644)
		for _, child := range notebook.children {
			makeDirStructure(child, createdPath+"/", fs)
		}
	}
}

var flatFileStructure = fileStructure{
	name: "root",
	children: []fileStructure{
		{
			name: "book1",
			children: []fileStructure{
				{name: "note1.md"},
			},
		},
		{
			name: "book2",
			children: []fileStructure{
				{name: "note2.md"},
			},
		},
		{
			name: "book3",
			children: []fileStructure{
				{name: "note3.md"},
			},
		},
	},
}

var nestedFileStructure = fileStructure{
	name: "root",
	children: []fileStructure{
		{
			name: "book1",
			children: []fileStructure{
				{name: "note1.md"},
				{
					name: "book2",
					children: []fileStructure{
						{name: "note2.md"},
						{
							name: "book3",
							children: []fileStructure{
								{name: "note3.md"},
							},
						},
					},
				},
			},
		},
	},
}
