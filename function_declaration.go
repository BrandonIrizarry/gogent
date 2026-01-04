package main

import "google.golang.org/genai"

var getFileContent = genai.FunctionDeclaration{
	Name:        "getFileContent",
	Description: "Read file contents",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"filepath": {
				Type:        genai.TypeString,
				Description: "Relative path to file",
			},
		},
	},
}
