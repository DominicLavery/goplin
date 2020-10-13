package ui

import (
	"dominiclavery/goplin/models"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func makeTreeNode(book models.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func add(notebooks []models.Notebook, target *tview.TreeNode, parentId string) {
	for _, book := range notebooks {
		if book.ParentId == parentId {
			node := makeTreeNode(book)
			target.AddChild(node)
		}
	}
}

func MakeNotebookView(notebooks []models.Notebook, notes []models.Note, updateNotesView func([]models.Note)) *tview.TreeView {
	var notebookRoot *tview.TreeNode
	if len(notebooks) == 0 {
		notebookRoot = makeTreeNode(models.Notebook{Id: "0", Name: "No notebooks found"})
	} else {
		notebookRoot = makeTreeNode(notebooks[0])
	}
	notebookRoot.SetColor(tcell.ColorRed)

	add(notebooks, notebookRoot, "1")
	updateNotesView(models.ByNotebookId(notes, "1"))

	notebookTree := tview.NewTreeView().
		SetRoot(notebookRoot).
		SetCurrentNode(notebookRoot).
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			updateNotesView(models.ByNotebookId(notes, reference.(string)))
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == "" {
				return // Selecting the notebookRoot node does nothing.
			}
			children := node.GetChildren()
			if len(children) == 0 {
				add(notebooks, node, reference.(string))
			} else {
				node.SetExpanded(!node.IsExpanded())
			}
		})
	notebookTree.SetBorder(true).SetTitle("Notebooks")
	return notebookTree
}
