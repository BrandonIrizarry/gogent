package functions

import (
	"log"
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
		absPath, err := normalizePath(args["filepath"], cliArgs.ToplevelDir)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		fileBuf := make([]byte, cliArgs.MaxFilesize)

		file, err := os.Open(absPath)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		defer file.Close()

		numBytes, err := file.Read(fileBuf)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		switch numBytes {
		case 0:
			log.Printf("Warning: read zero bytes from %s", absPath)
		case cliArgs.MaxFilesize:
			log.Printf("Warning: read maximum number of bytes (%d) from %s; possible truncation", numBytes, absPath)
		default:
			log.Printf("OK: read %d bytes from %s", numBytes, absPath)
		}

		fileContents := string(fileBuf)

		return responseOK(fnobj.Name(), fileContents)
	}
}
