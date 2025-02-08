package components

import (
	"context"
	"strings"

	"github.com/sgrumley/gotex/pkg/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// rootColor    = tcell.ColorRed
var unknownColor = tcell.ColorYellow

// var (
// 	LevelPackage  = 1
// 	LevelFile     = 2
// 	LevelFunction = 3
// 	LevelCase     = 4
// )

type TestTree struct {
	*tview.TreeView
}

func newTestTree(t *TUI) *TestTree {
	tt := &TestTree{
		TreeView: tview.NewTreeView(),
	}

	tt.setKeybinding(t)
	tt.SetTitle("Tests")
	tt.SetBorder(true)
	tt.Populate(t)

	return tt
}

func (tt *TestTree) setKeybinding(t *TUI) {
	tt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)

		// keybinding for single keys
		switch event.Rune() {
		// tree navigation
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'l':
			node := t.state.ui.testTree.GetCurrentNode()
			if node == nil {
				t.state.ui.result.RenderResults("Error can't get node " + node.GetReference().(models.Node).GetName())
			}
			node.ExpandAll()
		case 'h':
			node := t.state.ui.testTree.GetCurrentNode()
			if node == nil {
				t.state.ui.result.RenderResults("Error can't get node " + node.GetReference().(models.Node).GetName())
			}
			node.CollapseAll()
			// NOTE: This should jump through the list rather than scroll
		// case 'g':
		// 	if t.state.ui.lastKey == 'g' {
		// 		t.state.ui.result.ScrollTo(0, 0)
		// 		t.state.ui.lastKey = 0 // Reset last key
		// 	} else {
		// 		t.state.ui.lastKey = 'g'
		// 	}
		// 	return nil
		// case 'G':
		// 	t.state.ui.result.ScrollToEnd()
		// 	return nil
		case 'r':
			RunTest(t)
			return nil
		case 's':
			SyncProject(context.Background(), t)
			return nil
		case 'A':
			RunAllTests(t)
			return nil

		// search
		case '/':
			// TODO: update with page system
			// NOTE: this is an example of when to return the event rather than nil, as it will be passed through and still count as text input
			// upon close setGlobalKeybinding() is called to undo this
			t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Rune() {
				case 'c':
					return event
				}
				return event
			})

			t.state.ui.pages.ShowPage(searchPage)
			t.app.SetFocus(t.state.ui.search.input)
			return nil
		}

		// keybinding for special keys
		switch event.Key() {
		case tcell.KeyCtrlU:
			currentPosition, _ := t.state.ui.result.GetScrollOffset()
			t.state.ui.result.ScrollTo(currentPosition-10, 0)
			return nil
		case tcell.KeyCtrlD:
			currentPosition, _ := t.state.ui.result.GetScrollOffset()
			t.state.ui.result.ScrollTo(currentPosition+10, 0)
			return nil
			// NOTE: Temp fix until dynamic resizing
		case tcell.KeyCtrlH:
			_, horizontalPosition := t.state.ui.result.GetScrollOffset()
			t.state.ui.result.ScrollTo(0, horizontalPosition-10)
			return nil
			// NOTE: Temp fix until dynamic resizing
		case tcell.KeyCtrlL:
			_, horizontalPosition := t.state.ui.result.GetScrollOffset()
			t.state.ui.result.ScrollTo(0, horizontalPosition+10)
			return nil
		case tcell.KeyEsc:
			t.state.ui.pages.SwitchToPage(homePage)
			return nil
		}
		return event
	})
}

func (tt *TestTree) Populate(t *TUI) {
	models.GenerateTree(t.state.data.project)
	tree := t.state.data.project.Tree
	rootViewNode := convertNode(t, tree.RootNode)
	tt.SetRoot(rootViewNode)
	tt.SetCurrentNode(rootViewNode)

	for _, child := range rootViewNode.GetChildren() {
		child.CollapseAll()
	}

	// Set up the selection handler
	tt.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
}

