package functions

import (
	"fmt"
	"path/filepath"
)

func normalizePath(arg any) (string, error) {
	path, ok := arg.(string)

	if !ok {
		return "", fmt.Errorf("Couldn't normalize '%v'", arg)
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", fmt.Errorf("Couldn't get abspath: %w", err)
	}

	return absPath, nil
}
