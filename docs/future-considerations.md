# Future Considerations

Potential future workstreams. These are not yet scoped — they represent ideas, deferred decisions, and natural extensions that emerged during existing workstream work.

---

## LLM-as-Judge Evaluation

Automated evaluation of sandbox outputs using an LLM to compare results across sandboxes: which produced results sooner, which produced better results, and why. Replaces or augments human judgment in the derby report.

_Source: understanding phase — on the roadmap, not PoC/MVP_

## Agent Self-Retros

Agents within a sandbox running their own retrospectives on what worked well and what didn't — identifying what tripped them up, what could be better. Would feed structured self-assessment into the derby report alongside external evaluation.

_Source: understanding phase — on the roadmap, not PoC/MVP_

## Auto-Filing Issues Against Tool Repos

Derby findings published directly as issues against the repos that define the tools under test, rather than as a standalone markdown report. Requires attribution logic to determine which component of a loadout to file against.

_Source: understanding phase — future state_

## Non-Anthropic Model Support

Supporting agents powered by models other than Claude. The Docker isolation would work for any agent runtime. Claude is the first-class citizen; other providers would be additive.

_Source: understanding phase — future, Anthropic-first is a principle_

## CLI Porcelain

A Go CLI (`sbd`) with subcommands that wrap the Docker/Compose plumbing so the operator never runs docker commands directly. `sbd drive` stands up a sandbox and connects to it. `sbd coast` runs a sandbox autonomously against a course. `sbd run` executes a derby from a config file. The CLI is the first layer where the operator's interface is Sandbox Derby itself rather than Docker — fast follow after PoC, prerequisite to the TUI.

_Source: integration testing session — immediate post-PoC, before TUI_

## Charm TUI

A polished terminal UI built with Charm (Go) replacing the CLI as the primary interaction surface. Launch `sbd` with no subcommand to enter an interactive interface for driving, coasting, and running derbies. Builds on top of the CLI porcelain layer — same Go code underneath, richer interaction on top.

_Source: understanding phase — MVP scope, after CLI porcelain_

## Portable Loadouts

Standardizing the loadout format so a loadout tested in Sandbox Derby can be deployed as-is into other Claude agent systems (e.g., kubesat). The loadout already mirrors `~/.claude/` structure, but formal standardization and cross-project validation would make portability explicit.

_Source: planning phase — post-MVP goal_

## Build-Time Loadouts

Baking loadouts into the Docker image at build time (vs. volume-mounting at runtime). Provides reproducibility — the same image always has the same loadout. Useful for derby runs where exact reproducibility matters. May coexist with mount-time loadouts rather than replacing them.

_Source: solutioning phase — deferred, mount-time is sufficient for PoC_

## Workspace Publishing via PRs or Forks

Publishing sandbox workspace results as pull requests, worktrees, or repo forks rather than keeping them local. Would allow external tools (the "Derby Committee") to analyze results without access to the local machine.

_Source: understanding phase — MVP exploration, not PoC_

## Derby Committee

A separate tool or agent system that consumes derby output and performs deeper analysis — trend detection across multiple derby runs, longitudinal loadout performance tracking, automated recommendations for loadout changes.

_Source: understanding phase — implied by "another tool can do the analysis"_

## Sandbox IDs

Each sandbox in a derby receives a unique numeric identifier (e.g., Sandbox 42) — analogous to the number painted on a soapbox racer. The ID is assigned by the runner at derby configuration time, passed into the container as an environment variable, and used by the agent to identify itself in branch names, PR titles, commits, and issue references. The derby report correlates results by ID. Structural, not cosmetic — branch naming (`sandbox-42`), PR authorship, and result correlation all depend on it. Useful in both scrimmage and formal derby modes.

_Source: integration testing session — immediate, next implementation priority_

## Derby Scrimmage

A formal derby mode where the course is a GitHub repository with an issue backlog rather than a markdown file. The term "scrimmage" refers to the current informal mode (markdown course, quick and lightweight); a full "derby" is the structured version with more ceremony.

In a formal derby, each sandbox:
1. Gets a sandbox ID and creates a feature branch off main (`sandbox-42`)
2. Reads the issue backlog and works through it autonomously
3. For each issue (or group), branches off its feature branch, does the work, and opens a PR back into its feature branch — referencing issues without auto-closing them
4. When the backlog is exhausted, opens a final PR from its feature branch into main (without merging)
5. Exits

What the officiant sees on GitHub: main branch untouched, issue backlog still open, N pull requests into main (one per sandbox) each representing a complete attempt at the entire backlog. The repo is fully reusable for subsequent derbies.

_Source: integration testing session — post-sandbox-IDs, significant workstream_

## Docker SDK Migration

Replacing shell-out-to-`docker` in the Go derby orchestrator with the Docker SDK (`github.com/docker/docker/client`). Provides cleaner programmatic control, better error handling, and eliminates CLI parsing. Natural upgrade path from PoC.

_Source: solutioning/planning phase — MVP upgrade path_
