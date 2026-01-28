package functions

import (
	"fmt"

	"google.golang.org/genai"
)

var declarations map[string]functionDeclarationConfig

const (
	PropertyPath  = "path"
	PropertyDepth = "depth"
)

var ignoredPaths = map[string]bool{}

func Init(workingDir string, maxFilesize int) error {
	var err error
	ignoredPaths, err = ignoredFilesMap(workingDir)
	if err != nil {
		return err
	}

	declarations = map[string]functionDeclarationConfig{
		"getFileContent": {
			declaration: genai.FunctionDeclaration{
				Name:        "getFileContent",
				Description: "Read file contents",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						PropertyPath: {
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
						PropertyPath: {
							Type:        genai.TypeString,
							Description: "Relative path to directory where files and other directories are contained",
						},
						PropertyDepth: {
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
				Description: "List the names of files in a directory, along with their sizes, and whether they're a directory.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						PropertyPath: {
							Type: genai.TypeString,
							Description: `Relative path to the specified directory. If unclear, confirm with
the user whether they mean the working directory.`,
						},
					},
				},
			},

			fnObj: listDirectory{workingDir: workingDir},
		},
	}

	return nil
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
		return nil, fmt.Errorf("unknown function: %s", name)
	}

	return cfg.fnObj, nil
}
