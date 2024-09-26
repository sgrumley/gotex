package main

import (
	"sgrumley/test-tui/internal/cli"

	"sgrumley/test-tui/internal/finder"
)

// "sgrumley/test-tui/internal/runner"
// "sgrumley/test-tui/internal/cli"

func main() {

	project := finder.InitProject()
	// project.PrettyPrint()

	cli.Run(project.TestNameOut())
}
