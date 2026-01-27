package gogent

import (
	"fmt"
	"path/filepath"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

func handleFunCall(funCall *genai.FunctionCall, workingDir string) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		errUnknownFunction := fmt.Errorf("unknown function: %s", funCall.Name)
		return functions.ResponseError(fnObj, errUnknownFunction)
	}

	fn := fnObj.Function()

	// Canonicalize the "path" argument if present.
	pathArg, ok := funCall.Args[functions.PropertyPath]
	if ok {
		path, ok := pathArg.(string)
		if !ok {
			err := fmt.Errorf("path arg not found among Args: %v", funCall.Args)
			return functions.ResponseError(fnObj, err)
		}
		funCall.Args[functions.PropertyPath] = filepath.Join(workingDir, path)
	}

	log.Trace().
		Any("path", funCall.Args[functions.PropertyPath]).
		Send()

	return fn(funCall.Args)
}
