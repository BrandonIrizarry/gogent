package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// fileContent returns the contents of path, which is absolute.
func fileContent(path string, maxFilesize int) (string, error) {
	fileBuf := make([]byte, maxFilesize)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Read(fileBuf)
	if err != nil {
		return "", fmt.Errorf("fileContent: %w", err)
	}

	return string(fileBuf), nil
}

// ignoredFilesMap returns the set of filenames ignored by the project
// per the project's .gitignore file. Each filename is an absolute
// path.
func ignoredFilesMap(workingDir string) (map[string]bool, error) {
	_, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "clean", "-ndX")

	// Aim the script at the project's working directory.
	var outputIgnored strings.Builder

	cmd.Dir = workingDir
	cmd.Stdout = &outputIgnored

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error calling gitignore script: %w", err)
	}

	// Initialize the 'entries' map.
	//
	// Include the .git directory manually, since git ls-files
	// doesn't list it.
	gitDir := filepath.Join(workingDir, ".git")
	ignored := map[string]bool{
		gitDir: true,
	}

	for e := range strings.SplitSeq(outputIgnored.String(), "\n") {
		// Splitting creates empty-string entries, which later
		// get confused as referring to the top-level
		// directory. So we must skip them here.
		if e == "" {
			continue
		}

		// Remove the 'Would remove' prefix.
		e = strings.TrimPrefix(e, "Would remove ")

		ignored[filepath.Join(workingDir, e)] = true
	}

	return ignored, nil
}

// pathIsIgnored returns whether a file is untracked.
func pathIsIgnored(ignoredPaths map[string]bool, path string) bool {
	_, ok := ignoredPaths[path]
	if ok {
		return true
	}

	for p := range ignoredPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}
