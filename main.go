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

	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		msgBuf.Messages,
		&contentConfig,
	)

	if err != nil {
		log.Fatal(err)
	}

	funCall := response.Candidates[0].Content.Parts[0].FunctionCall
	fmt.Println(funCall.Args)
	fmt.Println(funCall.Name)

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

	msgBuf.AddMessage(response.Candidates[0].Content)
	msgBuf.AddMessage(&genai.Content{
		Role:  "tool",
		Parts: []*genai.Part{funCallResponsePart},
	})

	response, err = client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		msgBuf.Messages,
		&contentConfig,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.FunctionCalls())
	fmt.Println(response.Text())
}
