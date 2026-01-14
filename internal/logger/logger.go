package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogMode int

// There is no "LogSettingError" because logging errors is always
// enabled.
const (
	LogModeInfo LogMode = 1 << iota
	LogModeDebug
)

var (
	Info  *log.Logger
	Debug *log.Logger

	// Don't export the error logger, since we only log errors
	// with [ReportError].
	errorLogger *log.Logger
)

// New returns a new Logger. The log file's lifetime is scoped by the
// main function and so must be passed as the logFile parameter
// here. The verbositySetting is a bitfield specifying which loggers
// to use.
func Init(logFile *os.File, verbositySetting LogMode) {
	// Use 'dest' in case it ever becomes feasible to log to
	// stderr as well.
	dest := logFile

	infoWriter := io.Discard
	debugWriter := io.Discard
	errorWriter := dest

	if satisfies(verbositySetting, LogModeInfo) {
		infoWriter = dest
	}

	if satisfies(verbositySetting, LogModeDebug) {
		debugWriter = dest
	}

	logFlags := log.Llongfile

	Info = log.New(infoWriter, "INFO: ", logFlags)
	Debug = log.New(debugWriter, "DEBUG: ", logFlags)
	errorLogger = log.New(errorWriter, "ERROR: ", logFlags)

}

func ReportError(err error, msg string) error {
	errorLogger.Output(2, msg)

	return fmt.Errorf("%s: %w", msg, err)
}

// Set implements the flag.Value interface. It's used by the cliargs
// package to obtain the log settings from the command line.
func (s *LogMode) Set(value string) error {
	for setting := range strings.SplitSeq(value, ",") {
		switch setting {
		case "info":
			*s |= LogModeInfo
		case "debug":
			*s |= LogModeDebug
		default:
			return fmt.Errorf("invalid log setting: %s", setting)
		}
	}

	return nil
}

// String implements the flag.Value interface (see [LogMode.Set] above.)
func (s *LogMode) String() string {
	switch *s {
	case LogModeInfo:
		return "info"
	case LogModeDebug:
		return "debug"
	default:
		return ""
	}
}
