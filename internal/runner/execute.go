package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/kasjfulk/ollama-oneshot/internal/tools"
)

const defaultAutoExitTimeout = 30 * time.Minute

type Runner struct {
	tool            tools.Tool
	model           string
	prompt          string
	exitCode        int
	autoExit        bool
	autoExitTimeout time.Duration
}

func New(tool tools.Tool, model, prompt string) *Runner {
	return &Runner{
		tool:            tool,
		model:           model,
		prompt:          prompt,
		exitCode:        0,
		autoExit:        false,
		autoExitTimeout: defaultAutoExitTimeout,
	}
}

func (r *Runner) ExitCode() int {
	return r.exitCode
}

func (r *Runner) SetAutoExit(enabled bool) {
	r.autoExit = enabled
}

func (r *Runner) SetAutoExitTimeout(d time.Duration) {
	r.autoExitTimeout = d
}

func (r *Runner) Execute() error {
	cmdStr := r.tool.LaunchCommand(r.model)
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		r.exitCode = 1
		return fmt.Errorf("empty command for tool %s", r.tool.Name)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if r.autoExit {
		ctx, cancel = context.WithTimeout(context.Background(), r.autoExitTimeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	sh := setupSignalHandler(cancel)
	defer sh.Cleanup()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.WaitDelay = time.Second
	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGKILL)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		r.exitCode = 1
		return fmt.Errorf("stdin pipe: %w", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		r.exitCode = 1
		return fmt.Errorf("start %s: %w", r.tool.Name, err)
	}

	stdinDone := make(chan struct{})
	go func() {
		io.WriteString(stdin, r.prompt)
		stdin.Close()
		close(stdinDone)
	}()

	if err := cmd.Wait(); err != nil {
		<-stdinDone
		if ctx.Err() == context.DeadlineExceeded {
			r.exitCode = 124
			return fmt.Errorf("%s timed out after %v (auto-exit deadline)", r.tool.Name, r.autoExitTimeout)
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			r.exitCode = exitErr.ExitCode()
		} else {
			r.exitCode = 1
		}
		return fmt.Errorf("%s exited with error: %w", r.tool.Name, err)
	}

	<-stdinDone

	fmt.Println()
	fmt.Println("✔ Done")
	return nil
}
