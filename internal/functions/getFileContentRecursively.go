package functions

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"google.golang.org/genai"
)

// Name implements [functionObject].
func (g getFileContentRecursively) Name() string {
	return "getFileContentRecursively"
}

// Function implements [functionObject].
func (g getFileContentRecursively) Function() functionType {
	// This callback reads contents of all files under a
	// given directory. A depth parameter must be specified.
	return func(args map[string]any) *genai.Part {
		dir, err := canonicalize(args[PropertyPath], g.workingDir)
		if err != nil {
			return ResponseError(g, err)
		}

		trackedPaths := []string{}

		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			// It's a good idea to check 'd' for a nil value,
			// since it's possible that, for example, 'dir' was
			// malformed by some previous code and therefore the
			// current 'path' argument refers to a file or
			// directory that doesn't exist.
			if d == nil {
				return fmt.Errorf("%s nonexistent (possibly malformed)", path)
			}

			// Skip the rest of the directory if the
			// parent directory is ignored.
			parent := filepath.Dir(path)
			if pathIsIgnored(ignoredPaths, parent) {
				return filepath.SkipDir
			}

			// Only do this if 'path' is actually a file.
			if d.Type().IsRegular() {
				if !pathIsIgnored(ignoredPaths, path) {
					trackedPaths = append(trackedPaths, path)
				}
			}

			return nil
		})

		var bld strings.Builder
		for _, tracked := range trackedPaths {
			content, err := fileContent(tracked, g.maxFilesize)
			if err != nil {
				return ResponseError(g, err)
			}

			fmt.Fprintf(&bld, "Contents of %s: %s\n\n", tracked, content)
		}

		return ResponseOK(g, bld.String())
	}
}
