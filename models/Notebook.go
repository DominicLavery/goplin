package models

type Notebook struct {
	Id       int
	ParentId int
	Name     string
	Path     string
	Children []Notebook
}
