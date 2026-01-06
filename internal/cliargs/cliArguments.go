package cliargs

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CLIArguments struct {
	NumIterations int
	ToplevelDir   string
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments

	flag.IntVar(&cliArgs.NumIterations, "num", 20, "The number of times the function call loop should execute")
	flag.StringVar(&cliArgs.ToplevelDir, "dir", "", "The top-level project directory")

	flag.Parse()

	// Make sure the '-dir' argument was given.
	if cliArgs.ToplevelDir == "" {
		return CLIArguments{}, errors.New("Missing '-dir' argument")
	}

	// Make sure the provided working directory is, in fact, a
	// directory.
	finfo, err := os.Stat(cliArgs.ToplevelDir)

	if err != nil {
		return CLIArguments{}, err
	}

	if ok := finfo.IsDir(); !ok {
		return CLIArguments{}, fmt.Errorf("not a directory: %s", cliArgs.ToplevelDir)
	}

	return cliArgs, nil
}
