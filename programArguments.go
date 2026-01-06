package main

import (
	"flag"
)

type programArguments struct {
	numIterations int
}

func newProgramArguments() (programArguments, error) {
	var pargs programArguments

	flag.IntVar(&pargs.numIterations, "num", 20, "The number of times the function call loop should execute")

	flag.Parse()

	return pargs, nil
}
