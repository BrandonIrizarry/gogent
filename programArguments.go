package main

import (
	"errors"
	"flag"
)

type programArguments struct {
	initialPrompt string
	numIterations int
}

func newProgramArguments() (programArguments, error) {
	pargs := programArguments{}
	flag.StringVar(&pargs.initialPrompt, "prompt", "", "The initial user prompt")
	flag.IntVar(&pargs.numIterations, "num", 20, "The number of times the function call loop should execute")

	flag.Parse()

	if pargs.initialPrompt == "" {
		return programArguments{}, errors.New("-prompt needs an argument. Say something!")
	}

	return pargs, nil
}
