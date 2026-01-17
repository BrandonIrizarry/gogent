package gogent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/BrandonIrizarry/gogent/internal/functions"
	"google.golang.org/genai"
)

// FIXME: frontends are responsible for forwarding this stuff to
// gogent's own internal data model.
const (
	maxIterations = 20
	llmModel      = "gemini-2.5-flash-lite-preview-09-2025"
	workingDir    = "."
	maxFilesize   = 100_000
)

type Gogent struct {
	prompt string
}

func init() {
	functions.Init(workingDir, maxFilesize)
}

// Write implements the [io.Writer] interface. It offloads the bytes
// of p into Gogent. For now, all this does is tell Gogent about the
// prompt; no thinking is done yet.
func (h *Gogent) Write(p []byte) (n int, err error) {
	h.prompt = string(p)

	return len(p), nil
}

// Read implements the [io.Reader] interface. It loads the LLM's
// response into p. This is where the LLM performs its thinking
// operation.
func (h Gogent) Read(p []byte) (n int, err error) {
	response, err := think(string(p))
	if err != nil {
		return 0, err
	}

	return len(response), nil
}

func think(prompt string) (string, error) {
	msgBuf := NewMsgBuf()
	msgBuf.AddText(prompt)

	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	tools := []*genai.Tool{
		{FunctionDeclarations: functions.FunctionDeclarations()},
	}

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	var response *genai.GenerateContentResponse

	for i := range maxIterations {
		slog.Info("New iteration", slog.Int("iteration", i))

		response, err = client.Models.GenerateContent(
			ctx,
			llmModel,
			msgBuf.Messages,
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
			msgBuf.AddMessage(candidate.Content)
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
			msgBuf.AddToolPart(funCallResponsePart)
		}
	}

	return "", errors.New("LLM didn't generate a text response")
}

func handleFunCall(funCall *genai.FunctionCall) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)
	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args)
}
