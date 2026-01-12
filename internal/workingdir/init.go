package workingdir

import (
	"log"
	"os"
)

var logger *log.Logger
var logFile *os.File

func init() {
	var err error

	logFile, err = os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic("Can't open logfile")
	}

	logger = log.New(logFile, "WORKINGDIR: ", log.LstdFlags|log.Lshortfile)
}
