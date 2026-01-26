package functions

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

// ResponseError accepts an error and returns it as a
// [*genai.Part] object consumable by the LLM.
func ResponseError(fnObj functionObject, err error) *genai.Part {
	message := err.Error()

	log.Error().
		Err(err).
		Str("function", fnObj.Name()).
		Msg("response error")

	return genai.NewPartFromFunctionResponse(fnObj.Name(), map[string]any{
		"error":    message,
		"function": fnObj.Name(),
	})
}

// ResponseOK accepts a string message representing the output
// of an LLM function call, and returns it as a [*genai.Part]
// object consumable by the LLM.
func ResponseOK(fnObj functionObject, content string) *genai.Part {
	log.Info().
		Str("function", fnObj.Name()).
		Msg("response OK")

	return genai.NewPartFromFunctionResponse(fnObj.Name(), map[string]any{
		"result":   content,
		"function": fnObj.Name(),
	})
}
