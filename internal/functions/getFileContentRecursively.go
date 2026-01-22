package functions

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"google.golang.org/genai"
)

func (g getFileContentRecursively) Name() string {
	return "getFileContentRecursively"
}

// FIXME: inline allFilesMap filesystem walking code into this
// callback, using pathIsIgnored instead.
func (g getFileContentRecursively) Function() functionType {
	// This callback reads contents of all files under a
	// given directory. A depth parameter must be specified.
	return func(args map[string]any) *genai.Part {
		dir := args[PropertyPath].(string)
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

			slog.Debug("Current path:", slog.String("path", path))

			// Skip the rest of the directory if the
			// parent directory is ignored.
			parent := filepath.Dir(path)
			if pathIsIgnored(ignoredPaths, parent) {
				slog.Debug("Skipping because parent is ignored:", slog.String("path", path))
				return filepath.SkipDir
			}

			if !pathIsIgnored(ignoredPaths, path) {
				trackedPaths = append(trackedPaths, path)
			}

			return nil
		})

		slog.Debug("After getting all tracked files:", slog.Any("tracked", trackedPaths))

		var bld strings.Builder
		for _, tracked := range trackedPaths {
			content, err := fileContent(tracked, g.maxFilesize)
			if err != nil {
				return g.ResponseError(err)
			}

			fmt.Fprintf(&bld, "Contents of %s: %s\n\n", tracked, content)
		}

		return g.ResponseOK(bld.String())
	}
}

func (g getFileContentRecursively) ResponseError(err error) *genai.Part {
	message := err.Error()

	slog.Error("Response error:", slog.String("error", message))

	return genai.NewPartFromFunctionResponse(g.Name(), map[string]any{
		"error": message,
	})
}

func (g getFileContentRecursively) ResponseOK(content string) *genai.Part {
	return genai.NewPartFromFunctionResponse(g.Name(), map[string]any{
		"result": content,
	})
}
