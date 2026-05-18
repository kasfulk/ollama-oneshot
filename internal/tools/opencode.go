package tools

func init() {
	registry["opencode"] = Tool{
		Name:         "opencode",
		Command:      "ollama launch opencode",
		YoloArgs:     []string{},
		AutoExitArgs: []string{"run"},
	}
	registry["openclaw"] = Tool{
		Name:         "openclaw",
		Command:      "ollama launch openclaw",
		YoloArgs:     []string{},
		AutoExitArgs: []string{"run"},
	}
	registry["hermes"] = Tool{
		Name:         "hermes",
		Command:      "ollama launch hermes",
		YoloArgs:     []string{},
		AutoExitArgs: []string{"run"},
	}
}
