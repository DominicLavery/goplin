package data

import "dominiclavery/goplin/models"

type Source interface {
	Dataset() ([]models.Notebook, []models.Note)
}
