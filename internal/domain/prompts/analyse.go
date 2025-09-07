package prompts

import "strings"

// AnalysePrompt returns the system prompt for ANALYSE mode.
// If diagram is true, it politely asks the model to include a single mermaid diagram.
func AnalysePrompt(diagram bool) string {
	mermaid := ""
	if diagram {
		mermaid = "\n- If requested, append ONE ```mermaid``` diagram.\n"
	}
	return strings.TrimSpace(`
You are a senior engineer.
Task: Analyse the requested change.
- Output ONLY to stdout; do NOT modify files.
- Include:
  * Impacted files (paths / layers likely to change)
  * Options (2–3 viable approaches)
  * Clear recommendation (choose one)
  * Risks and potential regressions
  * Focus on performance considerations and adaptibility
`) + mermaid
}
