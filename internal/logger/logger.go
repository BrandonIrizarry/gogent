package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Logger struct {
	Info    *log.Logger
	Verbose *log.Logger
}

var logFlags = log.LstdFlags | log.Lshortfile

// New returns a new Logger. The file pointer is passed as a
// parameter, since closing it has to be managed from the main
// function in main.go. The verbose flag enables clients to perform
// verbose logging outside of conditional blocks.
func New(logFile *os.File, verbose bool, packageName string) Logger {
	var lg Logger
	var verboseWriter io.Writer

	if verbose {
		verboseWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		verboseWriter = io.Discard
	}

	prefix := fmt.Sprintf("%s: ", strings.ToUpper(packageName))
	verbosePrefix := fmt.Sprintf("VERBOSE %s: ", strings.ToUpper(packageName))

	lg.Info = log.New(logFile, prefix, logFlags)
	lg.Verbose = log.New(verboseWriter, verbosePrefix, logFlags)

	return lg
}
