# IMPLEMENTATION PLAN — `ollama-oneshoot`

> CLI Orchestration Tool: Prompt Enhancement · Documentation Injection · AI Agent Launcher · Multi-Tool Abstraction · Model Routing

---

## Overview

`ollama-oneshoot` adalah CLI orchestration layer di atas `ollama launch`. Bukan sekadar LLM wrapper, melainkan:

- **Launcher** — menjalankan agent tools via `ollama launch`
- **Prompt Enhancer** — mengubah prompt vague menjadi structured engineering task
- **Context Injector** — auto-inject dokumentasi project
- **Agent Orchestrator** — multi-tool, multi-model pipeline

**Stack:** Go · Cobra · Ollama API · godotenv · Viper

---

## Architecture

```
User Prompt
    ↓
Prompt Enhancement (Ollama API → deepseek-v4-flash)
    ↓
Documentation Injection (AGENT.md, PRD.md, DATABASE.md, ...)
    ↓
Prompt Assembly [SYSTEM + DOCS + ENHANCED + USER]
    ↓
Tool Launcher Resolver (Registry lookup)
    ↓
ollama launch <tool> --model <model>
    ↓
Agent Execution (stdout stream)
```

### Module Map

```
CLI (Cobra)
 ├── Config Loader        internal/config/
 ├── Prompt Enhancer      internal/enhancer/
 ├── Docs Loader          internal/docs/
 ├── Prompt Assembler     internal/prompt/
 ├── Tool Registry        internal/tools/
 ├── Ollama Client        internal/ollama/
 └── Runner               internal/runner/
```

---

## Project Structure

```
ollama-oneshoot/
├── cmd/
│   ├── root.go
│   └── run.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── enhancer/
│   │   ├── enhancer.go
│   │   └── prompts.go
│   │
│   ├── docs/
│   │   ├── loader.go
│   │   ├── discover.go
│   │   └── formatter.go
│   │
│   ├── ollama/
│   │   ├── client.go
│   │   ├── request.go
│   │   └── response.go
│   │
│   ├── prompt/
│   │   ├── assembler.go
│   │   └── templates.go
│   │
│   ├── tools/
│   │   ├── registry.go
│   │   ├── claude.go
│   │   ├── codex.go
│   │   ├── opencode.go
│   │   └── runner.go
│   │
│   └── runner/
│       ├── execute.go
│       ├── stream.go
│       └── process.go
│
├── .env
├── go.mod
├── main.go
└── README.md
```

---

## Environment Configuration

```bash
# .env

OLLAMA_HOST=127.0.0.1:11434

DEFAULT_MODEL=kimi-k2.6:cloud
DEFAULT_TOOL=claude

PROMPT_ENHANCEMENT=true
PROMPT_ENHANCEMENT_MODEL=deepseek-v4-flash
```

---

## CLI Interface

### Commands

| Command                    | Description              |
| -------------------------- | ------------------------ |
| `ollama-oneshoot`          | root / help              |
| `ollama-oneshoot run`      | main execution command   |

### Flags

| Flag           | Type     | Default           | Description                      |
| -------------- | -------- | ----------------- | -------------------------------- |
| `--prompt`     | string   | —                 | User prompt                      |
| `--tool`       | string   | `claude`          | Selected agent tool              |
| `--model`      | string   | `kimi-k2.6:cloud` | Execution model                  |
| `--docs`       | []string | auto-discover     | Documentation files to inject    |
| `--dry-run`    | bool     | false             | Preview final prompt, no execute |
| `--no-enhance` | bool     | false             | Skip prompt enhancement          |
| `--system`     | string   | default template  | Custom system prompt override    |
| `--profile`    | string   | —                 | Load YAML profile preset         |
| `--include`    | string   | —                 | Glob pattern for source files    |

### Example Usage

