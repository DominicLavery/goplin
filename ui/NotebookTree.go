package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/models"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

func makeTreeNode(book models.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func add(notebooks []models.Notebook, target *tview.TreeNode) {
	for _, book := range notebooks {
		node := makeTreeNode(book)
		if book.Children != nil && len(book.Children) > 0 {
			add(book.Children, node)
		}
		node.SetExpanded(false)
		target.AddChild(node)
	}
}

func findExpanded(nodes []*tview.TreeNode) []int {
	var expanded []int
	for _, node := range nodes {
		if node.IsExpanded() {
			expanded = append(expanded, node.GetReference().(int))
		}
		if len(node.GetChildren()) > 0 {
			expanded = append(expanded, findExpanded(node.GetChildren())...)
		}
	}
	return expanded
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func applyExpand(expanded []int, nodes []*tview.TreeNode) {
	for _, node := range nodes {
		if contains(expanded, node.GetReference().(int)) {
			node.SetExpanded(true)
		}
		if len(node.GetChildren()) > 0 {
			applyExpand(expanded, node.GetChildren())
		}
	}
}

func MakeNotebookView(source data.Source) *tview.TreeView {
	notebookTree := tview.NewTreeView()

	source.Notebooks(func(notebook models.Notebook) {
		var notebookRoot *tview.TreeNode
		notebookRoot = makeTreeNode(notebook).
			SetExpanded(true).
			SetColor(tcell.ColorRed)
		add(notebook.Children, notebookRoot)
		notebookRoot.SetExpanded(true)
		if notebookTree.GetRoot() != nil {
			expanded := findExpanded(notebookTree.GetRoot().GetChildren())
			applyExpand(expanded, notebookRoot.GetChildren())
		}
		notebookTree.SetRoot(notebookRoot).
			SetCurrentNode(notebookRoot)
	})

	notebookTree.
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			source.OpenBook(reference.(int))
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == 0 {
				return // Selecting the notebookRoot node does nothing.
			}
			children := node.GetChildren()
			if len(children) != 0 {
				node.SetExpanded(!node.IsExpanded())
			}
		})
	notebookTree.SetBorder(true).SetTitle("Notebooks")
	return notebookTree
}
