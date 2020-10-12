package models

type Notebook struct {
	Id       string
	ParentId string
	Path     string
	Name     string
}

//TODO handle missing
//ByPath finds a notebook for a given path.
func ByPath(notebooks []Notebook, path string) Notebook {
	var book Notebook
	for _, notebook := range notebooks {
		if notebook.Path == path {
			book = notebook
		}
	}
	return book
}