package functions

import (
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func (g getFileContentRecursively) Name() string {
	return "getFileContentRecursively"
}

func (g getFileContentRecursively) Function() functionType {
	// This callback reads contents of all files under a
	// given directory. A depth parameter must be specified.
	return func(args map[string]any) *genai.Part {
		dir, err := normalizePath(args["dir"], g.workingDir)
		if err != nil {
			return ResponseError(g.Name(), err.Error())
		}

		all, err := allFilesMap(g.workingDir, dir)
		if err != nil {
			return ResponseError(g.Name(), err.Error())
		}

		var bld strings.Builder
		for path := range all {
			content, err := fileContent(path, g.maxFilesize)
			if err != nil {
				return ResponseError(g.Name(), err.Error())
			}

			fmt.Fprintf(&bld, "Contents of %s: %s\n\n", path, content)
		}

		return responseOK(g.Name(), bld.String())
	}
}
