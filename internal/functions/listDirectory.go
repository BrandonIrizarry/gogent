package functions

import (
	"fmt"
	"log/slog"
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
			return fnobj.ResponseError(err)
		}

		bld := strings.Builder{}

		for _, file := range files {
			info, err := file.Info()

			if err != nil {
				return fnobj.ResponseError(err)
			}

			snippet := fmt.Sprintf("- %s: size=%d bytes, isDir: %v\n", info.Name(), info.Size(), info.IsDir())

			if _, err := bld.WriteString(snippet); err != nil {
				return fnobj.ResponseError(err)
			}
		}

		return fnobj.ResponseOK(bld.String())
	}
}

func (fnobj listDirectory) ResponseError(err error) *genai.Part {
	message := err.Error()

	slog.Error("Response error:", slog.String("error", message))

	return genai.NewPartFromFunctionResponse(fnobj.Name(), map[string]any{
		"error": message,
	})
}

func (fnobj listDirectory) ResponseOK(content string) *genai.Part {
	return genai.NewPartFromFunctionResponse(fnobj.Name(), map[string]any{
		"result": content,
	})
}
