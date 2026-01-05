package functions

import "google.golang.org/genai"

func ResponseError(funCallName, message string) *genai.Part {
	return genai.NewPartFromFunctionResponse(funCallName, map[string]any{
		"error": message,
	})
}