```bash
# Basic
ollama-oneshoot run --prompt "build ERP application"

# With specific tool and model
ollama-oneshoot run \
  --prompt "build currency exchange app" \
  --tool claude \
  --model kimi-k2.6:cloud \
  --docs AGENT.md,PRD.md

# Dry run (preview prompt only)
ollama-oneshoot run \
  --prompt "build CRM" \
  --dry-run

# Using profile
ollama-oneshoot run --prompt "optimize database queries" --profile backend
```

---

## Tool Registry

### Supported Tools

| Tool        | Launch Command              |
| ----------- | --------------------------- |
| `claude`    | `ollama launch claude`      |
| `claude-code` | `ollama launch claude`    |
| `codex`     | `ollama launch codex`       |
| `codex-app` | `ollama launch codex-app`   |
| `opencode`  | `ollama launch opencode`    |
| `openclaw`  | `ollama launch openclaw`    |
| `hermes`    | `ollama launch hermes`      |

---

## Implementation Phases

---

### Phase 1 — MVP Foundation

**Timeline: 6 Weeks**

---

#### Week 1 — Project Bootstrap + Cobra CLI

**Tasks:**

1. Initialize Go module:
   ```bash
   go mod init github.com/yourname/ollama-oneshoot
   ```

2. Install dependencies:
   ```bash
   go get github.com/spf13/cobra
   go get github.com/joho/godotenv
   go get github.com/spf13/viper
   ```

3. Scaffold base directory structure (`cmd/`, `internal/`, `main.go`, `.env`)

4. Implement `cmd/root.go` — root command with version info

5. Implement `cmd/run.go` — run command with all flags

**Deliverables:** `main.go`, `cmd/root.go`, `cmd/run.go`

---

#### Week 1 (cont.) — Environment Configuration

**Goal:** Flexible config from `.env` + flags + defaults

**Implementation:** `internal/config/config.go`

```go
type Config struct {
    OllamaHost             string
    OllamaModel            string
    DefaultTool            string
    PromptEnhancement      bool
    PromptEnhancementModel string
}
```

Config priority: CLI flags → environment variables → hardcoded defaults

**Deliverables:** `internal/config/config.go`

---

#### Week 2 — Ollama API Client

**Goal:** HTTP client untuk communicate dengan Ollama API

**Endpoint:** `POST http://<OLLAMA_HOST>/api/generate`

**Features:**
- Send prompt to model
- Parse JSON response
- Support streaming (chunked response)
- Configurable timeout

**Deliverables:** `internal/ollama/client.go`, `request.go`, `response.go`

---

#### Week 2 (cont.) — Tool Registry

**Goal:** Map tool names → launch commands

**Deliverables:** `internal/tools/registry.go`, `claude.go`, `codex.go`, `runner.go`

---

#### Week 3 — Prompt Enhancement

**Goal:** Transform vague prompts into structured engineering tasks

**Flow:**
```
user prompt (vague)
    ↓
Enhancement Model (deepseek-v4-flash via Ollama)
    ↓
structured prompt (implementation-ready)
```

**Enhancement Strategy:**
- Short prompt (< 50 words) → aggressive enhancement
- Long prompt (≥ 50 words) → minimal normalization only

**System Prompt for Enhancer:**
```
You are an expert software architect.
Convert vague software requests into:
- implementation plans
- architecture
- engineering tasks
- scalable structure
Focus on: backend, frontend, database, deployment, security, maintainability.
```

**Deliverables:** `internal/enhancer/enhancer.go`, `prompts.go`

---

#### Week 4 — Documentation Injection

**Goal:** Auto-inject project context into prompt

**Auto-discovered files (in order):**

| File              | Purpose                    |
| ----------------- | -------------------------- |
| `AGENT.md`        | Agent rules & behavior     |
| `PRD.md`          | Product requirements       |
| `DATABASE.md`     | Schema & data model        |
| `ARCHITECTURE.md` | System architecture        |
| `TASKS.md`        | Current task context       |

**Manual include:** `--docs=AGENT.md,PRD.md`

**Injection Format:**
```
<PROJECT_DOCUMENTATION>

[FILE: AGENT.md]
<content>

[FILE: PRD.md]
<content>

</PROJECT_DOCUMENTATION>
```

