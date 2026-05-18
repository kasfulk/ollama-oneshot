package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kasjfulk/ollama-oneshot/internal/config"
	"github.com/kasjfulk/ollama-oneshot/internal/docs"
	"github.com/kasjfulk/ollama-oneshot/internal/enhancer"
	"github.com/kasjfulk/ollama-oneshot/internal/prompt"
	"github.com/kasjfulk/ollama-oneshot/internal/runner"
	"github.com/kasjfulk/ollama-oneshot/internal/tools"
	"github.com/spf13/cobra"
)

var (
	flagPrompt    string
	flagTool      string
	flagModel     string
	flagDocs      []string
	flagDryRun    bool
	flagNoEnhance bool
	flagYoloMode  bool
	flagSystem    string
	flagProfile   string
	flagInclude   string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute an AI agent with enhanced prompt and project context",
	Long:  `Run orchestrates prompt enhancement, documentation injection, and agent launching in a single pipeline.`,
	RunE:  runExecute,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&flagPrompt, "prompt", "", "User prompt")
	runCmd.Flags().StringVar(&flagTool, "tool", "", "Agent tool to launch (claude, codex, opencode, etc.)")
	runCmd.Flags().StringVar(&flagModel, "model", "", "Execution model")
	runCmd.Flags().StringSliceVar(&flagDocs, "docs", []string{}, "Documentation files to inject (comma-separated)")
	runCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Preview final prompt without execution")
	runCmd.Flags().BoolVar(&flagNoEnhance, "no-enhance", false, "Skip prompt enhancement")
	runCmd.Flags().BoolVar(&flagYoloMode, "yolo-mode", false, "Auto-approve all tool permission prompts (bypassPermissions)")
	runCmd.Flags().StringVar(&flagSystem, "system", "", "Custom system prompt override")
	runCmd.Flags().StringVar(&flagProfile, "profile", "", "Load YAML profile preset")
	runCmd.Flags().StringVar(&flagInclude, "include", "", "Glob pattern for source files to include")

	runCmd.MarkFlagRequired("prompt")
}

func runExecute(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}
	cfg.ApplyFlags(flagModel, flagTool, flagNoEnhance, flagYoloMode)

	if cfg.YoloMode {
		tools.SetYoloMode(true)
	}

	fmt.Println("> Loading configuration...")

	var enhancedPrompt string
	if cfg.PromptEnhancement {
		fmt.Println("> Enhancing prompt...")
		client := enhancer.NewClient(cfg.OllamaURL(), cfg.PromptEnhancementModel)
		enhancedPrompt, err = client.Enhance(flagPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: prompt enhancement failed (%v), using original prompt\n", err)
			enhancedPrompt = flagPrompt
		} else {
			fmt.Println("> Prompt enhanced")
		}
	} else {
		enhancedPrompt = flagPrompt
	}

	fmt.Println("> Loading documentation...")
	var docFiles []string
	if len(flagDocs) > 0 {
		docFiles = flagDocs
	} else {
		docFiles, _ = docs.Discover(".")
	}
	docContent, err := docs.Load(docFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: docs loading failed (%v)\n", err)
	}

	if flagInclude != "" {
		patterns := strings.Split(flagInclude, ",")
		includedContent, err := docs.LoadGlob(patterns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: include glob failed (%v)\n", err)
		} else {
			docContent = docContent + "\n" + includedContent
		}
	}

	fmt.Println("> Assembling context...")
	assembler := prompt.NewAssembler(flagSystem)
	finalPrompt := assembler.Assemble(enhancedPrompt, flagPrompt, docContent)

	if flagDryRun {
		fmt.Println()
		fmt.Println("===== FINAL PROMPT =====")
		fmt.Println(finalPrompt)
		fmt.Println("========================")
		return nil
	}

	tool, ok := tools.Get(cfg.DefaultTool)
	if !ok {
		return fmt.Errorf("unknown tool: %s (available: %s)", cfg.DefaultTool, strings.Join(tools.List(), ", "))
	}

	fmt.Printf("> Launching %s...\n", tool.Name)
	fmt.Printf("> Using model %s\n", cfg.OllamaModel)

	execRunner := runner.New(tool, cfg.OllamaModel, finalPrompt)
	return execRunner.Execute()
}
