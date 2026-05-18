package tools

func init() {
	registry["codex"] = Tool{Name: "codex", Command: "ollama launch codex"}
	registry["codex-app"] = Tool{Name: "codex-app", Command: "ollama launch codex-app"}
}