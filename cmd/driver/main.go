package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/sgrumley/gotex/pkg/driver"
)

func main() {
	p, err := driver.InitProject(slog.Default())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dir: ", p.RootDir)
	p.Tree.Print()
}
