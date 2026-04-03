# CLI Porcelain Layer

## Overview

Wrap all three sandbox-derby modes (drive, coast, derby run) behind a unified `derby` CLI binary with proper subcommand routing. Users will never run docker commands directly — `derby drive`, `derby coast`, and `derby run` are the only entry points.

## Notes

- **Out of scope:** External CLI libraries (Cobra, urfave/cli). We use Go's `flag` package with FlagSets.
- **Out of scope:** Modifying existing `internal/derby/` files except to export `checkImage` (rename to `CheckImage`).
- **Key decision:** `checkImage` in `runner.go` needs to be exported so `drive.go` and `coast.go` can call it. This is the only change to an existing file.
- **Key decision:** Drive mode starts the container detached, execs bash into it interactively, then cleans up on exit. This matches the entrypoint-drive.sh pattern (`tail -f /dev/null`).
- **Key decision:** Coast mode uses `--rm` since there's no artifact extraction needed for single runs.

## Progress

- [x] Phase 1: Export checkImage helper
- [x] Phase 2: Implement drive mode
- [ ] Phase 3: Implement coast mode
- [ ] Phase 4: Rewrite main.go with subcommand routing

---

## Phase 1: Export checkImage helper

**LOE: 1** (Complexity: Low | Impact: Low — single function rename in one file)

### Step 1.1: Rename checkImage to CheckImage in runner.go
#### Task 1.1.1: Rename the function declaration from `checkImage` to `CheckImage`
#### Task 1.1.2: Update the call site within `Runner.Run()` to use `CheckImage`

### Critical Files
- **Modified:** `internal/derby/runner.go`

---

## Phase 2: Implement drive mode

**LOE: 2** (Complexity: Medium | Impact: Low — new file, interactive docker lifecycle)

### Step 2.1: Create internal/derby/drive.go
#### Task 2.1.1: Define a `DriveConfig` struct with fields: Image, Loadout, EnvFile
#### Task 2.1.2: Implement `Drive(cfg DriveConfig) error` function that:
  - Calls `CheckImage` to verify the image exists
  - Resolves absolute paths for loadout and env file
  - Generates a unique container name (e.g. `derby-drive-<pid>` or `derby-drive-<timestamp>`)
  - Starts the container in detached mode with `docker run -d`
  - Uses `defer` to ensure container cleanup (stop + rm) on function exit
  - Execs `docker exec -it <container> bash` with stdin/stdout/stderr connected to os.Stdin/os.Stdout/os.Stderr
  - Returns the exec's error (or nil on clean exit)

### Gotchas
- The `-it` flags on `docker exec` require a real TTY. The Go process must connect os.Stdin/os.Stdout/os.Stderr to the child process for this to work.
- Container cleanup must happen even if the exec fails or is interrupted.
- Container name must be unique to avoid conflicts if the user runs multiple drive sessions.

### Critical Files
- **Created:** `internal/derby/drive.go`

---

## Phase 3: Implement coast mode

**LOE: 2** (Complexity: Medium | Impact: Low — new file, straightforward docker run)

### Step 3.1: Create internal/derby/coast.go
#### Task 3.1.1: Define a `CoastConfig` struct with fields: Image, Loadout, Course, Repo, EnvFile, SkipPermissions
#### Task 3.1.2: Implement `Coast(cfg CoastConfig) error` function that:
  - Calls `CheckImage` to verify the image exists
  - Resolves absolute paths for loadout, course, and env file
  - Builds docker run args: `--rm`, `--env-file`, `-e TARGET_REPO=...`, optional `-e SKIP_PERMISSIONS=true`, volume mounts for loadout and course, resource limits, image, and entrypoint
  - Connects os.Stdout and os.Stderr to the command
  - Also connects os.Stdin (so the user can Ctrl+C gracefully)
  - Runs the command and returns its error

### Critical Files
- **Created:** `internal/derby/coast.go`

---

## Phase 4: Rewrite main.go with subcommand routing

**LOE: 2** (Complexity: Medium | Impact: Low — single file rewrite, standard flag.FlagSet pattern)

### Step 4.1: Implement subcommand dispatch
#### Task 4.1.1: Parse os.Args[1] as the subcommand (drive, coast, run)
#### Task 4.1.2: Print usage and exit 1 if no subcommand or unknown subcommand

### Step 4.2: Implement drive subcommand handler
#### Task 4.2.1: Create a `flag.FlagSet` for "drive" with flags: `--loadout`, `--image`, `--env-file`
#### Task 4.2.2: Parse flags from os.Args[2:]
#### Task 4.2.3: Build a `DriveConfig` and call `derby.Drive()`

### Step 4.3: Implement coast subcommand handler
#### Task 4.3.1: Create a `flag.FlagSet` for "coast" with flags: `--loadout`, `--course`, `--repo`, `--skip-permissions`, `--image`, `--env-file`
#### Task 4.3.2: Parse flags from os.Args[2:]
#### Task 4.3.3: Validate that `--course` and `--repo` are provided (required)
#### Task 4.3.4: Build a `CoastConfig` and call `derby.Coast()`

### Step 4.4: Keep existing run subcommand
#### Task 4.4.1: Keep the existing `derby run <config.yaml>` logic (LoadConfig, NewRunner, Run, GenerateReport, WriteReport)

### Critical Files
- **Modified:** `cmd/derby/main.go`

---

## Gotchas & Risks

1. **TTY passthrough for drive mode:** The `docker exec -it` command needs a real terminal. If `derby drive` is run from a non-TTY context (e.g., piped), docker will error. This is acceptable — drive mode is inherently interactive.
2. **Container cleanup race:** If the process is killed with SIGKILL, the deferred cleanup won't run and a detached container will be left behind. Acceptable for PoC — the container name is predictable enough to clean up manually.
3. **Existing runner.go change:** Renaming `checkImage` to `CheckImage` is the only change to existing files. All call sites within runner.go must be updated.

## Success Criteria

1. `derby drive` starts a detached container and drops the user into an interactive bash shell. On exit, the container is stopped and removed.
2. `derby coast --course <file> --repo <url>` runs the coast entrypoint with output streaming to the terminal.
3. `derby run <config.yaml>` continues to work exactly as before.
4. `derby` (no args) and `derby <unknown>` print helpful usage messages.
5. `go build ./cmd/derby` succeeds with no errors.
6. No external dependencies added.
