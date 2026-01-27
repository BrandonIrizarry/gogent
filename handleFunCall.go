package gogent

import (
	"fmt"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

func handleFunCall(funCall *genai.FunctionCall) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		errUnknownFunction := fmt.Errorf("unknown function: %s", funCall.Name)
		return functions.ResponseError(fnObj, errUnknownFunction)
	}

	fn := fnObj.Function()
	return fn(funCall.Args)
}
