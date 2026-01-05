package functions

import "google.golang.org/genai"

var declarations = map[string]genai.FunctionDeclaration{
	"getFileContent": genai.FunctionDeclaration{
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
}

func FunctionDeclarations() []*genai.FunctionDeclaration {
	buf := make([]*genai.FunctionDeclaration, 0)

	for _, decl := range declarations {
		buf = append(buf, &decl)
	}

	return buf
}
