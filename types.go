package gogent

type askerFn func(string) (string, []TokenCount, error)

type Gogent struct {
	WorkingDir    string
	MaxFilesize   int
	LLMModel      string
	MaxIterations int
}

type TokenCount struct {
	Prompt     int32
	Thoughts   int32
	Cached     int32
	Candidates int32
	ToolUse    int32
	Total      int32
}
