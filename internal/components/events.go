package components

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sgrumley/gotex/pkg/models"
	"github.com/sgrumley/gotex/pkg/runner"
	"github.com/sgrumley/gotex/pkg/scanner"
	"github.com/sgrumley/gotex/pkg/slogger"
)

func SyncProject(ctx context.Context, t *TUI) {
	// NOTE: this could happen on a timer or by watching the the test files for changes
	data, err := scanner.Scan(ctx, t.state.data.project.Config, t.state.data.project.RootDir)
	if err != nil {
		t.state.ui.result.RenderResults(err.Error())
	}
	t.state.data.project = data
	t.state.ui.testTree.Populate(t)
	t.state.ui.result.RenderResults("Project has successfully refreshed")
}

func RunTest(ctx context.Context, t *TUI) error {
	t.state.ui.result.RenderResults("Test is running")
	dataNode, ok := t.state.ui.testTree.GetCurrentNode().GetReference().(models.Node)
	if !ok {
		// t.log.Error("reference to current node is not a testable type")
		t.state.ui.result.RenderResults("Error selected node is not a test")
		return fmt.Errorf("invalid node")
	}

	go func() {
		log, err := slogger.FromContext(ctx)
		if err != nil {
			t.state.ui.result.RenderResults(err.Error())
			return
		}
		t.state.data.lastTest = dataNode
		output, err := dataNode.RunTest(ctx)
		if err != nil {
			log.Error("failed running test", err, slog.Any("output", output))
			if output == nil {
				t.state.ui.result.RenderResults("nil output")
				return
			}
			t.state.ui.result.RenderResults(output.Result)
			t.state.ui.console.panel.UpdateMeta(t, output)
			return
		}

		t.state.ui.result.RenderResults(output.Result)
		t.state.ui.console.panel.UpdateMeta(t, output)
	}()
	return nil
}

func RunAllTests(ctx context.Context, t *TUI) error {
	t.state.ui.result.RenderResults("Test is running")

	go func() {
		log, err := slogger.FromContext(ctx)
		if err != nil {
			t.state.ui.result.RenderResults(err.Error())
			return
		}
		output, err := runner.RunTest(ctx, runner.TestTypeProject, "", t.state.data.project.RootDir, t.state.data.project.Config)
		if err != nil {
			log.Error("failed running all tests", err)
			t.state.ui.console.panel.UpdateMeta(t, output)
			t.state.ui.result.RenderResults(err.Error())
			return
		}
		t.state.ui.result.RenderResults(output.Result)
		t.state.ui.console.panel.UpdateMeta(t, output)
	}()

	return nil
}

func RerunTest(ctx context.Context, t *TUI) error {
	t.state.ui.result.RenderResults("Rerunning test")
	log, err := slogger.FromContext(ctx)
	if err != nil {
		t.state.ui.result.RenderResults(err.Error())
		return err
	}
	node := t.state.data.lastTest
	if node == nil {
		t.state.ui.result.RenderResults("failed to run last test. Make sure you run a test before rerunning")
		log.Error("attempted test rerun, but no test has previously been run", nil)
		return fmt.Errorf("no previously run test")
	}

	output, err := node.RunTest(ctx)
	if err != nil {
		log.Error("failed to re run valid test", err)
		t.state.ui.console.panel.UpdateMeta(t, output)
		t.state.ui.result.RenderResults(err.Error())
		return err
	}
	t.state.ui.result.RenderResults(output.Result)
	t.state.ui.console.panel.UpdateMeta(t, output)

	return nil
}
