package runner

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/sgrumley/gotex/pkg/config"
	"github.com/sgrumley/gotex/pkg/slogger"
)

type TestType int

const (
	TestTypeProject TestType = iota
	TestTypePackage
	TestTypeFile
	TestTypeFunction
	TestTypeCase
)

type Response struct {
	TestName        string
	TestType        TestType
	TestDir         string
	CommandExecuted string
	Result          string
	Output          string
	Error           string
	ExitStatus      int
	External        bool
	ExternalOutput  string
	ExternalError   string
}

func RunTest(ctx context.Context, testType TestType, testName string, dir string, cfg config.Config) (*Response, error) {
	log, err := slogger.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	cmdStr := GetCommand(testType, testName)
	cmdStr = applyConfig(cfg, cmdStr)

	if cfg.PipeTo != "" {
		res, err := RunTestPiped(cmdStr, cfg.PipeTo, dir)
		if err != nil {
			res.TestName = testName
			res.TestType = testType
			res.TestDir = dir
			res.CommandExecuted = argsToString(cmdStr) + " | " + cfg.PipeTo
			log.Error("failed to run piped command", err)
			return res, err
		}
		res.TestName = testName
		res.TestType = testType
		res.TestDir = dir
		res.CommandExecuted = argsToString(cmdStr) + " | " + cfg.PipeTo
		return res, nil
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd := exec.Command("go", cmdStr...)
	cmd.Dir = dir
	cmd.Stdout = buf
	cmd.Stderr = errBuf

	res := &Response{
		TestName:        testName,
		TestType:        testType,
		CommandExecuted: argsToString(cmdStr),
		TestDir:         dir,
	}

	log.Info("os command executed",
		slog.Any("args", cmd.Args),
		slog.String("dir", cmd.Dir),
		slog.Any("type", testType),
	)

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				// This should mean that the tests failed (not an actual error)
				res.Output = buf.String()
				res.Error = errBuf.String()
				res.ExitStatus = exitErr.ExitCode()
				return res, nil
			}
		} else {
			// Actual execution error e.g. command not found
			res.Output = buf.String()
			res.Error = errBuf.String()
			return res, err
		}
	}

	log.Info("test run successfully", slog.String("name", testName))
	res.Output = buf.String()
	res.Result = buf.String()
	res.Error = errBuf.String()
	res.ExitStatus = 0
	return res, nil
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

func GetCommand(typed TestType, testName string) []string {
	switch typed {
	case TestTypeProject:
		return []string{"test", "./..."}
	case TestTypePackage:
		return []string{"test"}
	case TestTypeFile, TestTypeFunction, TestTypeCase:
		regName := fmt.Sprintf("^%s$", testName)
		return []string{"test", "-run", regName}
	default:
		return []string{}
	}
}

func RunTestPiped(goTestCmdStr []string, externalCmdStr string, dir string) (*Response, error) {
	var goTestOutput bytes.Buffer
	var goTestErrBuf bytes.Buffer

	goTestCmd := exec.Command("go", goTestCmdStr...)
	goTestCmd.Stdout = &goTestOutput
	goTestCmd.Stderr = &goTestErrBuf
	goTestCmd.Dir = dir
	res := &Response{
		External: true,
	}

	if err := goTestCmd.Start(); err != nil {
		res.Output = goTestOutput.String()
		res.Error = goTestErrBuf.String()
		return res, fmt.Errorf("failed to start command 1: %w", err)
	}

	if err := goTestCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				// if it is an exit status code 1, continue to pipe the failed tests
				res.Result = goTestErrBuf.String()
				res.Output = goTestOutput.String()
				res.Error = goTestErrBuf.String()
				res.ExitStatus = exitErr.ExitCode()
				return res, err
			}
		}
	}

	// piping to external program
	r, w, err := os.Pipe()
	if err != nil {
		return res, fmt.Errorf("failed to create pipe: %w", err)
	}
	defer r.Close()

	go func() {
		_, _ = w.Write(goTestOutput.Bytes())
		w.Close()
	}()

	var externalCmdOutput bytes.Buffer
	var externalErrBuf bytes.Buffer
	externalCmd := exec.Command(externalCmdStr)
	externalCmd.Stdin = r
	externalCmd.Stdout = &externalCmdOutput
	externalCmd.Stderr = &externalErrBuf

	if err := externalCmd.Run(); err != nil {
		res.Output = goTestOutput.String()
		res.Error = goTestErrBuf.String()
		res.ExternalOutput = externalCmdOutput.String()
		res.ExternalError = externalErrBuf.String()
		res.Result = externalCmdOutput.String()
		return res, err
	}
	res.ExternalOutput = externalCmdOutput.String()
	res.ExternalError = externalErrBuf.String()
	res.Result = externalCmdOutput.String()

	return res, nil
}

func argsToString(args []string) string {
	return "go " + strings.Join(args, " ")
}
