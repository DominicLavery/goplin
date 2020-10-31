package models

import "io"

type Note struct {
	Id         int
	NotebookId int
	Name       string
	Body       io.Reader
	Path       string
}
