package main

import (
	"os"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

type localConfig struct {
	logFile *os.File
	log     logger.Logger
}

func initConfig(logFile *os.File, verbose bool) localConfig {
	var cfg localConfig

	cfg.logFile = logFile
	cfg.log = logger.New(logFile, verbose, "main")

	return cfg
}
