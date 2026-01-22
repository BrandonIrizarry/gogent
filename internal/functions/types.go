package functions

import "google.golang.org/genai"

// functionType is the mandatory signature of each user-defined
// function. The sole argument is the map of LLM-provided arguments. A
// Part object is returned for consumption by the LLM. If an error is
// encountered, a Part encapsulating the error is returned.
type functionType func(map[string]any) *genai.Part

// functionObject groups all LLM functions under a common
// interface.
type functionObject interface {
	// Name returns the name of the function doing the dirty
	// work. Originally, this used to be the name of an actual Go
	// function. This method is really more of a cheat to require
	// the moral equivalent of a mandatory "name" field.
	Name() string

	// Function returns the function that does the dirty work. The client
	// using an implementor of functionObject should fetch this
	// return value and then invoke it with the necessary
	// arguments.
	Function() functionType
}

type functionDeclarationConfig struct {
	declaration genai.FunctionDeclaration
	fnObj       functionObject
}

// Types corresponding to LLM function calls.
type getFileContent struct {
	workingDir  string
	maxFilesize int
}

type getFileContentRecursively struct {
	workingDir  string
	maxFilesize int
}

type listDirectory struct {
	workingDir string
}
