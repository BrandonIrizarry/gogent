package functions

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// ignoredFilesMap returns the set of filenames ignored by the project
// per the project's .gitignore file. Each filename is an absolute
// path.
func ignoredFilesMap(workingDir string) (map[string]bool, error) {
	var bld strings.Builder

	cmd := exec.Command("./ignored.sh")

	// Aim the script at the project's working directory (not
	// Gogent's working directory), and send the script's output
	// to our string builder.
	cmd.Dir = workingDir
	cmd.Stdout = &bld

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Initialize the 'entries' map.
	//
	// Include the .git directory manually, since git ls-files
	// doesn't list it.
	gitDir := filepath.Join(workingDir, ".git")
	entries := map[string]bool{
		gitDir: true,
	}

	for e := range strings.SplitSeq(bld.String(), "\n") {
		// Splitting creates empty-string entries, which later
		// get confused as referring to the top-level
		// directory. So we must skip them here.
		//
		// FIXME: defend against this where it gets
		// seen by other functions.
		if e == "" {
			continue
		}

		ne := filepath.Join(workingDir, e)
		entries[ne] = true
	}

	return entries, nil
}

// allFilesMap walks the filesystem starting at dir (an absolute path)
// and returns a set of absolute pathnames corresponding to files
// underneath dir. This function uses ignoreFilesMap to avoid walking
// down certain directories.
func allFilesMap(workingDir, path string) (map[string]bool, error) {
	ignored, err := ignoredFilesMap(workingDir)
	if err != nil {
		return nil, err
	}

	slog.Debug("After getting ignored files:", slog.Any("ignored", ignored))

	allFiles := make(map[string]bool)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		// It's a good idea to check 'd' for a nil value,
		// since it's possible that, for example, 'dir' was
		// malformed by some previous code and therefore the
		// current 'path' argument refers to a file or
		// directory that doesn't exist.
		if d == nil {
			return fmt.Errorf("nil direntry object for path %s", path)
		}

		slog.Debug("Current path:", slog.String("path", path))

		_, parentIsIgnored := ignored[filepath.Dir(path)]
		if parentIsIgnored {
			slog.Debug("Skipping because parent is ignored:", slog.String("path", path))
			return filepath.SkipDir
		}

		// FIXME: for now we check for "regular files", though
		// I'm not 100% sure this is what will always be
		// sufficient.
		if d.Type().IsRegular() {
			_, fileIsIgnored := ignored[path]
			if fileIsIgnored {
				slog.Debug("Skipping because file is ignored:", slog.String("path", path))
			} else {
				allFiles[path] = true
			}
		}

		return nil
	})

	slog.Debug("After getting all tracked files:", slog.Any("tracked", allFiles))

	return allFiles, nil

}
