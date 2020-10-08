package data

import "dominiclavery/goplin/models"

var notebooks = []models.Notebook{
	{"1", "", "."},
	{"2", "1", "Book 1"},
	{"3", "2", "Book 1.1"},
	{"4", "1", "Book 2"},
	{"5", "2", "Book 1.2"},
	{"6", "5", "Book 1.2.1"},
	{"7", "4", "Book 2.1"},
}

var notes = []models.Note{
	{"1", "1", "MyGreatNote1", "Stuff is cool1. Here are more words \n# Test"},
	{"2", "1", "MyGreatNote2", "Stuff is cool2"},
	{"3", "1", "MyGreatNote3", "Stuff is cool3"},
	{"4", "1", "MyGreatNote4", "Stuff is cool4"},
	{"5", "2", "MyGreatNote5", "Stuff is cool5"},
	{"6", "5", "subbooking", "Stuff is cool6"},
}

type DummySource struct {

}

func NewDummySource() *DummySource {
	return &DummySource{}
}

func (b DummySource) GetDataset() ([]models.Notebook,[]models.Note) {
	return notebooks, notes
}
