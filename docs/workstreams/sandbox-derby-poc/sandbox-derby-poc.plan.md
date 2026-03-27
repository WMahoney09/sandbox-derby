# Sandbox Derby PoC

## Overview

Build a proof-of-concept for Sandbox Derby: containerized Claude agent workspaces with configurable loadouts, two execution modes (drive and coast), and a derby system for structured comparison across sandboxes. The PoC uses Docker Compose for single-sandbox operations and a Go CLI for derby orchestration.

## Notes

**Out of scope:**
- Charm TUI (MVP)
- LLM-as-judge evaluation (roadmap)
- Agent self-retros (roadmap)
- Auto-filing issues against tool repos (roadmap)
- Non-Anthropic model support (future)
- Worktrees, PRs, or repo forking as publish mechanisms (MVP)

**Key decisions:**
- Approach C: Compose for sandbox lifecycle, Go for derby orchestration
- Loadouts are volume-mounted read-only to a staging path, then copied into `/home/agent/.claude/` by the entrypoint so Claude Code can write freely to `~/.claude/` at runtime
- Workspace is a git repo cloned inside the container via TARGET_REPO env var
- Course is always a local markdown file, mounted read-only to a staging path and copied into the workspace by the entrypoint (so the agent can check off TODOs as it works)
- Model selection lives in the loadout (settings.json or agent definitions), not at the derby level
- Derby config is YAML
- Base image adapted from kubesat: `debian:bookworm-slim` + Claude Code CLI, gh, git, Node.js, Python, non-root `agent` user (no kubectl)

**Assumptions:**
- All courses target git-based projects (workspace = cloned repo with commit history)
- ANTHROPIC_API_KEY and GITHUB_TOKEN provided via env vars / .env file
- The operator has Docker and Go installed locally

## Progress
- [ ] Phase 1: Foundation
- [ ] Phase 2: Drive
- [ ] Phase 3: Coast
- [ ] Phase 4: Derby

---

## Phase 1: Foundation

Scaffold the project, build the base image, and define the canonical loadout structure. At the end of this phase, you can build the image and it runs.

### Step 1.1: Project structure

#### Task 1.1.1: Initialize Go module
Create `go.mod` with module path `github.com/WMahoney09/sandbox-derby`.

#### Task 1.1.2: Create directory layout
```
sandbox/
  Dockerfile
  entrypoint-common.sh   # shared setup: env validation, git identity, loadout/course copy-in
  entrypoint-drive.sh
  entrypoint-coast.sh
docker-compose.yml
.env.example
loadouts/
  bare/
    .gitkeep       # empty loadout — baseline for comparison
  example/         # reference loadout — mirrors .claude/ structure
    CLAUDE.md      # example agent guidance
    settings.json  # example settings
courses/
  example.md       # reference course
derby/
  derby.yaml.example
cmd/
  derby/
    main.go        # Go CLI entrypoint
internal/
  derby/           # derby orchestration logic
docs/
```

#### Task 1.1.3: Define canonical loadout structure
A loadout directory mirrors the `.claude/` structure. It is mounted read-only to a staging path (`/home/agent/loadout/`) and the entrypoint copies its contents into `/home/agent/.claude/` at startup. This lets Claude Code write freely to `~/.claude/` at runtime while keeping the loadout source clean.
```
loadout-name/
  CLAUDE.md       # agent guidance and config
  settings.json   # Claude Code settings
  skills/         # skill definitions
  agents/         # custom sub-agent definitions
  teams/          # team definitions
```
All slots are optional. A bare loadout is an empty directory (with `.gitkeep` for git tracking).

#### Task 1.1.4: Write entrypoint-common.sh
Shared setup script sourced by both drive and coast entrypoints. Responsibilities:
1. Validate `ANTHROPIC_API_KEY` is set (exit with error if not)
2. Configure git identity from `GIT_USER_NAME` / `GIT_USER_EMAIL` env vars (default to `Sandbox Derby Agent` / `sandbox-derby[bot]@noreply.github.com`)
3. If `/home/agent/loadout/` exists and is non-empty, copy its contents into `/home/agent/.claude/`
4. If `/home/agent/course/` exists, copy the course file into `/home/agent/workspace/`

All copy operations are conditional — the script handles whatever is present without erroring on what's absent.

### Step 1.2: Base image

