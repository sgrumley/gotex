package main

import (
	"sgrumley/test-tui/internal/finder"
)

func main() {

	project := finder.InitProject()
	project.PrettyPrint()
	// testpath := "./testdata/case_test.go"
	// _, err := finder.ListAll(testpath)
	// if err != nil {
	// 	log.Fatalf("failed finding tests: %s", err)
	// }

	// TODO: these functions may need new implementations since the existing functionality requires passing in ast types

	// _, err := finder.ListFunctions(testpath)
	// if err != nil {
	// 	log.Fatalf("failed finding test functions")
	// }

	// _, err := finder.ListTestCases(testpath, testFunction)
	// if err != nil {
	// 	log.Fatalf("failed finding test functions")
	// }
}
