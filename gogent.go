package gogent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

type Gogent struct {
	WorkingDir    string
	MaxFilesize   int
	LLMModel      string
	MaxIterations int
	Debug         bool
}

type askerFn func(string) (string, error)

// Init initializes state used by the LLM, providing both the values
// of the Gogent struct's own fields, as well as setting up any state
// the LLM client needs to persist across prompt/response cycles. It
// returns a function that clients can use to initiate a single
// prompt/response cycle, likely in the context of some kind of REPL.
func (g Gogent) Init() (askerFn, error) {
	// Set the appropriate logging level.
	if g.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

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

	msgbuf := NewMsgBuf()

	// This is the actual code that processes the user prompt.
	//
	// Note that this function captures many of the configuration
	// parameters defined just above (like 'client'.)
	asker := func(prompt string) (string, error) {
		// Initialize the message buffer with the user
		// prompt. I'm making a careful note that this should
		// be outside the function-call loop.
		msgbuf.AddText(prompt)

		for i := range g.MaxIterations {
			slog.Info("start of function-call loop:", slog.Int("iteration", i))

			response, err := client.Models.GenerateContent(
				ctx,
				g.LLMModel,
				msgbuf.Messages,
				&contentConfig,
			)
			if err != nil {
				return "", err
			}

			slog.Info(
				"Token Counts:",
				slog.Int("prompt", int(response.UsageMetadata.PromptTokenCount)),
				slog.Int("response", int(response.UsageMetadata.ThoughtsTokenCount)),
				slog.Int("cached", int(response.UsageMetadata.CachedContentTokenCount)),
				slog.Int("candidates", int(response.UsageMetadata.CandidatesTokenCount)),
				slog.Int("tool_use", int(response.UsageMetadata.ToolUsePromptTokenCount)),
				slog.Int("total", int(response.UsageMetadata.TotalTokenCount)),
			)

			// Add the candidates to the message buffer. This
			// conforms both to the Gemini documentation, as well
			// as the Boot.dev AI Agent project.
			for _, candidate := range response.Candidates {
				msgbuf.AddMessage(candidate.Content)
			}

			// Check if the LLM has proposed any function calls to
			// act upon.
			funCalls := response.FunctionCalls()

			// The LLM is ready to give a textual response.
			if len(funCalls) == 0 {
				slog.Info("printing text response:")

				return response.Text(), nil
			}

			for _, funCall := range funCalls {
				for arg, val := range funCall.Args {
					slog.Info(
						"which function call",
						slog.String("name", funCall.Name),
						slog.String("arg", arg),
						slog.Any("value", val),
					)
				}

				funCallResponsePart := handleFunCall(funCall)
				msgbuf.AddToolPart(funCallResponsePart)
			}
		}

		return "", errors.New("LLM didn't generate a text response")
	}

	return asker, nil
}

func handleFunCall(funCall *genai.FunctionCall) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args)
}
