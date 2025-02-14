package testdata

import "github.com/davecgh/go-spew/spew"

func AddNumbers(a, b int) int {
	if a == -10 {
		spew.Dump(a)
	}
	return a + b
}
