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

## Charm TUI

A polished terminal UI built with Charm (Go) replacing CLI commands and docker compose as the primary interaction surface. Covers driving, coasting, and running derbies.

_Source: understanding phase — MVP scope, not PoC_

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

## Docker SDK Migration

Replacing shell-out-to-`docker` in the Go derby orchestrator with the Docker SDK (`github.com/docker/docker/client`). Provides cleaner programmatic control, better error handling, and eliminates CLI parsing. Natural upgrade path from PoC.

_Source: solutioning/planning phase — MVP upgrade path_
