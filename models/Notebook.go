package models

type Notebook struct {
	Id       int
	ParentId int
	Name     string
	Children []Notebook
}
