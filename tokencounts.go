package gogent

import (
	"github.com/rs/zerolog/log"
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
	log.Info().
		Int32("prompt", metadata.PromptTokenCount).
		Int32("prompt", metadata.PromptTokenCount).
		Int32("thoughts", metadata.ThoughtsTokenCount).
		Int32("cached", metadata.CachedContentTokenCount).
		Int32("candidates", metadata.CandidatesTokenCount).
		Int32("tool_use", metadata.ToolUsePromptTokenCount).
		Int32("total", metadata.TotalTokenCount).
		Msg("Token counts:")

	g.tokenCounts.Cached += metadata.CachedContentTokenCount
	g.tokenCounts.Candidates += metadata.CandidatesTokenCount
	g.tokenCounts.ToolUse += metadata.ToolUsePromptTokenCount
	g.tokenCounts.Prompt += metadata.PromptTokenCount
	g.tokenCounts.Thoughts += metadata.ThoughtsTokenCount
	g.tokenCounts.Total += metadata.TotalTokenCount

	// Also log the running totals.
	log.Info().
		Int32("prompt", g.tokenCounts.Prompt).
		Int32("thoughts", g.tokenCounts.Thoughts).
		Int32("cached", g.tokenCounts.Cached).
		Int32("candidates", g.tokenCounts.Candidates).
		Int32("tool_use", g.tokenCounts.ToolUse).
		Int32("total", g.tokenCounts.Total).
		Msg("Running totals:")
}
