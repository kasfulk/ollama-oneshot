package tools

import (
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type Tool struct {
	Name         string
	Command      string
	YoloArgs     []string
	AutoExitArgs []string
}

var registry = map[string]Tool{}

var (
	yoloMode    bool
	yoloMu      sync.RWMutex
	yoloLog     = log.New(os.Stderr, "[YOLO] ", log.Lmsgprefix)
	autoExit    bool
	autoExitMu  sync.RWMutex
	autoExitLog = log.New(os.Stderr, "[AUTO-EXIT] ", log.Lmsgprefix)
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

func SetAutoExitMode(enabled bool) {
	autoExitMu.Lock()
	defer autoExitMu.Unlock()
	autoExit = enabled
	if enabled {
		autoExitLog.Printf("mode enabled at %s", time.Now().Format(time.RFC3339))
	}
}

func AutoExitMode() bool {
	autoExitMu.RLock()
	defer autoExitMu.RUnlock()
	return autoExit
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

func LogAutoExitActivation(toolName, callerContext string) {
	autoExitLog.Printf("activated for tool=%q caller=%q at=%s", toolName, callerContext, time.Now().Format(time.RFC3339))
}

func ResetGlobals() {
	yoloMu.Lock()
	yoloMode = false
	yoloMu.Unlock()

	autoExitMu.Lock()
	autoExit = false
	autoExitMu.Unlock()
}
