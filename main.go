package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var systemInstruction = `
You are a helpful AI coding agent.

When a user asks a question or makes a request, make a function call
plan. You can perform the following operations:

- Read file contents
`

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	tools := []*genai.Tool{
		{FunctionDeclarations: []*genai.FunctionDeclaration{&getFileContent}},
	}

	msgBuf := msgbuf.NewMsgBuf()
	msgBuf.AddText("Describe the contents of main.go in this same directory")

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	// FIXME: 20 is hardcoded as the number of times to attempt
	// the function-call loop.
	for range 20 {
		response, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			msgBuf.Messages,
			&contentConfig,
		)

		if err != nil {
			log.Fatal(err)
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
			fmt.Println(response.Text())
			break
		}

		for _, funCall := range funCalls {
			var funCallResponsePart *genai.Part

			switch funCall.Name {
			case "getFileContent":
				// Read the contents of the given file.
				path := funCall.Args["filepath"].(string)
				dat, err := os.ReadFile(path)

				if err != nil {
					log.Fatal(err)
				}

				fileContents := string(dat)

				funCallResponsePart = genai.NewPartFromFunctionResponse("getFileContent", map[string]any{
					"result": fileContents,
				})
			default:
				funCallResponsePart = genai.NewPartFromFunctionResponse(funCall.Name, map[string]any{
					"error": fmt.Sprintf("Unknown function: %s", funCall.Name),
				})
			}

			msgBuf.AddMessage(&genai.Content{
				Role:  "tool",
				Parts: []*genai.Part{funCallResponsePart},
			})
		}
	}
}
