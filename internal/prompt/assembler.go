package prompt

import "fmt"

type Assembler struct {
	systemPrompt string
}

func NewAssembler(systemOverride string) *Assembler {
	sys := DefaultSystemPrompt
	if systemOverride != "" {
		sys = systemOverride
	}
	return &Assembler{systemPrompt: sys}
}

func (a *Assembler) Assemble(enhancedPrompt, userPrompt, docContent string) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[SYSTEM]\n%s", a.systemPrompt))

	if docContent != "" {
		parts = append(parts, fmt.Sprintf("[DOCUMENTATION]\n%s", docContent))
	}

	if enhancedPrompt != userPrompt {
		parts = append(parts, fmt.Sprintf("[ENHANCED PROMPT]\n%s", enhancedPrompt))
	}

	parts = append(parts, fmt.Sprintf("[USER PROMPT]\n%s", userPrompt))

	joined := ""
	for i, part := range parts {
		if i > 0 {
			joined += "\n\n"
		}
		joined += part
	}
	return joined
}