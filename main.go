package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var systemInstruction = `
You are a helpful AI coding agent.

When a user asks a question or makes a request, make a function call
plan. You can perform the following operations:

- Read file contents

All paths you provide are relative to some working directory. You must
not specify the working directory in your function calls; for security
reasons, the tool dispatch code will handle that.

`

func getPrompt() (string, bool) {
	fmt.Println()
	fmt.Println("Ask the agent something (press Enter twice to submit your prompt)")
	fmt.Println("Submit a blank prompt to exit the agent REPL")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	bld := strings.Builder{}

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

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// For now, only get 'numIterations' from the user.
	pargs, err := newProgramArguments()

	if err != nil {
		log.Fatal(err)
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

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	// The REPL.
	for {
		initialPrompt, quit := getPrompt()
		msgBuf.AddText(initialPrompt)

		if quit {
			os.Exit(0)
		}

		for range pargs.numIterations {
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

				msgBuf.AddToolPart(funCallResponsePart)
			}
		}
	}
}
