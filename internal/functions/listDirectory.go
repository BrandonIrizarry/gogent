package functions

import (
	"fmt"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"google.golang.org/genai"
)

type listDirectoryType struct{}

var listDirectory listDirectoryType

func (fnobj listDirectoryType) Name() string {
	return "listDirectory"
}

func (fnobj listDirectoryType) Function() functionType {
	return func(args map[string]any, baseCfg baseconfig.BaseConfig) *genai.Part {
		dir, err := normalizePath(args["dir"], baseCfg.WorkingDir)

		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

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