#### Task 1.2.1: Write Dockerfile
Adapted from kubesat. `debian:bookworm-slim` base. Install git, curl, ca-certificates, gnupg, python3, nodejs, npm, gh CLI. Create non-root `agent` user. Install Claude Code CLI. No kubectl. No build-time skill cloning (loadouts are mounted, not baked). Conventional image name: `sandbox-derby`.

#### Task 1.2.2: Write .env.example
Document required env vars: `ANTHROPIC_API_KEY`, `GITHUB_TOKEN`, `TARGET_REPO` (optional), `GIT_USER_NAME` (optional), `GIT_USER_EMAIL` (optional).

### Step 1.3: Verify foundation

#### Task 1.3.1: Build image and confirm Claude Code is available
Build the image as `sandbox-derby:latest`, run a throwaway container, verify `claude --version` works.

**Critical files created:** `sandbox/Dockerfile`, `sandbox/entrypoint-common.sh`, `go.mod`, `.env.example`, `loadouts/bare/.gitkeep`, `loadouts/example/`, `courses/example.md`

**Gotchas:**
- Claude Code install script may change — pin to a known-good approach or accept latest
- `debian:bookworm-slim` may not have everything Claude Code needs at runtime — verify interactively after build

---

## Phase 2: Drive

Interactive sandbox mode. Build a container, keep it alive, SSH (docker exec) in and work with the Claude agent. At the end of this phase, an operator can drive a sandbox with a loadout of their choice.

### Step 2.1: Drive entrypoint

#### Task 2.1.1: Write entrypoint-drive.sh
Source `entrypoint-common.sh` (validates env vars, configures git identity, copies loadout from staging to `~/.claude/`). Print available tools and instructions. Keep container alive (`exec tail -f /dev/null` or `exec bash`).

### Step 2.2: Compose service for drive mode

#### Task 2.2.1: Write docker-compose.yml with drive service
Build context: `./sandbox` (keeps build context scoped to the sandbox directory). Service `sandbox` using the drive entrypoint. Volume mounts:
- Loadout directory → `/home/agent/loadout/` (read-only staging)
- Workspace directory → `/home/agent/workspace` (if local mount desired)
Env file: `.env`. Resource limits: 2 CPUs, 4GB memory.

#### Task 2.2.2: Document drive workflow
In a brief README or usage doc: how to build, configure a loadout, start the sandbox, exec in, and interact with the agent.

### Step 2.3: Verify drive mode

#### Task 2.3.1: Build and drive a sandbox
Start the sandbox, exec in, confirm Claude Code CLI works, confirm loadout files are visible at the expected paths, run a simple `claude` interaction.

**Critical files created:** `sandbox/entrypoint-drive.sh`, `docker-compose.yml`
**Also depends on:** `sandbox/entrypoint-common.sh` (created in Phase 1)

**Gotchas:**
- The Dockerfile should not write to `/home/agent/.claude/` — the entrypoint populates it from the loadout staging path at runtime
- TTY allocation: `docker exec -it` requires the container to have a TTY-compatible entrypoint

---

## Phase 3: Coast

Autonomous sandbox mode. Hand it a course and a workspace, it runs to completion. At the end of this phase, an operator can coast a sandbox and inspect the resulting workspace and git history.

### Step 3.1: Coast entrypoint

#### Task 3.1.1: Write entrypoint-coast.sh
Source `entrypoint-common.sh` (validates env vars, configures git identity, copies loadout from staging to `~/.claude/`). Validate `TARGET_REPO` is set (required for coast mode). Clone `TARGET_REPO` into `/home/agent/workspace`. The common script then copies the course file from staging (`/home/agent/course/`) into the workspace so the agent can modify it (e.g., check off TODOs). Read course content. Construct prompt and execute via `claude -p "<prompt>"`. Exit when done.

### Step 3.2: Compose profile for coast mode

#### Task 3.2.1: Add coast service to docker-compose.yml
Service `coast` using the coast entrypoint, activated via compose profile. Volume mounts: loadout → `/home/agent/loadout/` (read-only staging), course file → `/home/agent/course/` (read-only staging). Same env file and resource limits as drive.

### Step 3.3: Output collection

#### Task 3.3.1: Design workspace output mechanism
After coast completes, the workspace (with its git commit history) must be accessible. Options: named volume, bind mount to a local output directory, or `docker cp`. For PoC, use a named volume per sandbox and `docker cp` to extract results.

### Step 3.4: Verify coast mode

#### Task 3.4.1: Coast a sandbox on an example course
Run a sandbox in coast mode with the example loadout and example course. Verify the agent executed the course, committed work to the workspace, and the results can be extracted.

