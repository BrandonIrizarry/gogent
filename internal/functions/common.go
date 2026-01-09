package functions

import (
	"fmt"
	"os"
	"strings"
)

// fileContent returns the contents of filename (an absolute path) as
// a string, as well as a log message. If there is nothing to log, an
// empty string is returned.
func fileContent(path string, maxFilesize int) (string, string, error) {
	var logMsgs strings.Builder

	fileBuf := make([]byte, maxFilesize)

	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	numBytes, err := file.Read(fileBuf)
	if err != nil {
		return "", "", err
	}

	switch numBytes {
	case 0:
		fmt.Fprintf(&logMsgs, "Warning: read zero bytes from %s", path)
	case maxFilesize:
		fmt.Fprintf(&logMsgs, "Warning: read maximum number of bytes (%d) from %s; possible truncation", numBytes, path)
	default:
		fmt.Fprintf(&logMsgs, "OK: read %d bytes from %s", numBytes, path)
	}

	fileContents := string(fileBuf)
	logContents := logMsgs.String()

	return fileContents, logContents, nil
}
