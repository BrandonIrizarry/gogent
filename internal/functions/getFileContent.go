package functions

import (
	"log"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
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
	return func(args map[string]any, baseCfg baseconfig.BaseConfig) *genai.Part {
		path, err := normalizePath(args["filepath"], baseCfg.WorkingDir)
		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		content, logs, err := fileContent(path, baseCfg.MaxFilesize)
		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		log.Println(logs)

		return responseOK(fnobj.Name(), content)
	}
}
