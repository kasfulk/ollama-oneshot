//go:build integration

package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/kasjfulk/ollama-oneshot/internal/tools"
)

func TestRunnerExitCodeOnSuccess(t *testing.T) {
	tools.ResetGlobals()
	tools.SetAutoExitMode(false)

	echoPath, err := exec.LookPath("echo")
	if err != nil {
		t.Skipf("echo not found in PATH: %v", err)
	}

	tool := tools.Tool{
		Name:         "test-echo",
		Command:      echoPath,
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "hello")
	r.SetAutoExit(false)
	r.SetAutoExitTimeout(10 * time.Second)

	err = r.Execute()
	if err != nil {
		t.Logf("Execute returned error: %v", err)
	}
	if r.ExitCode() != 0 {
		t.Errorf("expected exit code 0 on success, got %d", r.ExitCode())
	}
}

func TestRunnerExitCodeOnFailure(t *testing.T) {
	tools.ResetGlobals()
	tools.SetAutoExitMode(false)

	falsePath, err := exec.LookPath("false")
	if err != nil {
		t.Skipf("false command not found: %v", err)
	}

	tool := tools.Tool{
		Name:         "test-fail",
		Command:      falsePath,
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "")
	r.SetAutoExit(false)
	r.SetAutoExitTimeout(10 * time.Second)

	_ = r.Execute()
	if r.ExitCode() == 0 {
		t.Errorf("expected non-zero exit code on failure, got %d", r.ExitCode())
	}
}

func TestAutoExitTimeoutExpiry(t *testing.T) {
	tools.ResetGlobals()
	tools.SetAutoExitMode(false)

	sleepScript := filepath.Join(t.TempDir(), "sleep-long.sh")
	scriptContent := "#!/bin/sh\nsleep 300\n"
	if err := os.WriteFile(sleepScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write script: %v", err)
	}

	tool := tools.Tool{
		Name:         "timeout-test",
		Command:      sleepScript,
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "")
	r.SetAutoExit(true)
	r.SetAutoExitTimeout(500 * time.Millisecond)

	err := r.Execute()
	if err == nil {
		t.Error("expected timeout error for long-running subprocess")
	}
	if r.ExitCode() != 124 {
		t.Errorf("expected exit code 124 for timeout, got %d", r.ExitCode())
	}
}

func TestAutoExitModeDoesNotTimeoutWithoutFlag(t *testing.T) {
	tools.ResetGlobals()
	tools.SetAutoExitMode(false)

	sleepScript := filepath.Join(t.TempDir(), "sleep-short.sh")
	scriptContent := "#!/bin/sh\nexit 0\n"
	if err := os.WriteFile(sleepScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to write script: %v", err)
	}

	tool := tools.Tool{
		Name:         "no-timeout-test",
		Command:      sleepScript,
		AutoExitArgs: []string{},
		YoloArgs:     []string{},
	}
	r := New(tool, "test-model", "")
	r.SetAutoExit(false)

	err := r.Execute()
	if err != nil {
		t.Errorf("expected no error for quick command without auto-exit, got: %v", err)
	}
	if r.ExitCode() != 0 {
		t.Errorf("expected exit code 0, got %d", r.ExitCode())
	}
}