**Critical files created:** `sandbox/entrypoint-coast.sh`, updated `docker-compose.yml`

**Gotchas:**
- `claude -p` behavior: confirm it exits cleanly after completing the prompt (needed for coast to terminate)
- Git identity must be configured before the agent starts working, or commits will fail
- Course file must be readable by the `agent` user inside the container

---

## Phase 4: Derby

Structured comparison across multiple sandboxes. A Go CLI reads a derby config, launches N sandboxes concurrently with different loadout/course/replica combinations, waits for completion, collects results, and generates a markdown report.

### Step 4.1: Derby config schema

#### Task 4.1.1: Define derby YAML schema
```yaml
name: <string>
image: <string, default "sandbox-derby:latest">
concurrency: <int, optional — max parallel sandboxes, defaults to total sandbox count>
workspace:
  repo: <git URL>
entries:
  - name: <string>
    loadout: <path to loadout dir>
    course: <path to course file>
    replicas: <int, default 1>
    resources:
      cpus: <string, default "2">
      memory: <string, default "4g">
```
Note: model selection lives in the loadout (settings.json or agent definitions), not at the derby level. `workspace.repo` is shared across all entries for PoC; per-entry repos are a known limitation to revisit later.

#### Task 4.1.2: Write derby.yaml.example
Demonstrate a comparison: same course, two loadouts, 2 replicas each.

### Step 4.2: Go CLI scaffold

#### Task 4.2.1: Write cmd/derby/main.go
CLI entrypoint with subcommands. For PoC, one command: `derby run <config.yaml>`. Parse the config, hand off to the orchestration layer.

#### Task 4.2.2: Write internal/derby/config.go
YAML config parsing and validation. Struct definitions matching the schema.

### Step 4.3: Orchestration

#### Task 4.3.1: Write internal/derby/runner.go
Core orchestration loop:
1. Parse config
2. Build image if not already built
3. For each entry × replicas, launch a container via `docker run` with the appropriate loadout mount, course mount, env vars, and resource limits
4. Use goroutines for concurrent execution, bounded by a configurable concurrency limit
5. Wait for all containers to complete
6. Collect exit codes and workspace contents

#### Task 4.3.2: Write internal/derby/sandbox.go
Encapsulate single-sandbox lifecycle: container creation, start, wait, output extraction. Wraps Docker CLI commands (shell out to `docker` for PoC; Docker SDK is an MVP upgrade path).

### Step 4.4: Reporting

#### Task 4.4.1: Write internal/derby/report.go
Generate a markdown report from collected results:
- Derby metadata (name, date, config summary)
- Per-sandbox results: entry name, replica number, exit code, duration
- Per-sandbox workspace summary: git log (commits, messages), files created/modified
- Comparison section: group by entry, highlight differences across replicas and across entries

#### Task 4.4.2: Write report to output directory
Save report as `derby-results/<derby-name>-<timestamp>/report.md`. Include raw git logs and any other collected artifacts alongside the report.

### Step 4.5: Verify derby

#### Task 4.5.1: Run an example derby
Configure a derby with 2 entries (bare vs example loadout), 2 replicas each, same course. Run it. Verify all sandboxes execute, results are collected, and the markdown report is generated and readable.

**Critical files created:** `cmd/derby/main.go`, `internal/derby/config.go`, `internal/derby/runner.go`, `internal/derby/sandbox.go`, `internal/derby/report.go`, `derby/derby.yaml.example`

**Gotchas:**
- Concurrent container launches may hit Docker resource limits — add a concurrency cap
- Container naming must be unique per sandbox (include derby name + entry name + replica number)
- Output extraction timing: must wait for container to fully stop before copying workspace
- If TARGET_REPO requires auth, GITHUB_TOKEN must be passed to each container
- Report generation should be resilient to individual sandbox failures (some may crash; report the failure, don't abort the derby)

---

## Success Criteria

- [ ] An operator can **drive** a sandbox: configure a loadout, start a container, exec in, and interact with a Claude agent
- [ ] An operator can **coast** a sandbox: point it at a course and a repo, let it run, and inspect the resulting workspace and git history
- [ ] An operator can run a **derby**: define a YAML config with multiple entries and replicas, run it with a single command, and get a markdown report comparing results
- [ ] Loadouts are swappable without rebuilding the image
- [ ] The system works with zero loadout (bare sandbox) as a baseline
