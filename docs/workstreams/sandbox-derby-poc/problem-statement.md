# Problem Statement

## Problem

Improving agent tooling — skills, agent definitions, team definitions — is a manual, intuitive process with no repeatable feedback loop. There's no systematic way to test whether a set of tools helps an agent perform better, compare toolkits against each other, or isolate the effect of non-determinism from the effect of tooling choices.

## Why It Matters / Why Now

Agent tooling that isn't empirically tested doesn't compound. Teams improve tools by feel, not evidence. As the number of skills and agent configurations grows, the combinatorial space makes intuition insufficient. A structured testing and comparison system is needed to turn agent tooling from craft into engineering.

## Key Constraints

- **Greenfield.** No existing codebase — starting from scratch.
- **Anthropic-first.** Claude agents are the first-class citizen. Other model providers are future state.
- **Docker-based isolation.** Sandboxes are containers. Isolation is structural, not behavioral.
- **Go + Charm for TUI.** MVP interaction surface will be a Charm-based TUI in Go. PoC uses CLI commands and `docker compose` directly.
- **Base image reference.** Adapted from kubesat project (`debian:bookworm-slim` + Claude Code CLI, gh, Node.js, Python, git, non-root `agent` user).
- **API key auth.** Passed as environment variables.
- **Phased delivery.** PoC (rough but functional) → MVP (smooth, TUI, reliable).

## Success Criteria

**PoC is done when:**
- A user can **drive** a sandbox — build and enter an interactive container with a Claude agent and a configurable loadout, and work with the agent via SSH.
- A user can **coast** a sandbox — hand it a task (string or markdown file), let it run autonomously, and collect the output.
- A user can run a **derby** — launch multiple sandboxes concurrently with varying loadouts, courses, and replica counts, and receive a markdown report comparing results.

**MVP is done when:**
- All three modes work reliably through a Charm TUI.
- The experience of driving, coasting, and running derbies is smooth and polished.

## Assumptions Surfaced

- **Loadout mechanism is TBD.** Build-time (baked into image) vs. mount-time (volume mounts) each have tradeoffs. May need both — build-time for reproducible derby runs, mount-time for fast iteration during development. Needs further exploration.
- **Workspace is open-ended.** Could be a cloned repo, an empty directory, or anything the operator specifies. Start with the simplest thing that works.
- **Evaluation is human-only for PoC/MVP.** LLM-as-judge and agent self-retros are roadmap items, not PoC/MVP scope.
- **Derby patterns are known.** Three patterns identified: vary course (generalization testing), vary loadout (comparison testing), hold constant (non-determinism/baseline testing).

## Workstream Slug

`sandbox-derby-poc`
