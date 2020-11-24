package data

import (
	"dominiclavery/goplin/models"
	"sync"
)

type Notes struct {
	notes         []models.Note
	mu            sync.Mutex
	openBookId    int
	openNote      *models.Note
	highestNoteId int
}
