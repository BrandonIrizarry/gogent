package functions

import "google.golang.org/genai"

func responseOK(funCallName, content string) *genai.Part {
	return genai.NewPartFromFunctionResponse(funCallName, map[string]any{
		"result": content,
	})
}
