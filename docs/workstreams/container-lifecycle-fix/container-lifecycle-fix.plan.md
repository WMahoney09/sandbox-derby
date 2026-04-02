# Container Lifecycle Fix

## Overview

Fix the sandbox container lifecycle so that workspace artifacts (git log, file list) can be extracted after the agent finishes its run. Currently `docker run --rm` destroys the container immediately on exit, making artifact extraction impossible. Replace this with a keep-then-extract-then-cleanup approach using `docker cp`, local git commands, and `filepath.Walk`.

Additionally, add a workspace summary section to the coast entrypoint script so that structured output appears in stdout as a fallback.

## Notes

- The `extractSection` stub function is deleted entirely -- artifact extraction is now done via filesystem operations on the copied workspace, not stdout parsing.
- `docker exec` does not work on stopped containers, so the approach is: `docker cp` the workspace out, inspect it locally, then `docker rm`.
- The `newCommand` helper from `exec.go` is used for shelling out to docker/git.
- `GitLog` and `FileList` fields exist on `SandboxResult` but are not yet rendered in the report -- that is out of scope for this change.

## Progress

- [x] Phase 1: Fix container lifecycle and artifact extraction in sandbox.go
- [x] Phase 2: Add workspace summary to coast entrypoint

---

## Phase 1: Fix container lifecycle and artifact extraction in sandbox.go

### Step 1.1: Remove --rm and add extractArtifacts function

#### Task 1.1.1: Remove `--rm` from docker run args
Remove the `"--rm"` line from the args slice in `RunSandbox`.

#### Task 1.1.2: Write the `extractArtifacts` function
Create a new function `extractArtifacts(containerName string) (gitLog string, fileList string)` that:
1. Creates a temp dir via `os.MkdirTemp`
2. Runs `docker cp <container>:/home/agent/workspace <tempdir>/workspace`
3. Runs `git -C <tempdir>/workspace log --oneline` and captures stdout into `gitLog`
4. Walks `<tempdir>/workspace` with `filepath.Walk`, excluding any path containing `.git/`, collecting relative file paths into `fileList` (newline-separated)
5. Runs `docker rm <container>` to clean up the container
6. Runs `os.RemoveAll` on the temp dir
7. If any extraction step fails, logs a warning to stderr, still cleans up, and returns what it got

#### Task 1.1.3: Replace extractSection calls with extractArtifacts call
Replace the two `extractSection` calls with a single `extractArtifacts(containerName)` call, assigning the results to `result.GitLog` and `result.FileList`.

#### Task 1.1.4: Delete the extractSection function
Remove the `extractSection` function and its doc comment.

### Critical Files
- `internal/derby/sandbox.go` (modified)

### Gotchas & Risks
- `docker cp` from a stopped container: this works -- Docker allows `cp` from stopped containers. No need to restart.
- Temp dir cleanup: must happen in a `defer` or at end-of-function, even if extraction fails.
- `git log --oneline` will fail if the workspace has no commits -- handle gracefully by returning empty string.
- `filepath.Walk` on the workspace should skip the `.git` directory entirely (return `filepath.SkipDir` when encountering it) for efficiency.

---

## Phase 2: Add workspace summary to coast entrypoint

### Step 2.1: Add summary section to entrypoint-coast.sh

#### Task 2.1.1: Add workspace summary block before "Coast complete" banner
After the `claude -p` command and before the "Coast complete" echo block, add:
- A separator banner: `============================================` / `  Workspace Summary` / `============================================`
- Run `git -C /home/agent/workspace log --oneline` and print output
- This should not cause the script to fail if git log returns nothing -- use `|| true` or similar

### Critical Files
- `sandbox/entrypoint-coast.sh` (modified)

### Gotchas & Risks
- The script uses `set -euo pipefail` so any command that fails will abort. The `git log` must be guarded so that an empty repo (no commits) does not kill the script.

---

## Success Criteria

1. `RunSandbox` no longer passes `--rm` to docker run
2. After the container exits, `extractArtifacts` copies the workspace, extracts git log and file list, cleans up the container and temp dir
3. `extractSection` is deleted
4. The coast entrypoint prints a workspace summary section with git log output before the "Coast complete" banner
5. The project compiles cleanly with `go build ./...`
