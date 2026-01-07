package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/BrandonIrizarry/gogent/internal/msgbuf"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
	"github.com/charmbracelet/glamour"
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

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cliArgs, err := cliargs.NewCLIArguments()
	if err != nil {
		log.Fatal(err)
	}

	yamlCfg, err := yamlconfig.NewYAMLConfig("gogent.yaml")
	if err != nil {
		log.Fatal(err)
	}

	baseCfg := baseconfig.BaseConfig{cliArgs, yamlCfg}

	if baseCfg.Verbose {
		fmt.Println()
		fmt.Println("Current settings:")
		fmt.Printf("Working directory: %s\n", baseCfg.WorkingDir)
		fmt.Printf("Max iterations: %d\n", baseCfg.MaxIterations)
		fmt.Printf("Max filesize: %d\n", baseCfg.MaxFilesize)
	}

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

		for i := 0; i < baseCfg.MaxIterations; i++ {
			if cliArgs.Verbose {
				fmt.Println()
				log.Printf("New iteration: %d", i+1)
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

			if cliArgs.Verbose {
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
				if cliArgs.Verbose {
					log.Println("No more function calls; printing text response")
					fmt.Println()
				}

				text := response.Text()
				if out, err := glamour.Render(text, "light"); err != nil {
					log.Println("Glamour rendering failed, defaulting to plain text")
					fmt.Println(text)
				} else {
					fmt.Println(out)
				}

				break
			}

			for _, funCall := range funCalls {
				if cliArgs.Verbose {
					log.Printf("Function call name: %s", funCall.Name)
					log.Printf("Function call args: %#v", funCall.Args)
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
