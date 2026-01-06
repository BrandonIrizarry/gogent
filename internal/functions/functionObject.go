package functions

import (
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"google.golang.org/genai"
)

// functionType is the mandatory signature of each user-defined
// function. The first argument is the map of LLM-provided arguments,
// while the second is user configuation. A Part object is returned
// for consumption by the LLM. If an error is encountered, a Part
// encapsulating the error is returned.
type functionType func(map[string]any, cliargs.CLIArguments) *genai.Part

// functionObject groups all our LLM functions under a common
// interface. Originally, the purpose of such an interface is to be
// able to refer to a function's name without the boilerplate of
// hard-coding it inside each individual function definition. However,
// more uses may arise as needed.
type functionObject interface {
	// Returns the name of the function doing the dirty work. Originally,
	// this used to be the name of an actual Go function.
	Name() string

	// Returns the function that does the dirty work. The client
	// using an implementor of functionObject should fetch this
	// return value and then invoke it with the necessary
	// arguments.
	Function() functionType
}
