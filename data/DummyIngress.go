package data

import "dominiclavery/goplin/models"

var notebooks = []models.Notebook{
	{Id: "1", ParentId: "", Name: "."},
	{Id: "2", ParentId: "1", Name: "Book 1"},
	{Id: "3", ParentId: "2", Name: "Book 1.1"},
	{Id: "4", ParentId: "1", Name: "Book 2"},
	{Id: "5", ParentId: "2", Name: "Book 1.2"},
	{Id: "6", ParentId: "5", Name: "Book 1.2.1"},
	{Id: "7", ParentId: "4", Name: "Book 2.1"},
}

var notes = []models.Note{
	{Id: "1", NotebookId: "1", Name: "MyGreatNote1", Body: "Stuff is cool1. Here are more words \n# Test"},
	{Id: "2", NotebookId: "1", Name: "MyGreatNote2", Body: "Stuff is cool2"},
	{Id: "3", NotebookId: "1", Name: "MyGreatNote3", Body: "Stuff is cool3"},
	{Id: "4", NotebookId: "1", Name: "MyGreatNote4", Body: "Stuff is cool4"},
	{Id: "5", NotebookId: "2", Name: "MyGreatNote5", Body: "Stuff is cool5"},
	{Id: "6", NotebookId: "5", Name: "subbooking", Body: "Stuff is cool6"},
}

type DummySource struct {
}

func NewDummySource() *DummySource {
	return &DummySource{}
}

func (b DummySource) Dataset() ([]models.Notebook, []models.Note) {
	return notebooks, notes
}
