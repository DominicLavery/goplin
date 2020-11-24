package data

import (
	"dominiclavery/goplin/models"
	"sync"
)

type Notebooks struct {
	notebookRoot      models.Notebook
	mu                sync.Mutex
	highestNotebookId int
}
