package models

type Note struct {
	Id         string
	NotebookId string
	Name       string
	Body       string
}

func ByNotebookId(notes []Note, notebookId string) []Note {
	filtered := make([]Note, 0)
	for _, note := range notes {
		if note.NotebookId == notebookId {
			filtered = append(filtered, note)
		}
	}
	return filtered
}
