package data

import (
	"dominiclavery/goplin/models"
)

type Notes struct {
	notes         []models.Note
	openNote      *models.Note
	highestNoteId int
}
