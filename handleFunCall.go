package gogent

import (
	"fmt"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

func handleFunCall(funCall *genai.FunctionCall) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args)
}
