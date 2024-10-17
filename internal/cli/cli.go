package cli

import (
	"fmt"
	"os"
	"os/signal"
	"sgrumley/gotex/internal/config"
	"sgrumley/gotex/internal/runner"
	"syscall"

	fzf "github.com/junegunn/fzf/src"
)

func Run(tests []string, locationMapping map[string]string) {
	// cfg, err := config.GetConfig("./internal/config/example.yaml")
	cfg, err := config.GetConfig("")
	if err != nil {
		fmt.Println(err)
	}

	inputChan := make(chan string)
	go func() {
		for _, s := range tests {
			inputChan <- s
		}
	}()

	wait := make(chan bool)
	outputChan := make(chan string)
	go func() {
		// WARN: should this be range? should it ever run more than once
		for s := range outputChan {
			location, exists := locationMapping[s]
			if !exists {
				fmt.Printf("failed to match test with file\n")
				wait <- true
				return
			}

			testOutput, err := runner.RunTest(s, location, cfg)
			if err != nil {
				fmt.Printf("failed to execute test: %s\n", err.Error())
				wait <- true
				return
			}

			// TODO: need to do something with the output?
			fmt.Println("\n" + testOutput)
			wait <- true
		}
	}()

	exit := func(code int, err error) {
		close(inputChan)
		close(outputChan)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(code)
	}

	// get signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-sigc
		wait <- true
	}()

	// search options
	options, err := fzf.ParseOptions(
		true, // load defaults
		[]string{
			"--multi",
			"--reverse",
			"--border",
			"--height=40%",
			"--bind=ctrl-j:down,ctrl-k:up",
		},
	)
	if err != nil {
		fmt.Println("failed parsing options")
		exit(fzf.ExitError, err)
	}

	options.Input = inputChan
	options.Output = outputChan

	code, err := fzf.Run(options)

	<-wait
	exit(code, err)
}