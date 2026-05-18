package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "ollama-oneshot",
	Short: "CLI orchestration layer for ollama launch — prompt enhancement, docs injection, agent launching",
	Long: `ollama-oneshot is a CLI orchestration tool that sits on top of ollama launch.
It enhances prompts, injects project documentation, and launches AI agent tools
with structured, context-rich instructions.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Help for ollama-oneshot")
}