**Deliverables:** `internal/docs/loader.go`, `discover.go`, `formatter.go`

---

#### Week 4 (cont.) — Prompt Assembly

**Goal:** Assemble final context from all components

**Final Prompt Structure:**
```
[SYSTEM]

[DOCUMENTATION]

[ENHANCED PROMPT]

[USER PROMPT]
```

**Deliverables:** `internal/prompt/assembler.go`, `templates.go`

---

#### Week 5 — Tool Runner

**Goal:** Spawn agent tool process, pipe prompt, stream output

**Execution:**
```bash
ollama launch claude --model kimi-k2.6:cloud
```

**Responsibilities:**
- Spawn OS process via `os/exec`
- Pipe assembled prompt to stdin
- Stream stdout to terminal
- Handle SIGINT / cancellation
- Exit code propagation

**Deliverables:** `internal/runner/execute.go`, `stream.go`, `process.go`

---

#### Week 6 — Dry Run Mode + Polish

**Goal:** Preview final assembled prompt without execution

**Output:**
```
===== FINAL PROMPT =====
[SYSTEM]
...
[DOCUMENTATION]
...
[ENHANCED PROMPT]
...
[USER PROMPT]
...
========================
```

**Polish tasks:**
- Error handling & user-friendly messages
- `--no-enhance` flag implementation
- Progress output (`> Enhancing...`, `> Loading docs...`, etc.)
- README documentation

**Deliverables:** Full MVP working end-to-end

---

### Phase 2 — Quality Improvements

**Timeline: Post-MVP (2–3 weeks)**

---

#### Phase 2.1 — Streaming UI

**Goal:** Improve terminal UX with real-time feedback

**Library:** `github.com/charmbracelet/bubbletea`

**Features:**
- Loading spinner during enhancement
- Streaming token display
- Progress indicator per step

**Deliverables:** `internal/ui/spinner.go`, `stream.go`

---

#### Phase 2.2 — `.oneshootignore`

**Goal:** Exclude irrelevant files from auto-discovery and `--include`

**Format (similar to .gitignore):**
```
node_modules
vendor
dist
*.lock
*.sum
```

**Deliverables:** `internal/ignore/parser.go`

---

#### Phase 2.3 — Context Compression

**Goal:** Prevent context window overflow when docs are large

**Strategy:**
1. Calculate total token estimate
2. If docs exceed threshold → summarize via Ollama
3. Inject summary instead of raw content
4. Always keep AGENT.md full (critical rules)

**Deliverables:** `internal/context/compressor.go`, `summarizer.go`

---

#### Phase 2.4 — File Inclusion

**Goal:** Inject source code directly into context

**Usage:**
```bash
--include=src/**/*.go
--include=handlers/**/*.go,models/**/*.go
```

**Deliverables:** `internal/include/glob.go`, `reader.go`

---

### Phase 3 — Advanced Agentic Features

**Timeline: 4–6 weeks after Phase 2**

---

#### Phase 3.1 — Profiles

**Goal:** Reusable preset configurations via YAML

**Usage:**
```bash
ollama-oneshoot run --prompt "optimize auth flow" --profile backend
```

**Profile Format (`backend.yaml`):**
```yaml
tool: claude
model: kimi-k2.6:cloud

docs:
  - AGENT.md
  - DATABASE.md

enhancement: true
system: "Focus on Go backend, PostgreSQL, and REST API patterns."
```

**Deliverables:** `internal/profile/loader.go`, `schema.go`

---

#### Phase 3.2 — Multi-Agent Pipeline

**Goal:** Planner → Executor → Reviewer pipeline

**Usage:**
```bash
ollama-oneshoot run \
  --planner hermes \
  --executor claude \
  --reviewer codex \
  --prompt "build auth service"
```

**Flow:**
```
Planner (hermes)  → generates implementation plan
      ↓
Executor (claude) → implements based on plan
      ↓
Reviewer (codex)  → reviews and validates output
```

