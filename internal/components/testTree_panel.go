package components

import (
	"context"
	"strings"

	"github.com/sgrumley/gotex/pkg/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var unknownColor = tcell.ColorYellow

type TestTree struct {
	*tview.TreeView
}

func newTestTree(ctx context.Context, t *TUI) *TestTree {
	tt := &TestTree{
		TreeView: tview.NewTreeView(),
	}

	tt.setKeybinding(ctx, t)
	tt.SetTitle("Tests")
	tt.SetBorder(true)
	tt.Populate(t)
	tt.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
		dynamicResize(tt, t)
	})

	return tt
}

func (tt *TestTree) setKeybinding(ctx context.Context, t *TUI) {
	tt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(ctx, event)

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
			dynamicResize(tt, t)
		case 'h':
			node := t.state.ui.testTree.GetCurrentNode()
			if node == nil {
				t.state.ui.result.RenderResults("Error can't get node " + node.GetReference().(models.Node).GetName())
			}
			node.CollapseAll()
			dynamicResize(tt, t)
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
			RunTest(ctx, t)
			return nil
		case 's':
			SyncProject(context.Background(), t)
			return nil
		case 'A':
			RunAllTests(ctx, t)
			return nil

		// search
		case '/':
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

func (tt *TestTree) Populate(t *TUI) error {
	err := models.GenerateTree(t.state.data.project)
	if err != nil {
		return err
	}
	tree := t.state.data.project.Tree
	rootViewNode := convertNode(t, tree.RootNode)
	tt.SetRoot(rootViewNode)
	tt.SetCurrentNode(rootViewNode)

	for _, child := range rootViewNode.GetChildren() {
		child.CollapseAll()
	}

	return nil
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calculateOptimalTreeWidth(node *tview.TreeNode, depth int) int {
	if node == nil {
		return 0
	}

	// Calculate width for this node
	nodeText := tview.TaggedStringWidth(node.GetText())
	indentWidth := depth * 4 // Tree indentation
	nodeWidth := nodeText + indentWidth

	maxWidth := nodeWidth

	// Only include children if this node is expanded
	if node.IsExpanded() {
		children := node.GetChildren()
		for _, child := range children {
			childWidth := calculateOptimalTreeWidth(child, depth+1)
			if childWidth > maxWidth {
				maxWidth = childWidth
			}
		}
	}

	return maxWidth
}

func dynamicResize(tt *TestTree, t *TUI) {
	_, _, currentWidth, _ := tt.GetRect()
	minWidth := 45
	maxWidth := 100
	padding := 6
	optimalWidth := calculateOptimalTreeWidth(tt.GetRoot(), 0) + padding

	if optimalWidth < minWidth {
		optimalWidth = minWidth
	} else if optimalWidth > maxWidth {
		optimalWidth = maxWidth
	}

	// Only resize if there's a meaningful difference (reduces flickering)
	if abs(optimalWidth-currentWidth) > 2 {
		t.state.ui.homeLayout.ResizeItem(tt.TreeView, optimalWidth, 0)
	}
}
