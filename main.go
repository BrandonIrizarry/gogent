package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"github.com/BrandonIrizarry/gogent/internal/cliargs"
	"github.com/BrandonIrizarry/gogent/internal/logger"
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

	// The logger local to main.
	lg := logger.New(logFile, cliArgs.Verbose, "main")

	lg.Verbose.Println()
	lg.Verbose.Println("Current settings:")
	lg.Verbose.Printf("Working directory: %s\n", wdir)
	lg.Verbose.Printf("Max iterations: %d\n", yamlCfg.MaxIterations)
	lg.Verbose.Printf("Max filesize: %d\n", yamlCfg.MaxFilesize)
	lg.Verbose.Printf("Render style: %s\n", yamlCfg.RenderStyle)
	lg.Verbose.Printf("Model: %s\n", yamlCfg.Model)

	if err := repl(yamlCfg.MaxIterations, yamlCfg.Model, yamlCfg.RenderStyle, baseCfg, lg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}
