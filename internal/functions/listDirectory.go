package functions

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func (fnobj listDirectory) Name() string {
	return "listDirectory"
}

func (fnobj listDirectory) Function() functionType {
	return func(args map[string]any) *genai.Part {
		dir := args[PropertyPath].(string)
		files, err := os.ReadDir(dir)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		bld := strings.Builder{}

		for _, file := range files {
			info, err := file.Info()

			if err != nil {
				return ResponseError(fnobj.Name(), err.Error())
			}

			snippet := fmt.Sprintf("- %s: size=%d bytes, isDir: %v\n", info.Name(), info.Size(), info.IsDir())

			if _, err := bld.WriteString(snippet); err != nil {
				return ResponseError(fnobj.Name(), err.Error())
			}
		}

		return responseOK(fnobj.Name(), bld.String())
	}
}
