package cliargs

import (
	"flag"
)

type CLIArguments struct {
	Verbose bool
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments
	flag.BoolVar(&cliArgs.Verbose, "verbose", false, "Whether to print usage metadata")
	flag.Parse()

	return cliArgs, nil
}
