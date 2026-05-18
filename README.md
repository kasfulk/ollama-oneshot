# ollama-oneshot

CLI orchestration layer for `ollama launch` — prompt enhancement, documentation injection, model validation, and AI agent launching in a single pipeline.

## What it does

```
User Prompt
    ↓
Prompt Enhancement  (Ollama → deepseek-v4-flash:cloud)
    ↓
Documentation Injection  (AGENT.md, PRD.md, DATABASE.md, …)
    ↓
Model Validation  (checks model exists in Ollama)
    ↓
Prompt Assembly  [SYSTEM + DOCS + ENHANCED + USER]
    ↓
ollama launch <tool> --model <model>
    ↓
Agent output streams to terminal
    ↓
Auto-Exit (optional) — exit or timeout after completion
```

## Requirements

- [Go 1.21+](https://go.dev/dl/)
- [Ollama](https://ollama.com) running locally
- Target model pulled in Ollama (e.g., `ollama pull glm-5.1:cloud`)

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

DEFAULT_MODEL=glm-5.1:cloud
DEFAULT_TOOL=claude

PROMPT_ENHANCEMENT=true
PROMPT_ENHANCEMENT_MODEL=deepseek-v4-flash:cloud

YOLO_MODE=false
AUTO_EXIT=false
```

Config priority: **CLI flags → environment variables → `.env` → defaults**

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
  --model glm-5.1:cloud \
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

### YOLO mode — auto-approve tool permission prompts

```bash
ollama-oneshot run \
  --prompt "build a full-stack app" \
  --yolo-mode
```

When enabled, the tool appends per-approval flags:

| Tool          | YOLO argument                           |
| ------------- | --------------------------------------- |
| `claude`      | `--dangerously-skip-permissions`        |
| `codex`       | `--full-auto`                           |
| `opencode`    | _(none — no built-in auto-approve)_     |

### Auto-exit — automatically exit after task completion

```bash
ollama-oneshot run \
  --prompt "refactor this module" \
  --auto-exit
```

- Automatically exits `ollama launch` after the task completes.
- Uses a 30-minute timeout to prevent hanging processes.
- On timeout, exits with code **124**.
- When combined with `--yolo-mode`, enables fully unattended operation.

Each tool appends a per-tool auto-exit argument (`-p` for claude, `exec` for codex, `run` for opencode).

## Flags

| Flag           | Type       | Default            | Description                            |
| -------------- | ---------- | ------------------ | -------------------------------------- |
| `--prompt`     | string     | **required**       | The task or request                    |
| `--tool`       | string     | `claude`           | Agent tool to launch                   |
| `--model`      | string     | `glm-5.1:cloud`    | Execution model                        |
| `--docs`       | []string   | auto-discover      | Docs to inject (comma-separated)       |
| `--include`    | string     | —                  | Glob pattern for source files          |
| `--dry-run`    | bool       | false              | Preview final prompt, skip execution   |
| `--no-enhance` | bool       | false              | Skip prompt enhancement                |
| `--yolo-mode`  | bool       | false              | Auto-approve tool permission prompts   |
| `--auto-exit`  | bool       | false              | Exit after task completion             |
| `--system`     | string     | default template   | Override system prompt                 |
| `--profile`    | string     | —                  | Load a YAML profile preset             |

## Supported Tools

| Tool          | Launch Command              | YOLO args                          | Auto-exit args |
| ------------- | --------------------------- | ---------------------------------- | -------------- |
| `claude`      | `ollama launch claude`      | `--dangerously-skip-permissions`   | `-p`           |
| `claude-code` | `ollama launch claude`      | `--dangerously-skip-permissions`   | `-p`           |
| `codex`       | `ollama launch codex`       | `--full-auto`                      | `exec`         |
| `codex-app`   | `ollama launch codex-app`   | `--full-auto`                      | `exec`         |
| `opencode`    | `ollama launch opencode`    | —                                  | `run`          |
| `openclaw`    | `ollama launch openclaw`    | —                                  | `run`          |
| `hermes`      | `ollama launch hermes`      | —                                  | `run`          |

## Model Validation

Before launching a tool, `ollama-oneshot` checks that the specified model exists in your local Ollama instance using the `/api/tags` endpoint. If the model is not found, execution is aborted with a message indicating the required `ollama pull` command.

This check can be disabled implicitly — validation warnings do not block execution.

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
> Validating model...
> Launching claude...
> Using model glm-5.1:cloud

[agent output streams here...]

✔ Done
```
