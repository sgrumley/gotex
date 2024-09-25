package main

import (
	"log"
	"sgrumley/test-tui/internal/finder"
)

func main() {
	testpath := "./testdata/case_test.go"
	_, err := finder.FindSubTests(testpath)
	if err != nil {
		log.Fatalf("failed finding tests: %s", err)
	}
}
