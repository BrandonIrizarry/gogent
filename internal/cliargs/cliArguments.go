package cliargs

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type CLIArguments struct {
	NumIterations int
	WorkingDir    string
	MaxFilesize   int
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments
	var wdir string

	flag.IntVar(&cliArgs.NumIterations, "num", 20, "The number of times the function call loop should execute (defaults to 20)")
	flag.IntVar(&cliArgs.MaxFilesize, "maxsize", 10000, "The size limit for file-reading operations (defaults to 10KB)")
	flag.StringVar(&wdir, "dir", ".", "The top-level project directory (absolute path, else defaults to current directory)")

	flag.Parse()

	// Make sure we're using an absolute path, in case a relative one was given.
	wdir, err := filepath.Abs(wdir)

	if err != nil {
		return CLIArguments{}, err
	}

	cliArgs.WorkingDir = wdir

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
