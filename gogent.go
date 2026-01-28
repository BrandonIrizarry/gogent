package gogent

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

// Init initializes state used by the LLM, providing both the values
// of the Gogent struct's own fields, as well as setting up any state
// the LLM client needs to persist across prompt/response cycles. It
// returns a function that clients can use to initiate a single
// prompt/response cycle, likely in the context of some kind of REPL.
func (g *Gogent) Init() (askerFn, error) {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.TimeOnly,
	}).With().Caller().Logger()

	// Initialize any state needed by the function call objects
	// themselves.
	if err := functions.Init(g.WorkingDir, g.MaxFilesize); err != nil {
		return nil, err
	}

	// Initialize the LLM configuration.
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	contentConfig := genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: functions.FunctionDeclarations()},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	// The msgbuf slice has to be declared outside the definition
	// of the asker function, since this is where the context
	// across various invocations to the asker (as provided by,
	// say, a REPL in the frontend) will persist: it's the LLM's
	// "memory" of the conversation.
	msgbuf := []*genai.Content{}

	// This is the actual code that processes the user prompt. It
	// is the askerFn return value of Init.
	//
	// Note that this function captures many of the configuration
	// parameters defined just above (like 'client'.)
	asker := func(prompt string) (string, []TokenCount, error) {
		// The tokenCounts slice is where the token counts for
		// the current question/answer cycle are stored.
		var tokenCounts []TokenCount

		// Add the current user-prompt (the "question" being
		// asked, to complete the metaphor) to the message
		// buffer.
		//
		// This should be outside the function-call loop,
		// since it functions as the obvious preamble to the
		// LLM's subsequent thinking iterations..
		msgbuf = append(msgbuf, genai.NewContentFromText(prompt, genai.RoleUser))

		for i := range g.MaxIterations {
			log.Info().Msgf("Function-call iteration %d", i)

			response, err := client.Models.GenerateContent(
				ctx,
				g.LLMModel,
				msgbuf,
				&contentConfig,
			)
			if err != nil {
				msg := err.Error()

				// If we've hit a RESOURCE_EXHAUSTED
				// error, don't signal a quit to the
				// client; simply return the error
				// text as a valid response. In case
				// of any other error, treat it as an
				// actual show-stopping error and
				// return accordingly.
				if strings.HasPrefix(msg, "Error 429") {
					return msg, nil, nil
				} else {
					return "", nil, err
				}
			}

			tokenCounts = append(tokenCounts, convert(response.UsageMetadata))

			// Add the candidates to the message buffer. This
			// conforms both to the Gemini documentation, as well
			// as the Boot.dev AI Agent project.
			for _, candidate := range response.Candidates {
				msgbuf = append(msgbuf, candidate.Content)
			}

			// Check if the LLM has proposed any function calls to
			// act upon.
			funCalls := response.FunctionCalls()

			// The LLM is ready to give a textual response.
			if len(funCalls) == 0 {
				return response.Text(), tokenCounts, nil
			}

			for _, funCall := range funCalls {
				log.Trace().
					Any("args", funCall.Args).
					Msgf("Function call: %s", funCall.Name)

				// Handle the function call. If the
				// named function call isn't among our
				// declared functions, report an error
				// to the LLM. Else call the function
				// with the given arguments.
				var funCallResponsePart *genai.Part
				fnObj, err := functions.FunctionObject(funCall.Name)
				if err != nil {
					funCallResponsePart = functions.ResponseError(fnObj, err)
				} else {
					funCallResponsePart = fnObj.Function()(funCall.Args)
				}

				msgbuf = append(msgbuf, genai.NewContentFromParts([]*genai.Part{funCallResponsePart}, genai.RoleModel))
			}
		}

		return "", tokenCounts, errors.New("LLM didn't generate a text response")
	}

	return asker, nil
}

// convert returns a digest (a TokenCount struct) of relevant data
// inside the metadata argument.
func convert(metadata *genai.GenerateContentResponseUsageMetadata) TokenCount {
	return TokenCount{
		Prompt:     metadata.PromptTokenCount,
		Thoughts:   metadata.ThoughtsTokenCount,
		Cached:     metadata.CachedContentTokenCount,
		Candidates: metadata.CandidatesTokenCount,
		ToolUse:    metadata.ToolUsePromptTokenCount,
		Total:      metadata.TotalTokenCount,
	}
}
