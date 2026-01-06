package functions

import (
	"fmt"
	"path/filepath"
)

func normalizePath(arg any, workingDir string) (string, error) {
	path, ok := arg.(string)

	if !ok {
		return "", fmt.Errorf("Couldn't normalize '%v'", arg)
	}

	absPath := filepath.Join(workingDir, path)

	return absPath, nil
}
