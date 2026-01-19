package gogent

import (
	"fmt"
	"path/filepath"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

func handleFunCall(funCall *genai.FunctionCall, workingDir string) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	// Canonicalize the "path" argument if present.
	pathArg, ok := funCall.Args[functions.PropertyPath]
	if ok {
		path := pathArg.(string)
		funCall.Args[functions.PropertyPath] = filepath.Join(workingDir, path)
	}

	return fn(funCall.Args)
}