func convertNode(t *TUI, node *models.NodeTree) *tview.TreeNode {
	if node == nil {
		return nil
	}

	// Create the tview node
	viewNode := tview.NewTreeNode(node.Data.GetName())
	viewNode.SetReference(node.Data)
	viewNode.SetSelectable(true)
	nodeStyling(t, viewNode, node)

	// Convert all children
	for _, child := range node.Children {
		childViewNode := convertNode(t, child)
		if childViewNode != nil {
			viewNode.AddChild(childViewNode)
		}
	}

	return viewNode
}

func nodeStyling(t *TUI, node *tview.TreeNode, dnode *models.NodeTree) {
	switch dnode.Type {
	case models.NODE_TYPE_PROJECT:
		node.SetText("  " + node.GetText())
		node.SetColor(t.theme.Project)
	case models.NODE_TYPE_DIRECTORY:
		node.SetText(" " + node.GetText())
		node.SetColor(t.theme.Directory)
	case models.NODE_TYPE_PACKAGE:
		node.SetText(" " + node.GetText())
		node.SetColor(t.theme.Package)
	case models.NODE_TYPE_FILE:
		node.SetText(" " + node.GetText())
		node.SetColor(t.theme.File)
	case models.NODE_TYPE_FUNCTION:
		node.SetText("󰡱 " + node.GetText())
		node.SetColor(t.theme.Function)
	case models.NODE_TYPE_CASE:
		node.SetText("󰙨 " + node.GetText())
		node.SetColor(t.theme.Case)
	default:
		node.SetColor(unknownColor)
	}
}

// func (tt *TestTree) Populate(t *TUI) {
// 	data := t.state.data.project
// 	root := tview.NewTreeNode(data.GetName()).SetColor(rootColor)
// 	tt.SetRoot(root)
// 	tt.SetCurrentNode(root)
//
// 	prefillTree(t, root, data, 0)
// 	// allow level 1 to be expanded
// 	for _, child := range root.GetChildren() {
// 		child.CollapseAll()
// 	}
//
// 	tt.SetSelectedFunc(func(node *tview.TreeNode) {
// 		if node.GetReference() == nil {
// 			return
// 		}
//
// 		node.SetExpanded(!node.IsExpanded())
// 	})
// }
//
// func prefillTree(t *TUI, target *tview.TreeNode, n models.Node, lvl int) {
// 	children := n.GetChildren()
// 	for _, child := range children {
// 		node := tview.NewTreeNode(child.GetName())
// 		node.SetReference(child)
// 		node.SetSelectable(true)
//
// 		// node level styling
// 		// TODO: consider useing SetPrefixes: https://pkg.go.dev/github.com/rivo/tview#TreeView
//
// 		switch lvl + 1 {
// 		case LevelPackage:
// 			node.SetText(" " + node.GetText())
// 			node.SetColor(t.theme.Package)
// 		case LevelFile:
// 			node.SetText(" " + node.GetText())
// 			node.SetColor(t.theme.File)
// 		case LevelFunction:
// 			node.SetText("󰡱 " + node.GetText())
// 			node.SetColor(t.theme.Function)
// 		case LevelCase:
// 			node.SetText("󰙨 " + node.GetText())
// 			node.SetColor(t.theme.Case)
// 		default:
// 			node.SetColor(unknownColor)
// 		}
//
// 		target.AddChild(node)
// 		prefillTree(t, node, child, lvl+1)
// 	}
// }

func search(tree *tview.TreeView, searchString string) bool {
	var matchedNode *tview.TreeNode
	var searchAndExpand func(node *tview.TreeNode, parents []*tview.TreeNode) bool

	searchAndExpand = func(node *tview.TreeNode, parents []*tview.TreeNode) bool {
		if strings.Contains(strings.ToLower(node.GetText()), strings.ToLower(searchString)) {
			// Expand all parent nodes
			for _, parent := range parents {
				parent.SetExpanded(true)
			}
			matchedNode = node
			return true // Stop traversing
		}

		// Traverse children
		for _, child := range node.GetChildren() {
			if searchAndExpand(child, append(parents, node)) {
				return true // Stop traversing
			}
		}

		return false
	}

	// Start traversal from the root
	root := tree.GetRoot()
	if root != nil {
		searchAndExpand(root, nil)
	}

	if matchedNode != nil {
		tree.SetCurrentNode(matchedNode)
	}

	return true
}
