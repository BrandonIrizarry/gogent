package gogent

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

type Gogent struct {
	WorkingDir    string
	MaxFilesize   int
	LLMModel      string
	MaxIterations int

	tokenCounts tokenCounts
}

func (g Gogent) TokenCounts() tokenCounts {
	return g.tokenCounts
}

// Init initializes state used by the LLM, providing both the values
// of the Gogent struct's own fields, as well as setting up any state
// the LLM client needs to persist across prompt/response cycles. It
// returns a function that clients can use to initiate a single
// prompt/response cycle, likely in the context of some kind of REPL.
func (g *Gogent) Init() (askerFn, error) {
	zerolog.TimeFieldFormat = time.TimeOnly

	log.Info().
		Str("working_dir", g.WorkingDir).
		Int("max_file_size", g.MaxFilesize).
		Str("llm_model", g.LLMModel).
		Int("max_iterations", g.MaxIterations).
		Msg("Gogent configuration")

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

	msgbuf := []*genai.Content{}

	// This is the actual code that processes the user prompt.
	//
	// Note that this function captures many of the configuration
	// parameters defined just above (like 'client'.)
	asker := func(prompt string) (string, error) {
		// Initialize the message buffer with the user
		// prompt. I'm making a careful note that this should
		// be outside the function-call loop.
		msgbuf = append(msgbuf, genai.NewContentFromText(prompt, genai.RoleUser))

		for i := range g.MaxIterations {
			log.Info().
				Int("iteration", i).
				Msg("Start of function-call loop")

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
					return msg, nil
				} else {
					return "", err
				}

			}

			g.incTokenCounts(response.UsageMetadata)

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
				text := response.Text()
				log.Info().Msg("Printing text response")

				return text, nil
			}

			for _, funCall := range funCalls {
				for arg, val := range funCall.Args {
					log.Trace().
						Str("name", funCall.Name).
						Str("arg", arg).
						Any("value", val).
						Msg("Function call:")
				}

				funCallResponsePart := handleFunCall(funCall, g.WorkingDir)
				msgbuf = append(msgbuf, genai.NewContentFromParts([]*genai.Part{funCallResponsePart}, genai.RoleModel))
			}
		}

		return "", errors.New("LLM didn't generate a text response")
	}

	return asker, nil
}
