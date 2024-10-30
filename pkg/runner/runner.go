package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"sgrumley/gotex/pkg/config"
)

// TODO: no std -> write everything to a buf and return it
func RunTest(testName string, dir string, cfg config.Config) (string, error) {
	buf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmdStr := GetCommand(1, testName)
	cmdStr = applyConfig(cfg, cmdStr)

	cmd := exec.Command("go", cmdStr...)
	cmd.Dir = dir
	cmd.Stdout = buf
	cmd.Stderr = errBuf

	// return fmt.Sprintf("cmd:%v %s", cmd.Args, cmd.Dir), nil
	// if the pipe to value is found in config attempt to run both commands via RunTestPiped
	// if this fails it will run without the piped command
	if cfg.PipeTo != "" {
		res, err := RunTestPiped(cmdStr, cfg.PipeTo, dir)
		if err == nil {
			return res.String(), nil
		}

		// TODO: log to file
		// fmt.Println("failed piped command, running without it: ", err)
	}

	// TODO: log to file
	// fmt.Printf("running cmd: %v from dir: %s\n", cmd.Args, cmd.Dir)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		// TODO: go test returns error if tests fail. This needs a correct solution
		// some logic to determine if the output is exit 1. If so this does not mean an error within the command but could be that the test did not pass
		errStr := errBuf.String()
		if err.Error() != "exit status 1" {
			// fmt.Println("buf inside", buf.String())
			// fmt.Println("err inside", errStr)
			// fmt.Println("err ", err)
			return "", fmt.Errorf(errStr)
		}
	}

	return buf.String(), nil
}

// TODO: clean
func RunTestPiped(cmdStr1 []string, cmdStr2 string, dir string) (*bytes.Buffer, error) {
	var cmd1Output bytes.Buffer
	var errBuf1 bytes.Buffer

	cmd1 := exec.Command("go", cmdStr1...)
	cmd1.Stdout = &cmd1Output
	cmd1.Stderr = &errBuf1
	cmd1.Dir = dir

	if err := cmd1.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command 1: %w", err)
	}

	if err := cmd1.Wait(); err != nil {
		return nil, fmt.Errorf("go test command failed: %w, stderr: %s", err, errBuf1.String())
	}

	// leaving for debugging purposes
	// fmt.Println("Output of command 1:")
	// fmt.Println(cmd1Output.String())

	// Create a pipe for the second command
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}
	defer r.Close()

	go func() {
		_, _ = w.Write(cmd1Output.Bytes())
		w.Close()
	}()

	// Buffer to capture the output of the second command
	var cmd2Output bytes.Buffer
	var errBuf2 bytes.Buffer
	// Prepare the second command
	cmd2 := exec.Command(cmdStr2)
	cmd2.Stdin = r
	cmd2.Stdout = &cmd2Output
	cmd2.Stderr = &errBuf2

	// Run the second command
	if err := cmd2.Run(); err != nil {
		errStr := errBuf2.String()
		if errStr != "exit status 1" {
			fmt.Println("buf inside", cmd2Output.String())
			fmt.Println("err inside", errStr)

			return nil, fmt.Errorf("piped command failed: %w, stderr: %s", err, errBuf2.String())
		}
	}

	return &cmd2Output, nil
}

func applyConfig(cfg config.Config, cmd []string) []string {
	args := append([]string{}, cmd...)

	// Add options based on the config
	if cfg.Timeout != "" {
		args = append(args, "-timeout", cfg.Timeout)
	}
	if cfg.Json {
		args = append(args, "-json")
	}
	if cfg.Short {
		args = append(args, "-short")
	}
	if cfg.Verbose {
		args = append(args, "-v")
	}
	if cfg.FailFast {
		args = append(args, "-failfast")
	}
	if cfg.Cover {
		args = append(args, "-cover")
	}

	return args
}

func GetCommand(typed int, testName string) []string {
	switch typed {
	case 1:
		return []string{"test", "-run", testName}
	case 2:
		// TODO: test all
		return []string{"test", "./..."}
	case 3:
		// TODO: test single file
		return []string{"test", "path to file"}
	case 4:
		// TODO: test whole package
		return []string{"test", "package name"}
	default:
		return nil
	}
}
