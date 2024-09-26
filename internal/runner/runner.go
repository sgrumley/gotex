package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"sgrumley/test-tui/internal/config"
)

// NOTE: good resource: https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/
// TODO: check if go test tool can be imported or terminal exec - https://pkg.go.dev/testing#InternalExample -> this would come with issues piping into other commands, at the least it wouldn't remove the need
func RunTest(testName string, dir string, cfg config.Config) (string, error) {
	buf := new(bytes.Buffer)
	cmdStr := GetCommand(1, testName)
	cmdStr = applyConfig(cfg, cmdStr)

	cmd := exec.Command("go", cmdStr...)
	cmd.Dir = dir
	cmd.Stdout = buf

	// TODO: clean at some point
	if cfg.PipeTo != "" {
		res, err := RunTestPiped(cmdStr, cfg.PipeTo, dir)
		if err == nil {
			return res.String(), nil
		}

		fmt.Println("failed piped command, running without it: ", err)
	}

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TODO: clean
func RunTestPiped(cmdStr1 []string, cmdStr2 string, dir string) (*bytes.Buffer, error) {

	// Buffer to capture the output of the first command
	var cmd1Output bytes.Buffer

	// Prepare the first command
	cmd1 := exec.Command("go", cmdStr1...)
	cmd1.Stdout = &cmd1Output // Capture stdout into the buffer
	var errBuf1 bytes.Buffer
	cmd1.Stderr = &errBuf1 // Capture stderr
	cmd1.Dir = dir

	// Start the first command
	if err := cmd1.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command 1: %w", err)
	}

	// Wait for the first command to finish
	if err := cmd1.Wait(); err != nil {
		return nil, fmt.Errorf("command 1 failed: %w, stderr: %s", err, errBuf1.String())
	}

	// leaving for debugging purposes
	// fmt.Println("Output of command 1:")
	// fmt.Println(cmd1Output.String())

	// Create a pipe for the second command
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}
	defer r.Close() // Ensure the read end is closed after function exit

	// Write command 1's output to the write end of the pipe
	go func() {
		_, _ = w.Write(cmd1Output.Bytes()) // Write output to the pipe
		w.Close()                          // Close the write end after writing
	}()

	// Buffer to capture the output of the second command
	buf := new(bytes.Buffer)

	// Prepare the second command
	cmd2 := exec.Command(cmdStr2)
	cmd2.Stdin = r    // Set the stdin to the read end of the pipe
	cmd2.Stdout = buf // Set the stdout to capture the output
	var errBuf2 bytes.Buffer
	cmd2.Stderr = &errBuf2 // Capture stderr

	// Run the second command
	if err := cmd2.Run(); err != nil {
		return nil, fmt.Errorf("failed to run command 2: %w, stderr: %s", err, errBuf2.String())
	}

	return buf, nil
	// temp
	// r, w, err := os.Pipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create pipe: %w", err)
	// }
	// defer r.Close() // Ensure the read end is closed after function exit

	// // Prepare the first command
	// fmt.Println("cmdStr1 ", cmdStr1)
	// cmd1 := exec.Command("go", cmdStr1...)
	// cmd1.Stdout = w // Set the stdout to the write end of the pipe
	// var errBuf1 bytes.Buffer
	// cmd1.Stderr = &errBuf1 // Capture stderr

	// // Start the first command
	// if err := cmd1.Start(); err != nil {
	// 	return nil, fmt.Errorf("failed to start command 1: %w", err)
	// }

	// // Close the write end after starting the first command
	// go func() {
	// 	if err := cmd1.Wait(); err != nil {
	// 		fmt.Printf("command 1 error: %s\n", err)
	// 	}
	// 	w.Close() // Close the write end after command 1 completes
	// }()

	// // Buffer to capture the output of the second command
	// buf := new(bytes.Buffer)

	// // Prepare the second command
	// cmd2 := exec.Command(cmdStr2)
	// cmd2.Stdin = r    // Set the stdin to the read end of the pipe
	// cmd2.Stdout = buf // Set the stdout to capture the output
	// var errBuf2 bytes.Buffer
	// cmd2.Stderr = &errBuf2 // Capture stderr

	// // Run the second command
	// if err := cmd2.Run(); err != nil {
	// 	return nil, fmt.Errorf("failed to run command 2: %w, stderr: %s", err, errBuf2.String())
	// }

	// // Check if there were any errors from command 1
	// if errBuf1.Len() > 0 {
	// 	return nil, fmt.Errorf("command 1 stderr: %s", errBuf1.String())
	// }

	// return buf, nil

	// ORIGINAL
	// r, w, err := os.Pipe()
	// if err != nil {
	// 	fmt.Println("returning 1, ", err)
	// 	return nil, err
	// }
	// defer r.Close()
	// cmd1 := exec.Command("go", cmdStr1...)
	// cmd1.Stdout = w
	// err = cmd1.Start()
	// if err != nil {

	// 	fmt.Println("returning 2, ", err)
	// 	return nil, err
	// }
	// defer cmd1.Wait()
	// w.Close()

	// buf := new(bytes.Buffer)
	// cmd2 := exec.Command(cmdStr2)
	// cmd2.Stdin = r
	// cmd2.Stdout = buf
	// err = cmd2.Run()
	// if err != nil {
	// 	fmt.Println("returning 3, ", err.Error(), " detail ", buf.String())
	// 	return nil, err
	// }

	// return buf, nil
}

func applyConfig(cfg config.Config, cmd []string) []string {
	// Start with the original command
	args := append([]string{}, cmd...) // Create a copy of the command

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
		return []string{"test", "-run", fmt.Sprintf(`"%s"`, testName)}
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
