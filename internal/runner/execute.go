package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kasjfulk/ollama-oneshot/internal/tools"
)

type Runner struct {
	tool   tools.Tool
	model  string
	prompt string
}

func New(tool tools.Tool, model, prompt string) *Runner {
	return &Runner{
		tool:   tool,
		model:  model,
		prompt: prompt,
	}
}

func (r *Runner) Execute() error {
	cmdStr := r.tool.LaunchCommand(r.model)
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command for tool %s", r.tool.Name)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandler(cancel)

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", r.tool.Name, err)
	}

	go func() {
		io.WriteString(stdin, r.prompt)
		stdin.Close()
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s exited with error: %w", r.tool.Name, err)
	}

	fmt.Println()
	fmt.Println("✔ Done")
	return nil
}