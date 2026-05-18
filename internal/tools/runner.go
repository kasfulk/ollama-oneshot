package tools

import (
	"fmt"
	"runtime"
	"strings"
)

func (t Tool) LaunchCommand(model string) string {
	cmd := t.Command + " --model " + model
	if YoloMode() {
		if len(t.YoloArgs) == 0 {
			yoloLog.Printf("no yolo args for tool=%q — auto-approval may not take effect", t.Name)
		}
		for _, arg := range t.YoloArgs {
			cmd += " " + arg
		}
		caller := callerContext()
		LogYoloApproval(t.Name, caller)
	}
	if AutoExitMode() {
		if len(t.AutoExitArgs) == 0 {
			autoExitLog.Printf("no auto-exit args for tool=%q — auto-exit may not take effect", t.Name)
		}
		for _, arg := range t.AutoExitArgs {
			cmd += " " + arg
		}
		caller := callerContext()
		LogAutoExitActivation(t.Name, caller)
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
