package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/charmbracelet/glamour"
	"google.golang.org/genai"
)

// repl launches a chat REPL with the agent, using the configuration
// parameters found in baseCfg.
func repl(baseCfg baseconfig.BaseConfig) (err error) {
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
		for i := 0; i < baseCfg.MaxIterations; i++ {
			if baseCfg.Verbose {
				fmt.Println()
				log.Printf("New iteration: %d", i+1)
			}

			response, err = client.Models.GenerateContent(
				ctx,
				baseCfg.Model,
				msgBuf.Messages,
				&contentConfig,
			)

			if err != nil {
				return
			}

			if baseCfg.Verbose {
				log.Printf("Prompt tokens: %d", response.UsageMetadata.PromptTokenCount)
				log.Printf("Response tokens: %d", response.UsageMetadata.ThoughtsTokenCount)
			}

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
				if baseCfg.Verbose {
					log.Println("Printing text response:")
					fmt.Println()
				}

				text := response.Text()
				if out, err := glamour.Render(text, baseCfg.RenderStyle); err != nil {
					log.Println("Glamour rendering failed, defaulting to plain text")
					fmt.Println(text)
				} else {
					fmt.Println(out)
				}

				break
			}

			for _, funCall := range funCalls {
				if baseCfg.Verbose {
					log.Printf("Function call name: %s", funCall.Name)

					for arg, val := range funCall.Args {
						log.Printf(" - argument: %s", arg)
						log.Printf(" - value: %v", val)
					}
				}

				funCallResponsePart := handleFunCall(funCall, baseCfg)
				msgBuf.AddToolPart(funCallResponsePart)
			}
		}
	}
}
