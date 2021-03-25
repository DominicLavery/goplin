package ui

import (
	"dominiclavery/goplin/data"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
	"github.com/google/uuid"
)

type NotebookTree struct {
	*tview.TreeView
}

func makeTreeNode(book *data.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func add(notebooks []*data.Notebook, target *tview.TreeNode) {
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

func findExpanded(nodes []*tview.TreeNode) []uuid.UUID {
	var expanded []uuid.UUID
	for _, node := range nodes {
		if node.IsExpanded() {
			expanded = append(expanded, node.GetReference().(uuid.UUID))
		}
		if len(node.GetChildren()) > 0 {
			expanded = append(expanded, findExpanded(node.GetChildren())...)
		}
	}
	return expanded
}

func contains(s []uuid.UUID, e uuid.UUID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func applyExpand(expanded []uuid.UUID, nodes []*tview.TreeNode) {
	for _, node := range nodes {
		if contains(expanded, node.GetReference().(uuid.UUID)) {
			node.SetExpanded(true)
		}
		if len(node.GetChildren()) > 0 {
			applyExpand(expanded, node.GetChildren())
		}
	}
}

func (nt *NotebookTree) SetNotebook(notebook *data.Notebook) *tview.TreeNode {
	var notebookRoot *tview.TreeNode
	//var currentSelection uuid.UUID
	//if nt.GetCurrentNode() != nil {
	//	currentSelection = nt.GetCurrentNode().GetReference().(uuid.UUID)
	//} else {
	//	currentSelection = notebook.Id
	//}
	notebookRoot = makeTreeNode(notebook).
		SetExpanded(true).
		SetColor(tcell.ColorRed)
	add(notebook.Children, notebookRoot)

	if nt.GetRoot() != nil {
		// If there is an old tree, get its expanded and apply it to the new one
		expanded := findExpanded(nt.GetRoot().GetChildren())
		applyExpand(expanded, notebookRoot.GetChildren())
	}
	//nt.SetRoot(notebookRoot).
	//	SetCurrentNode(getSelectedNode(notebookRoot, currentSelection))
	//notesTree.SetNotes(data.GetBook(currentSelection).Notes)
	return notebookRoot
}

func MakeNotebookView() *NotebookTree {
	tree := NotebookTree{tview.NewTreeView()}

	tree.
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			book := data.GetBook(reference.(uuid.UUID))
			notesTree.SetNotes(book.Notes)
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == tree.GetRoot().GetReference() {
				return // Selecting the notebookRoot node does nothing.
			}
			children := node.GetChildren()
			if len(children) != 0 {
				node.SetExpanded(!node.IsExpanded())
			}
		})

	tree.SetBorder(true).SetTitle("Notebooks")
	return &tree
}
