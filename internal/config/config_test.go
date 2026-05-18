package config

import (
	"testing"
)

func TestApplyFlagsAutoExit(t *testing.T) {
	cfg := &Config{
		OllamaHost:    "localhost:11434",
		OllamaModel:   "test-model",
		DefaultTool:   "claude",
		AutoExit:      false,
	}

	cfg.ApplyFlags("", "", false, false, true)
	if !cfg.AutoExit {
		t.Error("expected AutoExit to be true after ApplyFlags with autoExit=true")
	}

	cfg.ApplyFlags("", "", false, false, false)
	if !cfg.AutoExit {
		t.Error("expected AutoExit to remain true when autoExit=false is passed (ApplyFlags only sets true)")
	}
}

func TestApplyFlagsAutoExitFalseByDefault(t *testing.T) {
	cfg := &Config{
		OllamaHost:    "localhost:11434",
		OllamaModel:   "test-model",
		DefaultTool:   "claude",
		AutoExit:      false,
	}

	cfg.ApplyFlags("", "", false, false, false)
	if cfg.AutoExit {
		t.Error("expected AutoExit to remain false when autoExit=false")
	}
}

func TestApplyFlagsAllFlags(t *testing.T) {
	cfg := &Config{
		OllamaHost:    "localhost:11434",
		OllamaModel:   "original-model",
		DefaultTool:   "codex",
		PromptEnhancement: true,
		YoloMode:      false,
		AutoExit:      false,
	}

	cfg.ApplyFlags("new-model", "claude", true, true, true)

	if cfg.OllamaModel != "new-model" {
		t.Errorf("expected model to be new-model, got %s", cfg.OllamaModel)
	}
	if cfg.DefaultTool != "claude" {
		t.Errorf("expected tool to be claude, got %s", cfg.DefaultTool)
	}
	if cfg.PromptEnhancement {
		t.Error("expected PromptEnhancement to be false after noEnhance=true")
	}
	if !cfg.YoloMode {
		t.Error("expected YoloMode to be true after yoloMode=true")
	}
	if !cfg.AutoExit {
		t.Error("expected AutoExit to be true after autoExit=true")
	}
}