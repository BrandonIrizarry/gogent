package main

import (
	"fmt"
	"log"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/workingdir"
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

	baseCfg, err := baseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if baseCfg.Verbose {
		fmt.Println()
		fmt.Println("Current settings:")
		fmt.Printf("Working directory: %s\n", baseCfg.WorkingDir)
		fmt.Printf("Max iterations: %d\n", baseCfg.MaxIterations)
		fmt.Printf("Max filesize: %d\n", baseCfg.MaxFilesize)
		fmt.Printf("Render style: %s\n", baseCfg.RenderStyle)
		fmt.Printf("Model: %s\n", baseCfg.Model)
	}

	if err := repl(baseCfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}

func baseConfig() (baseconfig.BaseConfig, error) {
	// CLI arguments.
	cliArgs, err := cliargs.NewCLIArguments()
	if err != nil {
		return baseconfig.BaseConfig{}, err
	}

	// YAML configuration.
	yamlCfg, err := yamlconfig.NewYAMLConfig(cliArgs.ConfigFilename)
	if err != nil {
		return baseconfig.BaseConfig{}, err
	}

	// Get the working directory from a TUI file picker.
	wdir, err := workingdir.SelectWorkingDir()
	if err != nil {
		return baseconfig.BaseConfig{}, err
	}

	baseCfg := baseconfig.BaseConfig{
		CLIArguments: cliArgs,
		YAMLConfig:   yamlCfg,
		WorkingDir:   wdir,
	}

	return baseCfg, nil
}
