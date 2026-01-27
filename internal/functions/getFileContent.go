package functions

import (
	"google.golang.org/genai"
)

// Name implements [functionObject].
func (g getFileContent) Name() string {
	return "getFileContent"
}

// Function implements [functionObject].
func (g getFileContent) Function() functionType {
	// This callback reads the contents of the relative filepath
	// mentioned in args, and returns the corresponding Part object. If
	// there was an error in the internal logic, a Part corresponding to
	// an error is returned.
	return func(args map[string]any) *genai.Part {
		path, err := canonicalize(args[PropertyPath], g.workingDir)
		if err != nil {
			return ResponseError(g, err)
		}

		content, err := fileContent(path, g.maxFilesize)
		if err != nil {
			return ResponseError(g, err)
		}

		return ResponseOK(g, content)
	}
}
