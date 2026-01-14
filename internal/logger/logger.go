package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogSetting int

// There is no "LogSettingError" because logging errors is always
// enabled.
const (
	LogSettingInfo LogSetting = 1 << iota
	LogSettingDebug
)

type Logger struct {
	info  *log.Logger
	debug *log.Logger
	error *log.Logger
}

var logger Logger

// New returns a new Logger. The log file's lifetime is scoped by the
// main function and so must be passed as the logFile parameter
// here. The verbositySetting is a bitfield specifying which loggers
// to use.
func Init(logFile *os.File, verbositySetting LogSetting) {
	// Use 'dest' in case it ever becomes feasible to log to
	// stderr as well.
	dest := logFile

	infoWriter := io.Discard
	debugWriter := io.Discard
	errorWriter := dest

	if satisfies(verbositySetting, LogSettingInfo) {
		infoWriter = dest
	}

	if satisfies(verbositySetting, LogSettingDebug) {
		debugWriter = dest
	}

	logFlags := log.Llongfile
	logger = Logger{
		info:  log.New(infoWriter, "INFO: ", logFlags),
		debug: log.New(debugWriter, "DEBUG: ", logFlags),
		error: log.New(errorWriter, "ERROR: ", logFlags),
	}
}

func Info() *log.Logger {
	return logger.info
}

func Debug() *log.Logger {
	return logger.debug
}

func Error(err error, msg string) error {
	logger.error.Output(2, msg)

	return fmt.Errorf("%s: %w", msg, err)
}

func satisfies(verbosity, mask LogSetting) bool {
	return verbosity&mask == mask
}
