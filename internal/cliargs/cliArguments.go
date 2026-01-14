package cliargs

import (
	"flag"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

type CLIArguments struct {
	LogMode        logger.LogMode
	ConfigFilename string
	LogFilename    string
}

func NewCLIArguments() (CLIArguments, error) {
	var cliArgs CLIArguments
	flag.Var(&cliArgs.LogMode, "logmode", `Comma-separated list of log-message types to output to logfile.

Examples:
-logmode debug
-logmode info,debug
`)
	flag.StringVar(&cliArgs.ConfigFilename, "config", "gogent.yaml", "YAML configuration file (defaults to gogent.yaml)")
	flag.StringVar(&cliArgs.LogFilename, "log", "logs.txt", "Path to logfile (defaults to logs.txt)")

	flag.Parse()

	return cliArgs, nil
}
