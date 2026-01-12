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
	//
	// Note that, since we don't have our custom logger yet, we're
	// using the default logger for now.
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

	// Now that we have the log file, define package main's local
	// config struct. We can finally use our custom logger from
	// now on.
	cfg := initConfig(logFile, cliArgs.Verbose)

	// YAML configuration.
	yamlCfg, err := yamlconfig.NewYAMLConfig(cliArgs.ConfigFilename)
	if err != nil {
		cfg.log.Info.Fatal(err)
	}

	// Get the working directory from a TUI file picker.
	wdirCfg := workingdir.InitConfig(logFile, cliArgs.Verbose)

	wdir, err := wdirCfg.SelectWorkingDir()
	if err != nil {
		cfg.log.Info.Fatal(err)
	}

	baseCfg := baseconfig.BaseConfig{
		WorkingDir:  wdir,
		MaxFilesize: yamlCfg.MaxFilesize,
	}

	cfg.log.Verbose.Println()
	cfg.log.Verbose.Println("Current settings:")
	cfg.log.Verbose.Printf("Working directory: %s\n", wdir)
	cfg.log.Verbose.Printf("Max iterations: %d\n", yamlCfg.MaxIterations)
	cfg.log.Verbose.Printf("Max filesize: %d\n", yamlCfg.MaxFilesize)
	cfg.log.Verbose.Printf("Render style: %s\n", yamlCfg.RenderStyle)
	cfg.log.Verbose.Printf("Model: %s\n", yamlCfg.Model)

	if err := cfg.repl(yamlCfg.MaxIterations, yamlCfg.Model, yamlCfg.RenderStyle, baseCfg); err != nil {
		cfg.log.Info.Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}
