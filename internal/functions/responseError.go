package functions

import (
	"log/slog"

	"google.golang.org/genai"
)

// ResponseError returns a Part describing an error in our (the
// programmer's) functions (for example, trying to read from a file
// that doesn't exist.) It accepts the name of the function call, as
// well as the error message to be included as part of the return
// value.
func ResponseError(funCallName, message string) *genai.Part {
	slog.Error("Response error:", slog.String("error", message))

	return genai.NewPartFromFunctionResponse(funCallName, map[string]any{
		"error": message,
	})
}
