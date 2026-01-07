package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var systemInstruction = `
You are a helpful AI coding agent.

When a user asks a question or makes a request, make a function call
plan. You can perform the following operations:

- Read file contents
- List directory contents

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

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cliArgs, err := cliargs.NewCLIArguments()

	if err != nil {
		if errors.Is(err, cliargs.ErrBadDefaultDir) {
			fmt.Println("OK, rerun with -dir $MY_PATH (can be relative)")
			os.Exit(0)
		}

		log.Fatal(err)
	}

	fmt.Printf("OK, setting working directory to %s\n", cliArgs.WorkingDir)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	tools := []*genai.Tool{
		{FunctionDeclarations: functions.FunctionDeclarations()},
	}

	msgBuf := msgbuf.NewMsgBuf()

	contentConfig := genai.GenerateContentConfig{
		Tools:             tools,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemInstruction}}},
	}

	// The REPL.
	for {
		initialPrompt, quit := getPrompt()

		if quit {
			fmt.Println("Bye, come again soon!")
			os.Exit(0)
		}

		msgBuf.AddText(initialPrompt)

		for range cliArgs.NumIterations {
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
				funCallResponsePart := handleFunCall(funCall, cliArgs)
				msgBuf.AddToolPart(funCallResponsePart)
			}
		}
	}
}

func handleFunCall(funCall *genai.FunctionCall, cliArgs cliargs.CLIArguments) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)

	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args, cliArgs)
}
