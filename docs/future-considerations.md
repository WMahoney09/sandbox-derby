# Future Considerations

Potential future workstreams. These are not yet scoped — they represent ideas, deferred decisions, and natural extensions that emerged during existing workstream work.

---

## ~~CLI Porcelain~~ (complete — v0.1)

Shipped as `derby drive`, `derby coast`, `derby run`. Binary name is `derby`.

## ~~Sandbox IDs~~ (complete — v0.1)

Sequential IDs assigned per sandbox, passed as `SANDBOX_ID` env var, shown in banners, prompts, container names, and reports.

---

## Scrimmage Rename

Rename `derby run` to `derby scrimmage` to distinguish the current informal mode (markdown course, one-shot execution, local report) from the formal `derby race` lifecycle. This is a prerequisite naming change before building the race commands.

_Source: integration testing session — immediate, next implementation step_

## Derby Race Lifecycle

A formal derby mode with a multi-step lifecycle, replacing the one-shot `derby run` for structured evaluations. The course is a GitHub repository with an issue backlog. Each sandbox gets a feature branch, works through the backlog via PRs into its own branch, and opens a final PR into main when done.

CLI commands:

```
derby race setup config.yaml    # validate repo, assign sandbox IDs, show schedule
derby race start                # launch sandboxes; reports results when all finish
derby race status               # check progress while running
derby race conclude             # force-stop all sandboxes, report results (DNF for incomplete)
derby race results              # view results — no args = most recent, or pass a name/path
```

`setup` is where the officiant reviews the lineup and confirms. `start` runs the race and produces a report when all sandboxes finish naturally. `conclude` ends the race early — incomplete sandboxes are marked DNF (did not finish). `results` is read-only, for viewing reports from past races.

Each sandbox:
1. Gets a sandbox ID and creates a feature branch off main (`sandbox-42`)
2. Reads the issue backlog and works through it autonomously
3. For each issue (or group), branches off its feature branch, does the work, and opens a PR back into its feature branch — referencing issues without auto-closing them
4. When the backlog is exhausted, opens a final PR from its feature branch into main (without merging)
5. Exits

What the officiant sees on GitHub: main branch untouched, issue backlog still open, N pull requests into main (one per sandbox) each representing a complete attempt at the entire backlog. The repo is fully reusable for subsequent races.

_Source: integration testing session — significant workstream, after scrimmage rename_

## Derby Events

Community-scale races where the course repo is also the venue. The repo contains everything needed to run a race: the issue backlog (the course), a `contestants/` directory (the loadouts), and event configuration (replicas, resources, etc.).

Enrollment is PR-based: anyone who wants to compete PRs their loadout into `contestants/`. The officiant reviews and merges before the race. The lineup is transparent — everyone can see who entered what. The repo is versioned, so the history is the record of every event.

This can be automated. A kubesat instance on a monthly orbital period could officiate: pull the course repo, see who's in `contestants/`, run `derby race setup` + `derby race start`, and publish results. No human at the keyboard.

Design implications for `derby race`: the setup/start/conclude lifecycle must be fully automatable from the start so an external system (kubesat, CI, cron) can drive it without interactive prompts.

_Source: integration testing session — future, builds on derby race lifecycle_

## Charm TUI

A polished terminal UI built with Charm (Go) replacing the CLI as the primary interaction surface. Launch `derby` with no subcommand to enter an interactive interface for driving, coasting, and running scrimmages/races. Builds on top of the CLI porcelain layer — same Go code underneath, richer interaction on top.

_Source: understanding phase — MVP scope, prerequisite (CLI porcelain) is complete_

## LLM-as-Judge Evaluation

Automated evaluation of sandbox outputs using an LLM to compare results across sandboxes: which produced results sooner, which produced better results, and why. Replaces or augments human judgment in the derby report.

_Source: understanding phase — on the roadmap, not PoC/MVP_

## Agent Self-Retros

Agents within a sandbox running their own retrospectives on what worked well and what didn't — identifying what tripped them up, what could be better. Would feed structured self-assessment into the derby report alongside external evaluation.

_Source: understanding phase — on the roadmap, not PoC/MVP_

## Auto-Filing Issues Against Tool Repos

Derby findings published directly as issues against the repos that define the tools under test, rather than as a standalone markdown report. Requires attribution logic to determine which component of a loadout to file against.

_Source: understanding phase — future state_

## Workspace Publishing via PRs or Forks

Publishing sandbox workspace results as pull requests, worktrees, or repo forks rather than keeping them local. Would allow external tools (the "Derby Committee") to analyze results without access to the local machine. Partially addressed by the derby race lifecycle (which produces PRs), but scrimmage mode still keeps results local.

_Source: understanding phase — MVP exploration, not PoC_

## Derby Committee

A separate tool or agent system that consumes derby output and performs deeper analysis — trend detection across multiple derby runs, longitudinal loadout performance tracking, automated recommendations for loadout changes.

_Source: understanding phase — implied by "another tool can do the analysis"_

## Portable Loadouts

Standardizing the loadout format so a loadout tested in Sandbox Derby can be deployed as-is into other Claude agent systems (e.g., kubesat). The loadout already mirrors `~/.claude/` structure, but formal standardization and cross-project validation would make portability explicit.

_Source: planning phase — post-MVP goal_

## Build-Time Loadouts

Baking loadouts into the Docker image at build time (vs. volume-mounting at runtime). Provides reproducibility — the same image always has the same loadout. Useful for derby runs where exact reproducibility matters. May coexist with mount-time loadouts rather than replacing them.

_Source: solutioning phase — deferred, mount-time is sufficient for PoC_

## Non-Anthropic Model Support

Supporting agents powered by models other than Claude. The Docker isolation would work for any agent runtime. Claude is the first-class citizen; other providers would be additive.

_Source: understanding phase — future, Anthropic-first is a principle_

## Docker SDK Migration

Replacing shell-out-to-`docker` in the Go derby orchestrator with the Docker SDK (`github.com/docker/docker/client`). Provides cleaner programmatic control, better error handling, and eliminates CLI parsing. Natural upgrade path from PoC.

_Source: solutioning/planning phase — MVP upgrade path_
