package data

import (
	"dominiclavery/goplin/models"
	"github.com/spf13/afero"
	"testing"
)

func TestNewFilesystemSourceDeepDirectory_flat(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)
	_ = afs.Mkdir("/root/book1", 0644)
	_ = afs.Mkdir("/root/book2", 0644)
	_ = afs.Mkdir("/root/book3", 0644)
	_ = afs.Mkdir("/root/.ignored", 0644)

	expectedBooks := models.Notebook{
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
	}

	//act
	fs := newFilesystemSource("/root", &afero.Afero{Fs: afs})

	//assert
	if fs.getNotebooks().highestNotebookId != 4 {
		t.Errorf("Filesystem did not find the correct amount of notebooks. Expected %d, got %d", 4, fs.getNotebooks().highestNotebookId)
	}
	assertBookIsLike(expectedBooks, fs.getNotebooks().notebookRoot, t)
}

func TestNewFilesystemSourceDeepDirectory_deep(t *testing.T) {
	//arrange
	afs := afero.NewMemMapFs()
	_ = afs.Mkdir("/root", 0644)
	_ = afs.Mkdir("/root/book1", 0644)
	_ = afs.Mkdir("/root/book1/book2", 0644)
	_ = afs.Mkdir("/root/book1/book2/book3", 0644)
	expectedBooks := models.Notebook{
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
	}

	//act
	fs := newFilesystemSource("/root", &afero.Afero{Fs: afs})

	//assert
	if fs.getNotebooks().highestNotebookId != 4 {
		t.Errorf("Filesystem did not find the correct amount of notebooks. Expected %d, got %d", 4, fs.getNotebooks().highestNotebookId)
	}

	assertBookIsLike(expectedBooks, fs.getNotebooks().notebookRoot, t)
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

func assertBookIsLike(expected models.Notebook, actual models.Notebook, t *testing.T) {
	if expected.Name != actual.Name {
		t.Errorf("Expected created book to have name [%s]. Had [%s]", expected.Name, actual.Name)
	}
	if expected.Id != actual.Id {
		t.Errorf("Expected created book to have ID [%d]. Had [%d]", expected.Id, actual.Id)
	}
	if expected.ParentId != actual.ParentId {
		t.Errorf("Expected created book to have ParentId [%d]. Had [%d]", expected.ParentId, actual.ParentId)
	}
	if expected.Children != nil {
		if actual.Children == nil {
			t.Errorf("Expected created book to have [%d] childrem. Had 0", len(expected.Children))
		}
		if len(expected.Children) != len(actual.Children) {
			t.Errorf("Expected created book to have [%d] childrem. Had [%d]", len(expected.Children), len(actual.Children))
		}
		for i, book := range expected.Children {
			assertBookIsLike(book, actual.Children[i], t)
		}
	}
}
