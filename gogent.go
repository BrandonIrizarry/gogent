package gogent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

type Gogent struct {
	WorkingDir    string
	MaxFilesize   int
	LLMModel      string
	MaxIterations int
	LogLevel      string

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
	// Set the appropriate logging level.
	logLevels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level, ok := logLevels[g.LogLevel]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %s", g.LogLevel)
	}

	slog.SetLogLoggerLevel(level)

	slog.Info("Gogent configuration:",
		slog.String("working_dir", g.WorkingDir),
		slog.Int("max_file_size", g.MaxFilesize),
		slog.String("llm_model", g.LLMModel),
		slog.Int("max_iterations", g.MaxIterations),
		slog.String("logging verbosity", g.LogLevel),
	)

	// Initialize any state needed by the function call objects
	// themselves.
	functions.Init(g.WorkingDir, g.MaxFilesize)

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
			slog.Info("Start of function-call loop:", slog.Int("iteration", i))

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
				slog.Info("Printing text response:", slog.Int("length", len(text)))

				return text, nil
			}

			for _, funCall := range funCalls {
				for arg, val := range funCall.Args {
					slog.Info(
						"Function call:",
						slog.String("name", funCall.Name),
						slog.String("arg", arg),
						slog.Any("value", val),
					)
				}

				funCallResponsePart := handleFunCall(funCall, g.WorkingDir)
				msgbuf = append(msgbuf, genai.NewContentFromParts([]*genai.Part{funCallResponsePart}, genai.RoleModel))
			}
		}

		return "", errors.New("LLM didn't generate a text response")
	}

	return asker, nil
}
