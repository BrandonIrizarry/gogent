package functions

import (
	"log/slog"

	"google.golang.org/genai"
)

// ResponseError accepts an error and returns it as a
// [*genai.Part] object consumable by the LLM.
func ResponseError(fnObj functionObject, err error) *genai.Part {
	message := err.Error()

	slog.Error("Response error:",
		slog.String("error", message),
		slog.String("function", fnObj.Name()),
	)

	return genai.NewPartFromFunctionResponse(fnObj.Name(), map[string]any{
		"error":    message,
		"function": fnObj.Name(),
	})
}

// ResponseOK accepts a string message representing the output
// of an LLM function call, and returns it as a [*genai.Part]
// object consumable by the LLM.
func ResponseOK(fnObj functionObject, content string) *genai.Part {
	slog.Info("Response OK:",
		slog.String("function", fnObj.Name()),
	)

	return genai.NewPartFromFunctionResponse(fnObj.Name(), map[string]any{
		"result":   content,
		"function": fnObj.Name(),
	})
}
