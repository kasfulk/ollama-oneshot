# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI orchestration layer for `ollama launch` — prompt enhancement, documentation injection, and AI agent launching in a single pipeline. Written in Go using cobra/viper.

## Build & Test

```bash
make build      # compile to ./dist/ollama-oneshot
make install    # build + install to ~/.local/bin
make test       # go test ./...
make lint       # go vet ./...
make run ARGS="run --prompt 'hello' --dry-run"  # build + run with args
make deps       # go mod tidy && go mod download
make clean      # remove ./dist/
```

## Architecture

```
main.go → cmd/ (cobra commands)
              ├── root.go       — base command, version
              └── run.go        — pipeline orchestrator (the core command)
                       │
                       ▼
              internal/         — pipeline stages as packages
              ├── config/       — viper + godotenv config loading (env → CLI flags override)
              ├── enhancer/     — prompt enhancement via Ollama LLM
              ├── docs/         — auto-discover & load project docs (AGENT.md, PRD.md, etc.)
              ├── prompt/       — final prompt assembly ([SYSTEM] + [DOCS] + [ENHANCED] + [USER])
              ├── tools/        — tool registry (claude, codex, opencode, hermes, etc.)
              ├── runner/       — subprocess execution with stdin pipe for prompt
              └── ollama/       — HTTP client for Ollama API (/api/generate)
```

### Pipeline Flow

1. **Config** loads from `.env` → environment → CLI flag overrides
2. **Enhancer** sends short prompts (<50 words) to Ollama for expansion into structured specs
3. **Docs** auto-discovers `{AGENT.md,PRD.md,DATABASE.md,ARCHITECTURE.md,TASKS.md}` from CWD, plus optional `--include` glob
4. **Prompt assembler** builds the final context: `[SYSTEM]` + `[DOCUMENTATION]` + `[ENHANCED PROMPT]` + `[USER PROMPT]`
5. **Runner** spawns `ollama launch <tool> --model <model>` as a subprocess and pipes the assembled prompt via stdin

### Tool Registry

Tools are registered in `internal/tools/registry.go` as name→command mappings. Each maps to `ollama launch <name>`. To add a new tool, add an entry to the `registry` map.

### Key Dependencies

- `spf13/cobra` — CLI commands and flags
- `spf13/viper` + `joho/godotenv` — configuration (env file + env vars + defaults)
- `fsnotify` — pulled in by viper (not used directly)

### Config Priority

CLI flags > environment variables > `.env` file > viper defaults (in `internal/config/config.go`)
