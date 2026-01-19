package functions

import (
	"fmt"

	"google.golang.org/genai"
)

var declarations map[string]functionDeclarationConfig

func Init(workingDir string, maxFilesize int) {
	declarations = map[string]functionDeclarationConfig{
		"getFileContent": {
			declaration: genai.FunctionDeclaration{
				Name:        "getFileContent",
				Description: "Read file contents",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"path": {
							Type:        genai.TypeString,
							Description: "Relative path to file",
						},
					},
				},
			},

			fnObj: getFileContent{workingDir: workingDir, maxFilesize: maxFilesize},
		},

		"getFileContentRecursively": {
			declaration: genai.FunctionDeclaration{
				Name:        "getFileContentRecursively",
				Description: "Read contents of non-blacklisted files whose ancestor is the given directory.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"path": {
							Type:        genai.TypeString,
							Description: "Relative path to directory",
						},
						"depth": {
							Type:        genai.TypeInteger,
							Description: `How many directories deep to read file contents. The user can specify "no limit" to use an unbounded depth.`,
						},
					},
				},
			},

			fnObj: getFileContentRecursively{workingDir: workingDir, maxFilesize: maxFilesize},
		},

		"listDirectory": {
			declaration: genai.FunctionDeclaration{
				Name:        "listDirectory",
				Description: "List names of files in a directory, along with their sizes, and whether they're a directory.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"path": {
							Type:        genai.TypeString,
							Description: "Relative path to directory",
						},
					},
				},
			},

			fnObj: listDirectory{workingDir: workingDir},
		},
	}
}

func FunctionDeclarations() []*genai.FunctionDeclaration {
	buf := make([]*genai.FunctionDeclaration, 0)

	for _, cfg := range declarations {
		buf = append(buf, &cfg.declaration)
	}

	return buf
}

func FunctionObject(name string) (functionObject, error) {
	cfg, ok := declarations[name]
	if !ok {
		return nil, fmt.Errorf("Unknown function: %s", name)
	}

	return cfg.fnObj, nil
}
