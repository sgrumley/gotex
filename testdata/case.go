package testdata

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func AddNumbers(a, b int) int {
	fmt.Println("some text for terminal")
	if a == -10 {
		spew.Dump(a)
	}
	return a + b
}
