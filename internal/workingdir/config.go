package workingdir

import (
	"os"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

type config struct {
	log logger.Logger
}

func InitConfig(logFile *os.File, verbose bool) config {
	var cfg config

	cfg.log = logger.New(logFile, verbose, "workingdir")

	return cfg
}
