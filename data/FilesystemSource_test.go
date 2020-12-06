package data

import (
	"dominiclavery/goplin/models"
	"github.com/spf13/afero"
	"strings"
	"testing"
)

func TestNewFilesystemSourceNotebooksLoading(t *testing.T) {
	//arrange
	tables := []struct {
		name          string
		fileStructure fileStructure
		expectedBooks models.Notebook
	}{
		{
			name:          "flat notebook structure",
			fileStructure: flatFileStructure,
			expectedBooks: models.Notebook{
				Name:     "root",
				Id:       0,
				ParentId: -1,
				Children: []models.Notebook{
					{
						Id:       1,
						ParentId: 0,
						Name:     "book1",
					},
					{
						Id:       2,
						ParentId: 0,
						Name:     "book2",
					},
					{
						Id:       3,
						ParentId: 0,
						Name:     "book3",
					},
				},
			},
		},
		{
			name:          "Nested notebook structure",
			fileStructure: nestedFileStructure,
			expectedBooks: models.Notebook{
				Name:     "root",
				Id:       0,
				ParentId: -1,
				Children: []models.Notebook{{
					Id:       1,
					ParentId: 0,
					Name:     "book1",
					Children: []models.Notebook{{
						Id:       2,
						ParentId: 1,
						Name:     "book2",
						Children: []models.Notebook{{
							Id:       3,
							ParentId: 2,
							Name:     "book3",
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
			expectedBooks: models.Notebook{
				Name:     "root",
				Id:       0,
				ParentId: -1,
				Children: []models.Notebook{
					{
						Id:       1,
						ParentId: 0,
						Name:     "book1",
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
			fs := newFilesystemSource("/root", &afero.Afero{Fs: afs})

			//assert
			assertBookIsLike(test.expectedBooks, fs.getNotebooks().notebookRoot, t)
		})
	}
}

func TestFilesystemWriter_MakeBook(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)

	fw := FilesystemWriter{
		rootPath: "/root",
		fs:       &afero.Afero{Fs: afs},
	}
	dr := DummyReader{notebooks: Notebooks{notebookRoot: models.Notebook{
		Name: "root",
		Path: "/root",
	}}}
	expectedBook := models.Notebook{
		Id:       1,
		ParentId: 0,
		Name:     "new",
		Children: nil,
	}

	//act
	_ = fw.makeBook(&dr, "new")

	//assert
	if len(dr.getNotebooks().notebookRoot.Children) != 1 {
		t.Errorf("Expected root book to have 1 children. Had %d children", len(dr.getNotebooks().notebookRoot.Children))
	}
	assertBookIsLike(expectedBook, dr.getNotebooks().notebookRoot.Children[0], t)
	//TODO test for dir existence once relative pathing is correctly set up (and therefore, the location of the dir is predictable)
}

func TestFilesystemWriter_MakeNote(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)

	fw := FilesystemWriter{
		rootPath: "/root",
		fs:       &afero.Afero{Fs: afs},
	}
	dr := DummyReader{notebooks: Notebooks{
		notebookRoot: models.Notebook{
			Name: "root",
			Path: "/root",
		}},
	}
	expectedNote := Notes{
		notes: []models.Note{
			{Id: 1, NotebookId: 0, Name: "new.md", Path: "/root/new.md"},
		},
	}

	//act
	_ = fw.makeNote(&dr, "new")

	//assert
	assertNotesAreLike(expectedNote, *dr.getNotes(), t)
	//TODO test for notes existence once relative pathing is correctly set up (and therefore, the location of the dir is predictable)
}

func TestNewFilesystemSourceNotesLoading(t *testing.T) {
	tests := []struct {
		name          string
		fileStructure fileStructure
		expectedNotes Notes
	}{
		{
			name:          "Notes in a flat structure",
			fileStructure: flatFileStructure,
			expectedNotes: Notes{
				notes: []models.Note{
					{
						Id:         0,
						NotebookId: 1,
						Name:       "note1.md",
						Path:       "/root/book1/note1.md",
					},
					{
						Id:         1,
						NotebookId: 2,
						Name:       "note2.md",
						Path:       "/root/book2/note2.md",
					},
					{
						Id:         2,
						NotebookId: 3,
						Name:       "note3.md",
						Path:       "/root/book3/note3.md",
					},
				},
			},
		},
		{
			name:          "Notes in a nested structure",
			fileStructure: nestedFileStructure,
			expectedNotes: Notes{
				notes: []models.Note{
					{
						Id:         0,
						NotebookId: 3,
						Name:       "note3.md",
						Path:       "/root/book1/book2/book3/note3.md",
					},
					{
						Id:         1,
						NotebookId: 2,
						Name:       "note2.md",
						Path:       "/root/book1/book2/note2.md",
					},
					{
						Id:         2,
						NotebookId: 1,
						Name:       "note1.md",
						Path:       "/root/book1/note1.md",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			afs := afero.NewMemMapFs()
			makeDirStructure(test.fileStructure, "/", afs)

			//act
			fs := newFilesystemSource("/root", &afero.Afero{Fs: afs})
			assertNotesAreLike(test.expectedNotes, *fs.getNotes(), t)
		})
	}
}

func assertNotesAreLike(expected Notes, actual Notes, t *testing.T) {
	if len(expected.notes) != len(actual.notes) {
		t.Fatalf("Expected there to be [%d] notes, there were [%d]", len(expected.notes), len(actual.notes))
	}
	for i, expectedNote := range expected.notes {
		actualNote := actual.notes[i]
		if expectedNote.Name != actualNote.Name {
			t.Errorf("Expected note to have name [%s]. Had [%s]", expectedNote.Name, actualNote.Name)
		}
		if expectedNote.Id != actualNote.Id {
			t.Errorf("Expected note to have ID [%d]. Had [%d]", expectedNote.Id, actualNote.Id)
		}
		if expectedNote.Path != actualNote.Path {
			t.Errorf("Expected note to have path [%s]. Had [%s]", expectedNote.Path, actualNote.Path)
		}
	}
}

func assertBookIsLike(expected models.Notebook, actual models.Notebook, t *testing.T) {
	if expected.Name != actual.Name {
		t.Errorf("Expected book to have name [%s]. Had [%s]", expected.Name, actual.Name)
	}
	if expected.Id != actual.Id {
		t.Errorf("Expected book to have ID [%d]. Had [%d]", expected.Id, actual.Id)
	}
	if expected.ParentId != actual.ParentId {
		t.Errorf("Expected book to have ParentId [%d]. Had [%d]", expected.ParentId, actual.ParentId)
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
