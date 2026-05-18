package tools

func init() {
	registry["claude"] = Tool{
		Name:         "claude",
		Command:      "ollama launch claude",
		YoloArgs:     []string{"--", "--dangerously-skip-permissions"},
		AutoExitArgs: []string{"-p"},
	}
	registry["claude-code"] = Tool{
		Name:         "claude-code",
		Command:      "ollama launch claude",
		YoloArgs:     []string{"--", "--dangerously-skip-permissions"},
		AutoExitArgs: []string{"-p"},
	}
}
