package functions

import (
	"os"

	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"google.golang.org/genai"
)

type getFileContentType struct{}

// getFileContent reads the contents of the relative filepath
// mentioned in args, and returns the corresponding Part object. If
// there was an error in the internal logic, a Part corresponding to
// an error is returned.
var getFileContent getFileContentType

func (fnobj getFileContentType) Name() string {
	return "getFileContent"
}

func (fnobj getFileContentType) Function() functionType {
	return func(args map[string]any, cliArgs cliargs.CLIArguments) *genai.Part {
		absPath, err := normalizePath(args["filepath"])

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		dat, err := os.ReadFile(absPath)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		fileContents := string(dat)

		return responseOK(fnobj.Name(), fileContents)
	}
}
