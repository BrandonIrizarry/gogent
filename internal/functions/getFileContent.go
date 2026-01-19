package functions

import (
	"google.golang.org/genai"
)

func (g getFileContent) Name() string {
	return "getFileContent"
}

func (g getFileContent) Function() functionType {
	// This callback reads the contents of the relative filepath
	// mentioned in args, and returns the corresponding Part object. If
	// there was an error in the internal logic, a Part corresponding to
	// an error is returned.
	return func(args map[string]any) *genai.Part {
		path, err := normalizePath(args["filepath"], g.workingDir)
		if err != nil {
			return ResponseError(g.Name(), err.Error())
		}

		content, err := fileContent(path, g.maxFilesize)
		if err != nil {
			return ResponseError(g.Name(), err.Error())
		}

		return responseOK(g.Name(), content)
	}
}
