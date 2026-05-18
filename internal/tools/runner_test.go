package tools

import (
	"strings"
	"testing"
)

func TestLaunchCommandWithoutModes(t *testing.T) {
	ResetGlobals()
	tool := Tool{
		Name:         "test",
		Command:      "ollama launch test",
		YoloArgs:     []string{"--", "--approve"},
		AutoExitArgs: []string{"--", "--exit"},
	}
	cmd := tool.LaunchCommand("my-model")
	expected := "ollama launch test --model my-model"
	if cmd != expected {
		t.Errorf("expected %q, got %q", expected, cmd)
	}
}

func TestLaunchCommandWithYoloMode(t *testing.T) {
	ResetGlobals()
	SetYoloMode(true)
	tool := Tool{
		Name:         "test",
		Command:      "ollama launch test",
		YoloArgs:     []string{"--", "--approve"},
		AutoExitArgs: []string{"--", "--exit"},
	}
	cmd := tool.LaunchCommand("my-model")
	if !strings.Contains(cmd, "--approve") {
		t.Errorf("expected command to contain --approve, got %q", cmd)
	}
	if strings.Contains(cmd, "--exit") {
		t.Errorf("expected command NOT to contain --exit when only yolo mode is on, got %q", cmd)
	}
}

func TestLaunchCommandWithAutoExitMode(t *testing.T) {
	ResetGlobals()
	SetAutoExitMode(true)
	tool := Tool{
		Name:         "test",
		Command:      "ollama launch test",
		YoloArgs:     []string{"--", "--approve"},
		AutoExitArgs: []string{"--", "--exit"},
	}
	cmd := tool.LaunchCommand("my-model")
	if !strings.Contains(cmd, "--exit") {
		t.Errorf("expected command to contain --exit, got %q", cmd)
	}
	if strings.Contains(cmd, "--approve") {
		t.Errorf("expected command NOT to contain --approve when only auto-exit mode is on, got %q", cmd)
	}
}

func TestLaunchCommandWithBothModes(t *testing.T) {
	ResetGlobals()
	SetYoloMode(true)
	SetAutoExitMode(true)
	tool := Tool{
		Name:         "test",
		Command:      "ollama launch test",
		YoloArgs:     []string{"--", "--approve"},
		AutoExitArgs: []string{"--", "--exit"},
	}
	cmd := tool.LaunchCommand("my-model")
	if !strings.Contains(cmd, "--approve") {
		t.Errorf("expected command to contain --approve, got %q", cmd)
	}
	if !strings.Contains(cmd, "--exit") {
		t.Errorf("expected command to contain --exit, got %q", cmd)
	}
}

func TestLaunchCommandAutoExitNoArgs(t *testing.T) {
	ResetGlobals()
	SetAutoExitMode(true)
	tool := Tool{
		Name:         "test-noargs",
		Command:      "ollama launch test-noargs",
		AutoExitArgs: []string{},
	}
	cmd := tool.LaunchCommand("my-model")
	expected := "ollama launch test-noargs --model my-model"
	if cmd != expected {
		t.Errorf("expected %q, got %q", expected, cmd)
	}
}

func TestLaunchCommandYoloNoArgs(t *testing.T) {
	ResetGlobals()
	SetYoloMode(true)
	tool := Tool{
		Name:     "test-noargs",
		Command:  "ollama launch test-noargs",
		YoloArgs: []string{},
	}
	cmd := tool.LaunchCommand("my-model")
	expected := "ollama launch test-noargs --model my-model"
	if cmd != expected {
		t.Errorf("expected %q, got %q", expected, cmd)
	}
}

func TestClaudeToolAutoExitArgs(t *testing.T) {
	ResetGlobals()
	tool, ok := Get("claude")
	if !ok {
		t.Skip("claude tool not registered")
	}
	if len(tool.AutoExitArgs) == 0 {
		t.Error("expected claude tool to have AutoExitArgs")
	}
}

func TestCodexToolAutoExitArgs(t *testing.T) {
	ResetGlobals()
	tool, ok := Get("codex")
	if !ok {
		t.Skip("codex tool not registered")
	}
	if len(tool.AutoExitArgs) == 0 {
		t.Error("expected codex tool to have AutoExitArgs")
	}
}