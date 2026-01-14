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
	//
	// Note that, since we don't have our custom logger yet, we're
	// using the default logger for now.
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// CLI arguments.
	cliArgs, err := cliargs.NewCLIArguments()
	if err != nil {
		panic(err)
	}

	// Open up the log file.
	logFile, err := os.OpenFile(cliArgs.LogFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger.Init(logFile, cliArgs.LogMode)

	// YAML configuration.
	yamlCfg, err := yamlconfig.NewYAMLConfig(cliArgs.ConfigFilename)
	if err != nil {
		logger.Info().Fatal(err)
	}

	wdir, err := workingdir.SelectWorkingDir()
	if err != nil {
		logger.Info().Fatal(err)
	}

	baseCfg := baseconfig.BaseConfig{
		WorkingDir:  wdir,
		MaxFilesize: yamlCfg.MaxFilesize,
	}

	logger.Info().Println()
	logger.Info().Println("Current settings:")
	logger.Info().Printf("Working directory: %s\n", wdir)
	logger.Info().Printf("Max iterations: %d\n", yamlCfg.MaxIterations)
	logger.Info().Printf("Max filesize: %d\n", yamlCfg.MaxFilesize)
	logger.Info().Printf("Render style: %s\n", yamlCfg.RenderStyle)
	logger.Info().Printf("Model: %s\n", yamlCfg.Model)

	if err := repl(yamlCfg.MaxIterations, yamlCfg.Model, yamlCfg.RenderStyle, baseCfg); err != nil {
		logger.Info().Fatal(err)
	}

	fmt.Println("Bye, come again soon!")
}
