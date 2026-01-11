package cliargs

import (
	"flag"
)

type CLIArguments struct {
	Verbose        bool
	ConfigFilename string
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments
	flag.BoolVar(&cliArgs.Verbose, "verbose", false, "Whether to print usage metadata")
	flag.StringVar(&cliArgs.ConfigFilename, "config", "gogent.yaml", "YAML configuration file (defaults to gogent.yaml)")
	flag.Parse()

	return cliArgs, nil
}
