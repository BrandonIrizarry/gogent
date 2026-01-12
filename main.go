package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/workingdir"
	"github.com/BrandonIrizarry/gogent/internal/yamlconfig"
	"github.com/joho/godotenv"
)

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// CLI arguments.
	cliArgs, err := cliargs.NewCLIArguments()
	if err != nil {
		log.Fatal(err)
	}

	// Open up the log file.
	logFile, err := os.OpenFile(cliArgs.LogFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// YAML configuration.
	yamlCfg, err := yamlconfig.NewYAMLConfig(cliArgs.ConfigFilename)
	if err != nil {
		log.Fatal(err)
	}

	// Get the working directory from a TUI file picker.
	wdirCfg := workingdir.InitConfig(logFile, cliArgs.Verbose)

	wdir, err := wdirCfg.SelectWorkingDir()
	if err != nil {
		log.Fatal(err)
	}

	baseCfg := baseconfig.BaseConfig{
		WorkingDir:  wdir,
		MaxFilesize: yamlCfg.MaxFilesize,
	}

	// package main's local config struct.
	cfg := initConfig(logFile, cliArgs.Verbose)

	cfg.log.Verbose.Println()
	cfg.log.Verbose.Println("Current settings:")
	cfg.log.Verbose.Printf("Working directory: %s\n", wdir)
	cfg.log.Verbose.Printf("Max iterations: %d\n", yamlCfg.MaxIterations)
	cfg.log.Verbose.Printf("Max filesize: %d\n", yamlCfg.MaxFilesize)
	cfg.log.Verbose.Printf("Render style: %s\n", yamlCfg.RenderStyle)
	cfg.log.Verbose.Printf("Model: %s\n", yamlCfg.Model)

	if err := cfg.repl(yamlCfg.MaxIterations, yamlCfg.Model, yamlCfg.RenderStyle, baseCfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}
