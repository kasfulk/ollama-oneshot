package tools

func init() {
	registry["opencode"] = Tool{Name: "opencode", Command: "ollama launch opencode"}
	registry["openclaw"] = Tool{Name: "openclaw", Command: "ollama launch openclaw"}
	registry["hermes"] = Tool{Name: "hermes", Command: "ollama launch hermes"}
}