package cliargs

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CLIArguments struct {
	NumIterations int
	WorkingDir    string
	MaxFilesize   int
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments

	flag.IntVar(&cliArgs.NumIterations, "num", 20, "The number of times the function call loop should execute (defaults to 20)")
	flag.IntVar(&cliArgs.MaxFilesize, "maxsize", 10000, "The size limit for file-reading operations (defaults to 10KB)")
	flag.StringVar(&cliArgs.WorkingDir, "dir", "", "The top-level project directory")

	flag.Parse()

	// Make sure the '-dir' argument was given.
	if cliArgs.WorkingDir == "" {
		return CLIArguments{}, errors.New("Missing '-dir' argument")
	}

	// Make sure the provided working directory is, in fact, a
	// directory.
	finfo, err := os.Stat(cliArgs.WorkingDir)

	if err != nil {
		return CLIArguments{}, err
	}

	if ok := finfo.IsDir(); !ok {
		return CLIArguments{}, fmt.Errorf("not a directory: %s", cliArgs.WorkingDir)
	}

	return cliArgs, nil
}
