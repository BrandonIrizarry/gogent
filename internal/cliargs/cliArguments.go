package cliargs

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CLIArguments struct {
	NumIterations int
	WorkingDir    string
	MaxFilesize   int
}

var ErrBadDefaultDir = errors.New("user didn't accept -dir default of current working directory")

func confirmUseCurrentDirectory() bool {
	// Warn the user that the current directory will be used if
	// that setting is detected.
	fmt.Println("No -dir argument was specified, so I'll default to the current directory.")

	for {
		fmt.Print("Is this OK? [Y/n] ")

		var confirm string
		fmt.Scanln(&confirm)

		confirm = strings.ToLower(confirm)

		switch {
		case confirm == "n":
			return false

		case len(confirm) == 0 || confirm == "y":
			return true

		default:
			fmt.Println("Please answer y/Y or n/N.")
			continue
		}
	}
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments
	var wdir string

	flag.IntVar(&cliArgs.NumIterations, "num", 20, "The number of times the function call loop should execute (defaults to 20)")
	flag.IntVar(&cliArgs.MaxFilesize, "maxsize", 200_000, "The size limit for file-reading operations (defaults to 200KB)")

	// Use the empty string as the default for -dir to later vet
	// whether using the current directory is acceptable. This is
	// to differentiate from the case where -dir is explicitly
	// provided as the current directory, in which case nothing
	// need be flagged to the user.
	flag.StringVar(&wdir, "dir", "", "The top-level project directory (absolute path, else defaults to current directory)")

	flag.Parse()

	// Check if the default argument for -dir is acceptable.
	if wdir == "" {
		if ok := confirmUseCurrentDirectory(); !ok {
			return CLIArguments{}, ErrBadDefaultDir
		}

		wdir = "."
	}

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
