package components

import (
	"fmt"
	"log/slog"

	"github.com/sgrumley/gotex/pkg/finder"
	"github.com/sgrumley/gotex/pkg/runner"
)

func SyncProject(t *TUI) {
	// NOTE: this could happen on a timer or by watching the the test files for changes
	data, err := finder.InitProject(t.log)
	if err != nil {
		t.state.ui.result.RenderResults(err.Error())
	}
	t.state.data.project = data
	t.state.ui.testTree.Populate(t)
	t.state.ui.result.RenderResults("Project has successfully refreshed")
}

func RunTest(t *TUI) error {
	t.state.ui.result.RenderResults("Test is running")
	dataNode, ok := t.state.ui.testTree.GetCurrentNode().GetReference().(finder.Node)
	if !ok {
		t.log.Error("reference to current node is not a testable type")
		t.state.ui.result.RenderResults("Error selected node is not a test")
		return fmt.Errorf("invalid node")
	}

	go func() {
		t.state.data.lastTest = dataNode
		output, err := dataNode.RunTest()
		if err != nil {
			t.log.Error("failed running test", slog.Any("error", err), slog.Any("output", output))
			t.state.ui.result.RenderResults(output.Result)
			t.state.ui.console.panel.UpdateMeta(t, output)
			return
		}

		t.state.ui.result.RenderResults(output.Result)
		t.state.ui.console.panel.UpdateMeta(t, output)
	}()
	return nil
}

func RunAllTests(t *TUI) error {
	t.state.ui.result.RenderResults("Test is running")

	go func() {
		output, err := runner.RunTest(runner.TestTypeProject, "", t.state.data.project.RootDir, t.state.data.project.Config)
		if err != nil {
			t.log.Error("failed running all tests", slog.Any("error", err))
			t.state.ui.console.panel.UpdateMeta(t, output)
			t.state.ui.result.RenderResults(err.Error())
			return
		}
		t.state.ui.result.RenderResults(output.Result)
		t.state.ui.console.panel.UpdateMeta(t, output)
	}()

	return nil
}

func RerunTest(t *TUI) error {
	t.state.ui.result.RenderResults("Rerunning test")
	t.log.Error("this should not have run")

	node := t.state.data.lastTest
	if node == nil {
		t.state.ui.result.RenderResults("failed to run last test. Make sure you run a test before rerunning")
		t.log.Error("attempted test rerun, but no test has previously been run")
		return fmt.Errorf("no previously run test")
	}

	output, err := node.RunTest()
	if err != nil {
		t.log.Error("failed to re run valid test", slog.Any("error", err))
		t.state.ui.console.panel.UpdateMeta(t, output)
		t.state.ui.result.RenderResults(err.Error())
		return err
	}
	t.state.ui.result.RenderResults(output.Result)
	t.state.ui.console.panel.UpdateMeta(t, output)

	return nil
}
