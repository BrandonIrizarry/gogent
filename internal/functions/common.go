package functions

import (
	"log/slog"
	"os"
)

// fileContent returns the contents of path.
func fileContent(path string, maxFilesize int) (string, error) {
	fileBuf := make([]byte, maxFilesize)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	numBytes, err := file.Read(fileBuf)
	if err != nil {
		return "", err
	}

	slog.Info("LLM file content read", slog.String("path", path), slog.Int("bytes", numBytes))

	return string(fileBuf), nil
}
