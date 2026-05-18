package runner

import (
	"context"
	"testing"
	"time"

	"github.com/kasjfulk/ollama-oneshot/internal/tools"
)

func TestRunnerNewDefaults(t *testing.T) {
	tools.ResetGlobals()
	tool := tools.Tool{
		Name:         "test",
		Command:      "echo",
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "test prompt")
	if r.ExitCode() != 0 {
		t.Errorf("expected default exit code 0, got %d", r.ExitCode())
	}
	if r.autoExit {
		t.Error("expected autoExit to be false by default")
	}
	if r.autoExitTimeout != defaultAutoExitTimeout {
		t.Errorf("expected default timeout %v, got %v", defaultAutoExitTimeout, r.autoExitTimeout)
	}
}

func TestRunnerSetAutoExit(t *testing.T) {
	tools.ResetGlobals()
	tool := tools.Tool{
		Name:    "test",
		Command: "echo",
	}
	r := New(tool, "test-model", "test prompt")
	r.SetAutoExit(true)
	if !r.autoExit {
		t.Error("expected autoExit to be true after SetAutoExit(true)")
	}
	r.SetAutoExit(false)
	if r.autoExit {
		t.Error("expected autoExit to be false after SetAutoExit(false)")
	}
}

func TestRunnerSetAutoExitTimeout(t *testing.T) {
	tools.ResetGlobals()
	tool := tools.Tool{
		Name:    "test",
		Command: "echo",
	}
	r := New(tool, "test-model", "test prompt")
	customTimeout := 5 * time.Minute
	r.SetAutoExitTimeout(customTimeout)
	if r.autoExitTimeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, r.autoExitTimeout)
	}
}

func TestRunnerEmptyCommand(t *testing.T) {
	tools.ResetGlobals()
	tool := tools.Tool{
		Name:         "empty",
		Command:      "",
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "test")
	r.SetAutoExit(false)

	err := r.Execute()
	if err == nil {
		t.Error("expected error for empty command")
	}
	if r.ExitCode() != 1 {
		t.Errorf("expected exit code 1 for empty command, got %d", r.ExitCode())
	}
}

func TestRunnerNonexistentCommand(t *testing.T) {
	tools.ResetGlobals()
	tool := tools.Tool{
		Name:         "nonexistent",
		Command:      "/nonexistent/path/to/binary",
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "test")
	r.SetAutoExit(false)
	r.SetAutoExitTimeout(5 * time.Second)

	err := r.Execute()
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
	if r.ExitCode() != 1 {
		t.Errorf("expected exit code 1 for nonexistent command, got %d", r.ExitCode())
	}
}

func TestSignalHandlerCleanup(t *testing.T) {
	cancel := func() {}
	sh := setupSignalHandler(cancel)
	sh.Cleanup()

	sh2 := setupSignalHandler(func() {})
	sh2.Cleanup()
	sh2.Cleanup()
}

func TestContextTimeoutIntegration(t *testing.T) {
	timeout := 800 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	<-ctx.Done()
	elapsed := time.Since(start)

	if elapsed < timeout/2 {
		t.Errorf("context cancelled too early: elapsed %v, timeout %v", elapsed, timeout)
	}
}
