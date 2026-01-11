package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/functions"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	baseCfg := baseconfig.BaseConfig{CLIArguments: cliArgs, YAMLConfig: yamlCfg}

	if baseCfg.Verbose {
		fmt.Println()
		fmt.Println("Current settings:")
		fmt.Printf("Working directory: %s\n", baseCfg.WorkingDir)
		fmt.Printf("Max iterations: %d\n", baseCfg.MaxIterations)
		fmt.Printf("Max filesize: %d\n", baseCfg.MaxFilesize)
		fmt.Printf("Render style: %s\n", baseCfg.RenderStyle)
	}

	if err := repl(baseCfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}

func handleFunCall(funCall *genai.FunctionCall, baseCfg baseconfig.BaseConfig) *genai.Part {
	fnObj, err := functions.FunctionObject(funCall.Name)

	if err != nil {
		return functions.ResponseError(funCall.Name, fmt.Sprintf("Unknown function: %s", funCall.Name))
	}

	fn := fnObj.Function()

	return fn(funCall.Args, baseCfg)
}
