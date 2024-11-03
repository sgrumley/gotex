package components

import (
	"sgrumley/gotex/pkg/finder"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	rootColor     = tcell.ColorRed
	fileColor     = tcell.ColorBlue
	functionColor = tcell.ColorPink
	caseColor     = tcell.ColorGreen
	unknownColor  = tcell.ColorYellow
)

var (
	LevelFile     = 1
	LevelFunction = 2
	LevelCase     = 3
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
		case 'r':
			t.state.result.RenderResults("Testing ....")
			dataNode, ok := tt.GetCurrentNode().GetReference().(finder.Node)
			if !ok {
				t.state.result.RenderResults("Error selected node is not a test")
				return event
			}
			t.state.lastTest = dataNode
			output, err := dataNode.RunTest()
			if err != nil {
				t.state.result.RenderResults(err.Error())
				return event
			}

			t.state.result.RenderResults(output)
		case '/':
			// TODO: search
		}
		// keybinding for special keys
		switch event.Key() {
		case tcell.KeyCtrlR:
			t.state.result.RenderResults("Testing ....")
			time.Sleep(1 * time.Second)

			node := t.state.lastTest
			if node == nil {
				t.state.result.RenderResults("Failed to run last test. Make sure you run a test before rerunning")
				return event
			}
			output, err := node.RunTest()
			if err != nil {
				t.state.result.RenderResults(err.Error())
				return event
			}

			t.state.result.RenderResults(output)
		case tcell.KeyCtrlU:
			currentPosition, _ := t.state.result.GetScrollOffset()
			t.state.result.ScrollTo(currentPosition-10, 0)
		case tcell.KeyCtrlD:
			currentPosition, _ := t.state.result.GetScrollOffset()
			t.state.result.ScrollTo(currentPosition+10, 0)
		}
		return event
	})
}

func (tt *TestTree) Populate(t *TUI) {
	data := t.state.resources.data
	root := tview.NewTreeNode(data.GetName()).SetColor(rootColor)
	tt.SetRoot(root)
	tt.SetCurrentNode(root)

	add(t, root, data)

	tt.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			dataNode := reference.(finder.Node)
			add(t, node, dataNode)
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})
}

func add(t *TUI, target *tview.TreeNode, n finder.Node) {
	children := n.GetChildren()
	for _, child := range children {
		node := tview.NewTreeNode(child.GetName())
		node.SetReference(child)
		node.SetSelectable(true)
		// node.SetSelectable(child.HasChildren()) // NOTE: this makes cases unselectable

		// node level styling
		switch target.GetLevel() + 1 {
		case LevelFile:
			node.SetText("  " + node.GetText())
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
	}
}
