# ollama-oneshot

CLI orchestration layer for `ollama launch` — prompt enhancement, documentation injection, and AI agent launching in a single pipeline.

## What it does

```
User Prompt
    ↓
Prompt Enhancement  (Ollama → deepseek-v4-flash)
    ↓
Documentation Injection  (AGENT.md, PRD.md, DATABASE.md, …)
    ↓
Prompt Assembly  [SYSTEM + DOCS + ENHANCED + USER]
    ↓
ollama launch <tool> --model <model>
    ↓
Agent output streams to terminal
```

## Requirements

- [Go 1.21+](https://go.dev/dl/)
- [Ollama](https://ollama.com) running locally

## Install

```bash
git clone https://github.com/kasjfulk/ollama-oneshot
cd ollama-oneshot
bash install.sh
```

Or with Make:

```bash
make install
```

The binary is installed to `~/.local/bin/ollama-oneshot`. Make sure that directory is in your `PATH`.

## Configuration

Copy or edit `.env` in the project root:

```bash
OLLAMA_HOST=127.0.0.1:11434

DEFAULT_MODEL=kimi-k2.6:cloud
DEFAULT_TOOL=claude

PROMPT_ENHANCEMENT=true
PROMPT_ENHANCEMENT_MODEL=deepseek-v4-flash
```

Config priority: **CLI flags → environment variables → defaults**

## Usage

### Basic

```bash
ollama-oneshot run --prompt "build ERP application"
```

### With explicit tool and model

```bash
ollama-oneshot run \
  --prompt "build currency exchange app" \
  --tool claude \
  --model kimi-k2.6:cloud \
  --docs AGENT.md,PRD.md
```

### Dry run — preview assembled prompt, no execution

```bash
ollama-oneshot run \
  --prompt "build CRM" \
  --dry-run
```

### Skip prompt enhancement

```bash
ollama-oneshot run \
  --prompt "refactor auth middleware" \
  --no-enhance
```

### Include source files in context

```bash
ollama-oneshot run \
  --prompt "review this code" \
  --include "internal/**/*.go"
```

## Flags

| Flag           | Type       | Default            | Description                            |
| -------------- | ---------- | ------------------ | -------------------------------------- |
| `--prompt`     | string     | **required**       | The task or request                    |
| `--tool`       | string     | `claude`           | Agent tool to launch                   |
| `--model`      | string     | `kimi-k2.6:cloud`  | Execution model                        |
| `--docs`       | []string   | auto-discover      | Docs to inject (comma-separated)       |
| `--include`    | string     | —                  | Glob pattern for source files          |
| `--dry-run`    | bool       | false              | Preview final prompt, skip execution   |
| `--no-enhance` | bool       | false              | Skip prompt enhancement                |
| `--system`     | string     | default template   | Override system prompt                 |
| `--profile`    | string     | —                  | Load a YAML profile preset             |

## Supported Tools

| Tool          | Launch Command              |
| ------------- | --------------------------- |
| `claude`      | `ollama launch claude`      |
| `claude-code` | `ollama launch claude`      |
| `codex`       | `ollama launch codex`       |
| `codex-app`   | `ollama launch codex-app`   |
| `opencode`    | `ollama launch opencode`    |
| `openclaw`    | `ollama launch openclaw`    |
| `hermes`      | `ollama launch hermes`      |

## Documentation Auto-Discovery

When `--docs` is not provided, the following files are auto-discovered from the current directory:

- `AGENT.md`
- `PRD.md`
- `DATABASE.md`
- `ARCHITECTURE.md`
- `TASKS.md`

## Prompt Enhancement

Short prompts (< 50 words) are aggressively expanded into structured implementation specs. Longer prompts are normalized and structured without changing the detail level.

Disable with `--no-enhance` or set `PROMPT_ENHANCEMENT=false` in `.env`.

## Make Targets

```bash
make build      # compile to ./dist/
make install    # build + install to ~/.local/bin
make uninstall  # remove installed binary
make test       # run tests
make lint       # go vet
make clean      # remove ./dist/
make run ARGS="run --prompt 'hello' --dry-run"
make help       # list all targets
```

## Example Output

```
> Loading configuration...
> Enhancing prompt...
> Prompt enhanced
> Loading documentation...
> Assembling context...
> Launching claude...
> Using model kimi-k2.6:cloud

[agent output streams here...]

✔ Done
```
