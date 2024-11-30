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
		// run test
		case 'l':
			node := t.state.testTree.GetCurrentNode()
			if node == nil {
				// TODO: this should print to the console..
				t.state.result.RenderResults("Error can't get node " + node.GetReference().(finder.Node).GetName())
			}
			node.ExpandAll()
		case 'h':
			node := t.state.testTree.GetCurrentNode()
			if node == nil {
				// TODO: this should print to the console..
				t.state.result.RenderResults("Error can't get node " + node.GetReference().(finder.Node).GetName())
			}
			node.CollapseAll()

		case 'r':
			t.state.result.RenderResults("Testing ....")
			dataNode, ok := tt.GetCurrentNode().GetReference().(finder.Node)
			if !ok {
				t.log.Error("reference to current node is not a testable type")
				t.state.result.RenderResults("Error selected node is not a test")
				return event
			}
			t.state.lastTest = dataNode
			output, err := dataNode.RunTest()
			if err != nil {
				t.log.Error("failed running test", slog.Any("error", err), slog.Any("output", output))
				t.state.result.RenderResults(err.Error())
				return event
			}

			t.state.result.RenderResults(output)
			return nil
		// rerun last test
		case 'R':
			// FIX: need a way to show the user that the test has been rerun/ is rerunning
			// maybe a job for the meta console?
			t.state.result.RenderResults("Rerunning test")
			t.log.Error("this should not have run")

			node := t.state.lastTest
			if node == nil {
				t.state.result.RenderResults("failed to run last test. Make sure you run a test before rerunning")
				t.log.Error("attempted test rerun, but no test has previously been run")
				return event
			}

			output, err := node.RunTest()
			if err != nil {
				t.log.Error("failed to re run valid test", slog.Any("error", err))
				t.state.result.RenderResults(err.Error())
				return event
			}
			t.state.result.RenderResults(output)
			return nil

		// sync tests
		case 's':
			// NOTE: this could happen on a timer or by watching the the test files for changes
			data, err := finder.InitProject(t.log)
			if err != nil {
				t.state.result.RenderResults(err.Error())
			}
			t.state.resources.data = data
			t.state.testTree.Populate(t)
			t.state.result.RenderResults("Project has successfully refreshed")
			return nil
		// run all
		case 'A':
			output, err := runner.RunTest(runner.TEST_TYPE_PROJECT, "", t.state.resources.data.RootDir, t.state.resources.data.Config)
			if err != nil {
				t.log.Error("failed running all tests", slog.Any("error", err))
				t.state.result.RenderResults(err.Error())
				return event
			}
			t.state.result.RenderResults(output)
			return nil

		// search
		case '/':
			t.state.pages.ShowPage(searchPage)
			t.app.SetFocus(t.state.search.input)
			return nil
		}
		// keybinding for special keys
		switch event.Key() {
		case tcell.KeyCtrlU:
			currentPosition, _ := t.state.result.GetScrollOffset()
			t.state.result.ScrollTo(currentPosition-10, 0)
			return nil
		case tcell.KeyCtrlD:
			currentPosition, _ := t.state.result.GetScrollOffset()
			t.state.result.ScrollTo(currentPosition+10, 0)
			return nil
		case tcell.KeyEsc:
			t.state.pages.SwitchToPage(homePage)
			return nil
		}
		return event
	})
}

func (tt *TestTree) Populate(t *TUI) {
	data := t.state.resources.data
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

func search(tree *tview.TreeView, searchString string, t *TUI) bool {
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
