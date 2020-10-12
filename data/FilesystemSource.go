package data

import (
	"dominiclavery/goplin/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FilesystemSource struct {
	path string
}

func NewFilesystemSource(path string) *FilesystemSource {
	return &FilesystemSource{path: path}
}

func (f FilesystemSource) Dataset() ([]models.Notebook, []models.Note) {
	root := f.path // TODO source this from somewhere
	notes := make([]models.Note, 0)
	notebooks := make([]models.Notebook, 0)

	notebookCount := 1
	walkFunction := func(path string, info os.FileInfo, err error) error {
		if err != nil { //TODO handle me?
			return err
		}

		if info.IsDir() && info.Name()[0:1] == "." {
			return filepath.SkipDir
		} else if info.IsDir() {
			parentId := ""
			// Root doesn't have a parentId
			if path != root {
				parentPath := filepath.Dir(path)
				parentId = models.ByPath(notebooks, parentPath).Id
			}
			notebooks = append(notebooks, models.Notebook{Name: info.Name(), Id: strconv.Itoa(notebookCount), ParentId: parentId, Path: path})
			notebookCount++
		} else if strings.HasSuffix(info.Name(), ".md") {
			parentPath := filepath.Dir(path)
			notebook := models.ByPath(notebooks, parentPath)
			notes = append(notes, models.Note{Name: info.Name(), Id: strconv.Itoa(notebookCount), NotebookId: notebook.Id, Body: path, Path: path})
		}
		return nil
	}
	filepath.Walk(root, walkFunction) //TODO handle error

	return notebooks, notes
}
