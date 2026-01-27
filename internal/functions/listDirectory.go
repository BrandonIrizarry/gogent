package functions

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

func (fnobj listDirectory) Name() string {
	return "listDirectory"
}

func (fnobj listDirectory) Function() functionType {
	return func(args map[string]any) *genai.Part {
		log.Trace().
			Any("args", args).
			Msg("Inside listDirectory")

		dir, err := canonicalize(args[PropertyPath], fnobj.workingDir)
		if err != nil {
			return ResponseError(fnobj, err)
		}

		log.Trace().
			Str("canonicalized_path", dir).
			Send()

		files, err := os.ReadDir(dir)

		if err != nil {
			return ResponseError(fnobj, err)
		}

		bld := strings.Builder{}

		for _, file := range files {
			info, err := file.Info()

			if err != nil {
				return ResponseError(fnobj, err)
			}

			snippet := fmt.Sprintf("- %s: size=%d bytes, isDir: %v\n", info.Name(), info.Size(), info.IsDir())

			if _, err := bld.WriteString(snippet); err != nil {
				return ResponseError(fnobj, err)
			}
		}

		return ResponseOK(fnobj, bld.String())
	}
}
