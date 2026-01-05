package functions

import (
	"fmt"

	"google.golang.org/genai"
)

type functionDeclarationConfig struct {
	declaration genai.FunctionDeclaration
	function    func(map[string]any) *genai.Part
}

var declarations = map[string]functionDeclarationConfig{
	"getFileContent": {
		declaration: genai.FunctionDeclaration{
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
		},

		function: getFileContent,
	},
}

func FunctionDeclarations() []*genai.FunctionDeclaration {
	buf := make([]*genai.FunctionDeclaration, 0)

	for _, cfg := range declarations {
		buf = append(buf, &cfg.declaration)
	}

	return buf
}

func Function(name string) (func(map[string]any) *genai.Part, error) {
	cfg, ok := declarations[name]

	if !ok {
		return nil, fmt.Errorf("Unknown function: %s", name)
	}

	return cfg.function, nil
}
