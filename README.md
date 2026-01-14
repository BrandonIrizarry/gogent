# Introduction

A CLI LLM coding agent using the Google Gemini API, written in Go.

This is currently a work in progress.

## Installation

```sh
    go install github.com/BrandonIrizarry/gogent@latest
```

## CLI Flags

### config

The path to the app's YAML configuration file. Defaults to `gogent.yaml`.

### log
The path to the app's logfile. Defaults to `logs.txt`.

### logmode
  
One or more comma-separated log-message types. The values specified
are the logging messages that will appear in the logfile.
    
Examples:

`gogent -logmode `
`gogent -logmode info,debug`
