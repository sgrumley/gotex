package main

import (
	"sgrumley/gotex/internal/cli"

	"sgrumley/gotex/internal/finder"
)

func main() {

	project := finder.InitProject()
	// project.PrettyPrint()

	cli.Run(project.TestNameOut())
}
