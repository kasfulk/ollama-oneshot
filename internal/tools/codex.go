package tools

func init() {
	registry["codex"] = Tool{
		Name:         "codex",
		Command:      "ollama launch codex",
		YoloArgs:     []string{"--full-auto"},
		AutoExitArgs: []string{"exec"},
	}
	registry["codex-app"] = Tool{
		Name:         "codex-app",
		Command:      "ollama launch codex-app",
		YoloArgs:     []string{"--full-auto"},
		AutoExitArgs: []string{"exec"},
	}
}
