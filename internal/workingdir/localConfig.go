package workingdir

import (
	"os"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

type localConfig struct {
	log logger.Logger
}

func InitConfig(logFile *os.File, verbose bool) localConfig {
	var cfg localConfig

	cfg.log = logger.New(logFile, verbose, "workingdir")

	return cfg
}
