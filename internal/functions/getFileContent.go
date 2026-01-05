package functions

import (
	"os"

	"google.golang.org/genai"
)

// getFileContent reads the contents of the relative filepath
// mentioned in args, and returns the corresponding Part object. If
// there was an error in the internal logic, a Part corresponding to
// an error is returned.
func getFileContent(args map[string]any) *genai.Part {
	path := args["filepath"].(string)
	dat, err := os.ReadFile(path)

	if err != nil {
		return ResponseError("getFileContent", err.Error())
	}

	fileContents := string(dat)

	return responseOK("getFileContent", fileContents)
}
