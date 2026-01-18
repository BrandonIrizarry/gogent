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

type functionCallLoop func(string) (string, error)

func (g Gogent) Init() (functionCallLoop, error) {
	// Set the appropriate logging level.
	if g.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// Initialize any state needed by the function call objects
	// themselves.
	functions.Init(g.WorkingDir, g.MaxFilesize)

	// Initialize the LLM configuration, which must persist across
	// calls to [Gogent.Ask].
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	tools := []*genai.Tool{
		{FunctionDeclarations: functions.FunctionDeclarations()},
	}

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	msgbuf := NewMsgBuf()

	// This is the actual code that processes the user prompt.
	fn := func(prompt string) (string, error) {
		// Initialize the message buffer with the user
		// prompt. I'm making a careful note that this should
		// be outside the function-call loop.
		msgbuf.AddText(prompt)

		for i := range g.MaxIterations {
			slog.Info("Start of function-call loop:", slog.Int("iteration", i))

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
				"metadata",
				slog.Int("prompt tokens", int(response.UsageMetadata.PromptTokenCount)),
				slog.Int("response tokens", int(response.UsageMetadata.ThoughtsTokenCount)),
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
				slog.Info("Printing text response:")

				return response.Text(), nil
			}

			for _, funCall := range funCalls {
				for arg, val := range funCall.Args {
					slog.Info(
						"Function call",
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

	return fn, nil
}

func handleFunCall(funCall *genai.FunctionCall) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args)
}
