package enhancer

const EnhancementSystemPrompt = `You are an expert software architect.
Convert vague software requests into structured, implementation-ready specifications.
Your output must include:
- Implementation plan with clear steps
- Architecture overview
- Engineering tasks breakdown
- Technology stack recommendations
- Security and scalability considerations

Focus on: backend, frontend, database, deployment, security, maintainability.
Output only the enhanced prompt, no meta-commentary.`

func buildEnhancementPrompt(userPrompt string) string {
	wordCount := len(splitWords(userPrompt))
	if wordCount < 50 {
		return `Expand this vague request into a comprehensive, implementation-ready specification. ` +
			`Include architecture, tasks, tech stack, and deployment strategy:\n\n` + userPrompt
	}
	return `Normalize and structure this request into a clear engineering specification. ` +
		`Preserve the detail level but add structure:\n\n` + userPrompt
}

func splitWords(s string) []string {
	count := 0
	inWord := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			inWord = false
		} else if !inWord {
			inWord = true
			count++
		}
	}
	return make([]string, count)
}