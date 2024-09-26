package runner

import (
	"bytes"
	"os/exec"
)

// NOTE: good resource: https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/
// TODO: check if go test tool can be imported or terminal exec - https://pkg.go.dev/testing#InternalExample
func RunTest(testName string, dir string) (string, error) {
	// TODO: dynamic dir
	// TODO: get config
	buf := new(bytes.Buffer)

	cmdStr := GetCommand(1, testName)
	cmd := exec.Command("go", cmdStr...)
	cmd.Dir = dir
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetCommand(typed int, testName string) []string {
	switch typed {
	case 1:
		return []string{"test", "-run", testName, "-json"}
	// TODO: test whole package
	case 2:
		return []string{"test", "./..."}
	// TODO: test all
	case 3:
		return []string{"test", "./..."}
	// TODO: test single file
	case 4:
		return []string{"test", "./..."}
	default:
		return nil
	}
}
