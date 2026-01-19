package gogent

import (
	"log/slog"

	"google.golang.org/genai"
)

type tokenCounts struct {
	Cached, Candidates, ToolUse, Prompt, Thoughts, Total int32
}

// incTokenCounts sums all individual token counts to provide an
// accessible report for the Gogent client. This is mainly so that the
// client can keep costs under control.
//
// The metadata parameter is the LLM response's metadata.
func (g *Gogent) incTokenCounts(metadata *genai.GenerateContentResponseUsageMetadata) {
	slog.Info(
		"Token Counts:",
		slog.Int("prompt", int(metadata.PromptTokenCount)),
		slog.Int("thoughts", int(metadata.ThoughtsTokenCount)),
		slog.Int("cached", int(metadata.CachedContentTokenCount)),
		slog.Int("candidates", int(metadata.CandidatesTokenCount)),
		slog.Int("tool_use", int(metadata.ToolUsePromptTokenCount)),
		slog.Int("total", int(metadata.TotalTokenCount)),
	)

	g.tokenCounts.Cached += metadata.CachedContentTokenCount
	g.tokenCounts.Candidates += metadata.CandidatesTokenCount
	g.tokenCounts.ToolUse += metadata.ToolUsePromptTokenCount
	g.tokenCounts.Prompt += metadata.PromptTokenCount
	g.tokenCounts.Thoughts += metadata.ThoughtsTokenCount
	g.tokenCounts.Total += metadata.TotalTokenCount

	// Also log the running totals.
	slog.Info(
		"Running Totals:",
		slog.Any("prompt", g.tokenCounts.Prompt),
		slog.Any("thoughts", g.tokenCounts.Thoughts),
		slog.Any("cached", g.tokenCounts.Cached),
		slog.Any("candidates", g.tokenCounts.Candidates),
		slog.Any("tool_use", g.tokenCounts.ToolUse),
		slog.Any("total", g.tokenCounts.Total),
	)
}
