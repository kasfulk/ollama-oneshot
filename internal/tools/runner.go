package tools

import (
	"fmt"
	"runtime"
	"strings"
)

func (t Tool) LaunchCommand(model string) string {
	cmd := t.Command + " --model " + model
	if YoloMode() {
		cmd += " --permission-mode bypassPermissions"
		caller := callerContext()
		LogYoloApproval(t.Name, caller)
	}
	return cmd
}

func callerContext() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	return fmt.Sprintf("%s:%d", file, line)
}