**Deliverables:** `internal/agents/planner.go`, `executor.go`, `reviewer.go`

---

#### Phase 3.3 — Auto Tool Routing

**Goal:** Automatically select best tool based on task type

**Routing Rules:**

| Task Keywords         | Routed Tool |
| --------------------- | ----------- |
| frontend, UI, React   | `codex-app` |
| backend, API, database | `claude`   |
| autonomous, plan, research | `hermes` |
| terminal, script, CLI | `opencode`  |
| review, audit         | `codex`     |

**Deliverables:** `internal/router/task_classifier.go`, `tool_router.go`

---

### Phase 4 — Production Ready

**Timeline: Future / Ongoing**

| Feature               | Description                                    |
| --------------------- | ---------------------------------------------- |
| Telemetry             | Usage metrics, latency tracking                |
| Crash Recovery        | Auto-restart on agent failure                  |
| Plugin System         | Custom tool adapters via interface             |
| MCP Integration       | Model Context Protocol support                 |
| Remote Execution      | Execute against remote Ollama hosts            |
| Distributed Agents    | Parallel multi-agent task execution            |

---

## MVP Timeline Summary

| Week | Goal                                      | Key Deliverables                             |
| ---- | ----------------------------------------- | -------------------------------------------- |
| 1    | Bootstrap + Cobra CLI + Config            | `main.go`, `cmd/`, `internal/config/`        |
| 2    | Ollama Client + Tool Registry             | `internal/ollama/`, `internal/tools/`        |
| 3    | Prompt Enhancement                        | `internal/enhancer/`                         |
| 4    | Docs Injection + Prompt Assembly          | `internal/docs/`, `internal/prompt/`         |
| 5    | Tool Runner                               | `internal/runner/`                           |
| 6    | Dry Run + Polish + README                 | Full MVP, end-to-end working                 |

---

## Expected Final UX (MVP)

```bash
$ ollama-oneshoot run \
    --prompt "build ERP application" \
    --tool claude \
    --model kimi-k2.6:cloud \
    --docs AGENT.md,PRD.md

> Enhancing prompt...
> Loading documentation...
> Assembling context...
> Launching claude...
> Using model kimi-k2.6:cloud

✔ Ready
[agent output streams here...]
```

### Dry Run Output

```bash
$ ollama-oneshoot run --prompt "build ERP app" --dry-run

===== FINAL PROMPT =====

[SYSTEM]
You are an expert software engineer...

[DOCUMENTATION]
[FILE: AGENT.md]
...

[ENHANCED PROMPT]
Build a production-ready ERP application with:
- Multi-module architecture (Finance, HR, Inventory, CRM)
- Role-based access control
- REST API with OpenAPI spec
- PostgreSQL with migration support
- Docker Compose for local dev
- CI/CD pipeline template

[USER PROMPT]
build ERP application

========================
```

---

## MVP Must-Have Checklist

- [ ] Cobra CLI with all flags
- [ ] `.env` config loading with fallbacks
- [ ] Ollama `/api/generate` client
- [ ] Prompt enhancement via local model
- [ ] Auto-discover docs (AGENT.md, PRD.md, DATABASE.md)
- [ ] Manual docs include via `--docs`
- [ ] Prompt assembly [SYSTEM + DOCS + ENHANCED + USER]
- [ ] Tool registry with `ollama launch` mapping
- [ ] Process runner with stdin pipe + stdout stream
- [ ] Dry-run mode
- [ ] `--no-enhance` flag
- [ ] README with install and usage guide

## Future Enhancement Checklist

- [ ] Streaming UI (bubbletea)
- [ ] `.oneshootignore` support
- [ ] Context compression
- [ ] Source file inclusion (`--include`)
- [ ] YAML profiles
- [ ] Multi-agent pipeline
- [ ] Auto tool routing
- [ ] MCP integration
- [ ] Plugin system

---

*Generated for `ollama-oneshoot` v0.1.0 MVP*
