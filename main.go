package main

import (
	"fmt"
	"log"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
	"github.com/joho/godotenv"
)

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
