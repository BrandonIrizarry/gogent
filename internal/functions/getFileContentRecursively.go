package functions

import (
	"fmt"
	"log/slog"
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
		all, err := allFilesMap(g.workingDir, dir)
		if err != nil {
			return g.ResponseError(err)
		}

		var bld strings.Builder
		for path := range all {
			content, err := fileContent(path, g.maxFilesize)
			if err != nil {
				return g.ResponseError(err)
			}

			fmt.Fprintf(&bld, "Contents of %s: %s\n\n", path, content)
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
