package functions

import (
	"fmt"
	"path/filepath"
)

// normalizePath takes the given path arg relative to workingDir, and
// returns the corresponding full path. The argument workingDir is
// assumed to already be an absolute path.
func normalizePath(arg any, workingDir string) (string, error) {
	path, ok := arg.(string)

	if !ok {
		return "", fmt.Errorf("Couldn't normalize '%v'", arg)
	}

	absPath := filepath.Join(workingDir, path)

	return absPath, nil
}
