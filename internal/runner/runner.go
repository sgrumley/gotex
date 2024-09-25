package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// TODO: execute go run and parse output
// NOTE: consider if this should be a method on project??
// NOTE: good resource: https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/
// TODO: check if go test tool can be imported or terminal exec - https://pkg.go.dev/testing#InternalExample
// unsure where to go with this
func RunTest() {
	testCMD := []string{"test"}
	testArgs := []string{"-json"}
	testList := "./..."
	cmdStr := append(testCMD, testArgs...)
	cmdStr = append(cmdStr, testList)
	cmd := exec.Command("go", cmdStr...)
	buf := new(bytes.Buffer)

	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	jsonobj := make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), jsonobj)
	fmt.Println(jsonobj)
}

/*
  permutations needed:
  go test -run ''        # Run all tests.
  go test -run Foo       # Run top-level tests matching "Foo", such as "TestFooBar".
go test -run Foo/A=    # For top-level tests matching "Foo", run subtests matching "A=".
go test -run /A=1      # For all top-level tests, run subtests matching "A=1".
go test -fuzz FuzzFoo  # Fuzz the target matching "FuzzFoo"
*/

func Commands() []*exec.Cmd {
	cmds := []*exec.Cmd{
		{
			// example params
			Path:         "",
			Args:         []string{},
			Env:          []string{},
			Dir:          "",
			Stdin:        nil,
			Stdout:       nil,
			Stderr:       nil,
			ExtraFiles:   []*os.File{},
			SysProcAttr:  &syscall.SysProcAttr{},
			Process:      &os.Process{},
			ProcessState: &os.ProcessState{},
			Err:          nil,
			Cancel: func() error {
				return nil
			},
			WaitDelay: 0,
		},
		// a single test case
		{
			Path: "project root",
			Args: []string{"-run"},
		},
		// a package
		{
			Path: "project root",
		},
		// a file

		// a function

	}

	return cmds
}
