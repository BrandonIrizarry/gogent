package workingdir

import (
	"os"

	"github.com/BrandonIrizarry/gogent/internal/logger"
)

var lg logger.Logger

func InitLogger(logFile *os.File, verbose bool) {
	lg = logger.New(logFile, verbose, "workingdir")
}
