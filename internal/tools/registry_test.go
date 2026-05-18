package tools

import (
	"sync"
	"testing"
)

func TestSetAutoExitMode(t *testing.T) {
	SetAutoExitMode(true)
	if !AutoExitMode() {
		t.Error("expected AutoExitMode to be true after SetAutoExitMode(true)")
	}
	SetAutoExitMode(false)
	if AutoExitMode() {
		t.Error("expected AutoExitMode to be false after SetAutoExitMode(false)")
	}
}

func TestAutoExitModeConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val bool) {
			defer wg.Done()
			SetAutoExitMode(val)
		}(i%2 == 0)
	}
	wg.Wait()
}

func TestAutoExitModeDefault(t *testing.T) {
	SetAutoExitMode(false)
	ResetGlobals()
	if AutoExitMode() {
		t.Error("expected AutoExitMode to default to false")
	}
}

func TestLogAutoExitActivation(t *testing.T) {
	SetAutoExitMode(true)
	tool := Tool{
		Name:         "test-tool",
		Command:      "echo",
		AutoExitArgs: []string{"--exit"},
	}
	_, ok := Get(tool.Name)
	if ok {
		t.Error("expected test-tool not to be registered yet")
	}
}

func TestRegistryAutoExitArgs(t *testing.T) {
	ResetGlobals()
	registry["test-autoexit"] = Tool{
		Name:         "test-autoexit",
		Command:      "echo",
		AutoExitArgs: []string{"--", "--exit"},
		YoloArgs:     []string{"--", "--approve"},
	}
	tool, ok := Get("test-autoexit")
	if !ok {
		t.Fatal("expected tool to be found")
	}
	if len(tool.AutoExitArgs) != 2 {
		t.Errorf("expected 2 AutoExitArgs, got %d", len(tool.AutoExitArgs))
	}
	if tool.AutoExitArgs[1] != "--exit" {
		t.Errorf("expected AutoExitArgs[1] to be --exit, got %s", tool.AutoExitArgs[1])
	}
}