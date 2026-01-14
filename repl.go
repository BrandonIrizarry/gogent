package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/BrandonIrizarry/gogent/internal/logger"
	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/charmbracelet/glamour"
	"google.golang.org/genai"
)

// repl launches a chat REPL with the agent, using the configuration
// parameters found in baseCfg.
func repl(maxIterations int, llmModel, renderStyle string, baseCfg baseconfig.BaseConfig) (err error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		return
	}

	tools := []*genai.Tool{
		{FunctionDeclarations: functions.FunctionDeclarations()},
	}

	msgBuf := msgbuf.NewMsgBuf()

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	for {
		initialPrompt, quit := getPrompt()

		if quit {
			return
		}

		msgBuf.AddText(initialPrompt)

		var response *genai.GenerateContentResponse
		for i := range maxIterations {
			logger.Info().Printf("New iteration: %d", i+1)

			response, err = client.Models.GenerateContent(
				ctx,
				llmModel,
				msgBuf.Messages,
				&contentConfig,
			)

			if err != nil {
				return
			}

			logger.Info().Printf("Prompt tokens: %d", response.UsageMetadata.PromptTokenCount)
			logger.Info().Printf("Response tokens: %d", response.UsageMetadata.ThoughtsTokenCount)

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
				logger.Info().Println("Printing text response:")

				text := response.Text()
				if out, err := glamour.Render(text, renderStyle); err != nil {
					logger.Info().Println("Glamour rendering failed, defaulting to plain text")
					fmt.Println(text)
				} else {
					fmt.Println(out)
				}

				break
			}

			for _, funCall := range funCalls {
				logger.Info().Printf("Function call name: %s", funCall.Name)

				for arg, val := range funCall.Args {
					logger.Info().Printf(" - argument: %s", arg)
					logger.Info().Printf(" - value: %v", val)
				}

				funCallResponsePart := handleFunCall(funCall, baseCfg)
				msgBuf.AddToolPart(funCallResponsePart)
			}
		}
	}
}

func handleFunCall(funCall *genai.FunctionCall, baseCfg baseconfig.BaseConfig) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)

	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args, baseCfg)
}

func getPrompt() (string, bool) {
	fmt.Println()
	fmt.Println("Ask the agent something (press Enter twice to submit your prompt)")
	fmt.Println("Submit a blank prompt to exit")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	var bld strings.Builder

	for scanner.Scan() {
		text := scanner.Text()

		if strings.TrimSpace(text) == "" {
			break
		}

		// Write an extra space, to make sure that words
		// across newline boundaries don't run on to each
		// other.
		bld.WriteString(" ")
		bld.WriteString(text)
	}

	// Nothing was written, meaning we must signal to our caller
	// to not invoke the agent REPL.
	if bld.Len() == 0 {
		return "", true
	}

	fmt.Println("Thinking...")
	return bld.String(), false
}
