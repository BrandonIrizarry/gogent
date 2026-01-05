package functions

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func listDirectory(args map[string]any) *genai.Part {
	dir := args["dir"].(string)

	files, err := os.ReadDir(dir)

	if err != nil {
		return ResponseError("listDirectory", err.Error())
	}

	bld := strings.Builder{}

	for _, file := range files {
		info, err := file.Info()

		if err != nil {
			return ResponseError("listDirectory", err.Error())
		}

		snippet := fmt.Sprintf("- %s: size=%d bytes, isDir: %v\n", info.Name(), info.Size(), info.IsDir())

		if _, err := bld.WriteString(snippet); err != nil {
			return ResponseError("listDirectory", err.Error())
		}
	}

	return responseOK("listDirectory", bld.String())
}
