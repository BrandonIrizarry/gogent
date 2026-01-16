package cliargs

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

type CLIArguments struct {
	LogMode        logger.LogMode
	ConfigFilename string
	LogFilename    string
	WorkingDir     string
}

func NewCLIArguments() (CLIArguments, error) {
	var err error
	var cliArgs CLIArguments
	flag.Var(&cliArgs.LogMode, "logmode", `Comma-separated list of log-message types to output to logfile.

Examples:
-logmode debug
-logmode info,debug
`)
	flag.StringVar(&cliArgs.ConfigFilename, "config", "gogent.yaml", "YAML configuration file")
	flag.StringVar(&cliArgs.LogFilename, "log", "logs.txt", "Path to logfile")

	// Handle the 'dir' argument.
	//
	// FIXME: use the home directory to confirm that cwd is in
	// fact a descendant of the user's home directory.
	hdir, err := os.UserHomeDir()
	if err != nil {
		return CLIArguments{}, err
	}

	flag.StringVar(&cliArgs.WorkingDir, "dir", ".", "Path to working directory")
	flag.Parse()

	cliArgs.WorkingDir, err = filepath.Abs(cliArgs.WorkingDir)
	if err != nil {
		return CLIArguments{}, err
	}

	if !strings.HasPrefix(cliArgs.WorkingDir, hdir) {
		return CLIArguments{}, errors.New("Working directory not inside user's $HOME")
	}

	info, err := os.Stat(cliArgs.WorkingDir)
	if err != nil {
		return CLIArguments{}, err
	}

	if !info.IsDir() {
		return CLIArguments{}, errors.New("Working directory argument not a directory")
	}

	return cliArgs, nil
}
