package main

import (
	// "sgrumley/test-tui/internal/finder"
	"sgrumley/test-tui/internal/runner"
)

func main() {

	// project := finder.InitProject()
	// project.PrettyPrint()

	runner.RunTest()
	// config -> InitProject -> fzf -> runner

	// TODO:
	// - setup config
	//		- allow project param or cwd
	//		- configure defaults (tc.Name, prettyPrint, ...)
	// - use fzf (as lib) to allow searching tests (this will be the interface for this repo, TUIs and other plugins can import the libs)
	//		- when selecting a test it should run
	//		- cntrl j,k or up down
	//		- enter to run
	//- how nice can I make results look in just the terminal
}
