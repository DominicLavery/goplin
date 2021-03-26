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

var expandedMap map[uuid.UUID]bool

func init() {
	expandedMap = make(map[uuid.UUID]bool)
}

func makeTreeNode(book *data.Notebook) *tview.TreeNode {
	return tview.NewTreeNode(book.Name).
		SetReference(book.Id).
		SetSelectable(true)
}

func addChildBooks(notebooks []*data.Notebook, target *tview.TreeNode) {
	for _, book := range notebooks {
		node := makeTreeNode(book)
		addChildBooks(book.Children, node)
		target.AddChild(node)
		setExpanded(book, node)
	}
}

func setExpanded(book *data.Notebook, node *tview.TreeNode) {
	if _, present := expandedMap[book.Id]; !present {
		expandedMap[book.Id] = false
	}
	node.SetExpanded(expandedMap[book.Id])
}

func findNode(target *tview.TreeNode, selected uuid.UUID) *tview.TreeNode {
	if target.GetReference() == selected {
		return target
	}
	for _, child := range target.GetChildren() {
		found := findNode(child, selected)
		if found != nil {
			return found
		}
	}
	return nil
}

func getNodeById(target *tview.TreeNode, selected uuid.UUID) *tview.TreeNode {
	found := findNode(target, selected)
	if found != nil {
		return found
	}
	//Nothing found, default to the root
	return target
}

func (nt *NotebookTree) SetDataTree(tree map[string]data.SourceDataTree) {
	root := tview.NewTreeNode(".").
		SetReference(uuid.New()).
		SetSelectable(false).
		SetColor(tcell.ColorWhite)

	var currentNode *tview.TreeNode
	var currentId uuid.UUID
	var roots []*tview.TreeNode
	for name, sdt := range tree {
		node := tview.NewTreeNode(name).
			SetReference(sdt.NotebookRoot.Id).
			SetSelectable(true).
			SetColor(tcell.ColorBlue)
		setExpanded(sdt.NotebookRoot, node)
		addChildBooks(sdt.NotebookRoot.Children, node)
		roots = append(roots, node)
	}
	root.SetChildren(roots)
	if nt.GetRoot() != nil {
		currentId = nt.GetCurrentNode().GetReference().(uuid.UUID)
		currentNode = findNode(root, currentId)
	}

	if currentNode == nil && roots != nil {
		currentId = roots[0].GetReference().(uuid.UUID)
		currentNode = roots[0]
	}
	nt.SetRoot(root).
		SetCurrentNode(currentNode)
	notesTree.SetNotes(data.GetBook(currentId).Notes)
}

func MakeNotebookView() *NotebookTree {
	tree := NotebookTree{tview.NewTreeView()}
	tree.SetTopLevel(1)
	tree.
		SetChangedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			book := data.GetBook(reference.(uuid.UUID))
			notesTree.SetNotes(book.Notes)
		}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			children := node.GetChildren()
			if len(children) != 0 {
				node.SetExpanded(!node.IsExpanded())
				expandedMap[node.GetReference().(uuid.UUID)] = node.IsExpanded()
			}
		})

	tree.SetBorder(true).SetTitle("Notebooks")
	return &tree
}

func (nt *NotebookTree) expandCurrentNode() {
	node := nt.GetCurrentNode()
	node.SetExpanded(!node.IsExpanded())
	expandedMap[node.GetReference().(uuid.UUID)] = node.IsExpanded()
}
