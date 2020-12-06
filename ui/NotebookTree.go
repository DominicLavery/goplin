package ui

import (
	"dominiclavery/goplin/data"
	"dominiclavery/goplin/models"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

type NotebookTree struct {
	*tview.TreeView
}

func makeTreeNode(book models.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func add(notebooks []models.Notebook, target *tview.TreeNode) {
	for _, book := range notebooks {
		node := makeTreeNode(book)
		add(book.Children, node)
		node.SetExpanded(false)
		target.AddChild(node)
	}
}

func lookForSelected(target *tview.TreeNode, selected interface{}) *tview.TreeNode {
	if target.GetReference() == selected {
		return target
	}
	for _, child := range target.GetChildren() {
		found := lookForSelected(child, selected)
		if found != nil {
			return found
		}
	}
	return nil
}

func getSelectedNode(target *tview.TreeNode, selected interface{}) *tview.TreeNode {
	found := lookForSelected(target, selected)
	if found != nil {
		return found
	}
	//Nothing found, default to the root
	return target
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

func (nt *NotebookTree) SetNotebook(notebook models.Notebook) {
	var notebookRoot *tview.TreeNode
	var currentSelection interface{}
	if nt.GetCurrentNode() != nil {
		currentSelection = nt.GetCurrentNode().GetReference()
	} else {
		currentSelection = 0
	}
	notebookRoot = makeTreeNode(notebook).
		SetExpanded(true).
		SetColor(tcell.ColorRed)
	add(notebook.Children, notebookRoot)
	notebookRoot.SetExpanded(true)
	if nt.GetRoot() != nil {
		expanded := findExpanded(nt.GetRoot().GetChildren())
		applyExpand(expanded, notebookRoot.GetChildren())
	}
	nt.SetRoot(notebookRoot).
		SetCurrentNode(getSelectedNode(notebookRoot, currentSelection))
}

func MakeNotebookView() *NotebookTree {
	notebookTree := NotebookTree{tview.NewTreeView()}

	notebookTree.
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			data.OpenNotebooksChan <- reference.(int)
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == data.RootId {
				return // Selecting the notebookRoot node does nothing.
			}
			children := node.GetChildren()
			if len(children) != 0 {
				node.SetExpanded(!node.IsExpanded())
			}
		})

	notebookTree.SetBorder(true).SetTitle("Notebooks")
	return &notebookTree
}
