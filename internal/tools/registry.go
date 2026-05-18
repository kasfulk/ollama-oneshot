package tools

import (
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type Tool struct {
	Name    string
	Command string
}

var registry = map[string]Tool{
	"claude":      {Name: "claude", Command: "ollama launch claude"},
	"claude-code": {Name: "claude-code", Command: "ollama launch claude"},
	"codex":       {Name: "codex", Command: "ollama launch codex"},
	"codex-app":   {Name: "codex-app", Command: "ollama launch codex-app"},
	"opencode":    {Name: "opencode", Command: "ollama launch opencode"},
	"openclaw":    {Name: "openclaw", Command: "ollama launch openclaw"},
	"hermes":      {Name: "hermes", Command: "ollama launch hermes"},
}

var (
	yoloMode bool
	yoloMu   sync.RWMutex
	yoloLog  = log.New(os.Stderr, "[YOLO] ", log.Lmsgprefix)
)

func SetYoloMode(enabled bool) {
	yoloMu.Lock()
	defer yoloMu.Unlock()
	yoloMode = enabled
	if enabled {
		yoloLog.Printf("mode enabled at %s", time.Now().Format(time.RFC3339))
	}
}

func YoloMode() bool {
	yoloMu.RLock()
	defer yoloMu.RUnlock()
	return yoloMode
}

func Get(name string) (Tool, bool) {
	t, ok := registry[name]
	return t, ok
}

func List() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func LogYoloApproval(toolName, callerContext string) {
	yoloLog.Printf("auto-approved tool=%q caller=%q at=%s", toolName, callerContext, time.Now().Format(time.RFC3339))
}
