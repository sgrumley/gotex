package components

import (
	"log/slog"
	"sgrumley/gotex/pkg/finder"
	"sgrumley/gotex/pkg/runner"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	rootColor    = tcell.ColorRed
	unknownColor = tcell.ColorYellow
)

var (
	LevelPackage  = 1
	LevelFile     = 2
	LevelFunction = 3
	LevelCase     = 4
)

type TestTree struct {
	*tview.TreeView
}

func newTestTree(t *TUI) *TestTree {
	tt := &TestTree{
		TreeView: tview.NewTreeView(),
	}

	tt.setKeybinding(t)
	tt.SetTitle("Tests")
	SetTreeStyling(t, tt.TreeView)
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
				// TODO: this should print to the console..
				t.state.ui.result.RenderResults("Error can't get node " + node.GetReference().(finder.Node).GetName())
			}
			node.ExpandAll()
		case 'h':
			node := t.state.ui.testTree.GetCurrentNode()
			if node == nil {
				// TODO: this should print to the console..
				t.state.ui.result.RenderResults("Error can't get node " + node.GetReference().(finder.Node).GetName())
			}
			node.CollapseAll()
		// run test
		case 'r':
			t.state.ui.result.RenderResults("Testing ....")
			dataNode, ok := tt.GetCurrentNode().GetReference().(finder.Node)
			if !ok {
				t.log.Error("reference to current node is not a testable type")
				t.state.ui.result.RenderResults("Error selected node is not a test")
				return event
			}
			t.state.data.lastTest = dataNode
			output, err := dataNode.RunTest()
			if err != nil {
				t.log.Error("failed running test", slog.Any("error", err), slog.Any("output", output))
				t.state.ui.result.RenderResults(output.Result)
				t.state.ui.console.panel.UpdateMeta(t, output)
				return event
			}

			t.state.ui.result.RenderResults(output.Result)
			t.state.ui.console.panel.UpdateMeta(t, output)
			return nil
		// sync tests
		case 's':
			// NOTE: this could happen on a timer or by watching the the test files for changes
			data, err := finder.InitProject(t.log)
			if err != nil {
				t.state.ui.result.RenderResults(err.Error())
			}
			t.state.data.project = data
			t.state.ui.testTree.Populate(t)
			t.state.ui.result.RenderResults("Project has successfully refreshed")
			return nil
		// run all
		case 'A':
			output, err := runner.RunTest(runner.TestTypeProject, "", t.state.data.project.RootDir, t.state.data.project.Config)
			if err != nil {
				t.log.Error("failed running all tests", slog.Any("error", err))
				t.state.ui.console.panel.UpdateMeta(t, output)
				t.state.ui.result.RenderResults(err.Error())
				return event
			}
			t.state.ui.result.RenderResults(output.Result)
			t.state.ui.console.panel.UpdateMeta(t, output)
			return nil

		// search
		case '/':
			t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Rune() {
				case 'c':
					// NOTE: this is an example of when to return the event rather than nil, as it will be passed through and still count as text input
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
		case tcell.KeyEsc:
			t.state.ui.pages.SwitchToPage(homePage)
			return nil
		}
		return event
	})
}

func (tt *TestTree) Populate(t *TUI) {
	data := t.state.data.project
	root := tview.NewTreeNode(data.GetName()).SetColor(rootColor)
	tt.SetRoot(root)
	tt.SetCurrentNode(root)

	prefillTree(t, root, data, 0)
	// allow level 1 to be expanded
	for _, child := range root.GetChildren() {
		child.CollapseAll()
	}

	tt.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetReference() == nil {
			return
		}

		node.SetExpanded(!node.IsExpanded())
	})
}

func prefillTree(t *TUI, target *tview.TreeNode, n finder.Node, lvl int) {
	children := n.GetChildren()
	for _, child := range children {
		node := tview.NewTreeNode(child.GetName())
		node.SetReference(child)
		node.SetSelectable(true)

		// node level styling
		// TODO: consider useing SetPrefixes: https://pkg.go.dev/github.com/rivo/tview#TreeView

		switch lvl + 1 {
		case LevelPackage:
			node.SetText(" " + node.GetText())
			node.SetColor(t.theme.Package)
		case LevelFile:
			node.SetText(" " + node.GetText())
			node.SetColor(t.theme.File)
		case LevelFunction:
			node.SetText("󰡱 " + node.GetText())
			node.SetColor(t.theme.Function)
		case LevelCase:
			node.SetText("󰙨 " + node.GetText())
			node.SetColor(t.theme.Case)
		default:
			node.SetColor(unknownColor)
		}

		target.AddChild(node)
		prefillTree(t, node, child, lvl+1)
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
