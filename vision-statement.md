# Sandbox Derby

## Vision

Sandbox Derby is a tool for running agents in clean, isolated, configurable containers — and learning from what they do.

A **sandbox** is a containerized workspace for a Claude agent. It starts clean. You decide what's on the workbench: which skills, which agent definitions, which team definitions — or none at all. The agent works inside the sandbox, and the sandbox walls are real. Permissions aren't a suggestion; they're structural. The agent can have broad autonomy because the isolation is physical, not behavioral.

A **derby** is what happens when you line up multiple sandboxes, give them the same task with different loadouts, and let them race. Same hill, different cars. At the finish line, you compare: which sandbox produced the best result? Which skills helped? Which got in the way? The derby synthesizes those learnings and publishes them — typically as issues filed against the repos that define the tools under test.

## Core Concepts

**Sandbox.** A Docker container with a Claude agent inside it. Configurable crew. Mountable workspace. Two modes: interactive (SSH in and drive the agent like your own shell) or autonomous (hand it a task and walk away). No orchestration opinions. No lifecycle management. It runs, it works, it stops.

**Crew.** The set of augmentations loaded into a sandbox: skills, agent definitions, team definitions. A sandbox can run bare — no crew at all — or fully loaded. The crew is the variable under test.

**Derby.** A structured comparison. N sandboxes, same task, different crews. The derby runs them, captures their outputs, evaluates the results, and distills what worked into actionable findings. Those findings leave the system as issues, ready to be picked up by whatever process maintains the tools under test.

## Evaluation

MVP evaluation is human judgment. The derby produces a markdown report; a human reads it and decides what to act on.

On the roadmap: agents running their own retros inside each sandbox — identifying what tripped them up, what could be better. Beyond that, LLM-as-judge evaluation comparing outcomes across sandboxes: which produced results sooner, which produced better results, and why.

Automated evaluation is the long game, but it's earned incrementally. Human judgment first, agent self-assessment second, cross-sandbox comparative judgment third.

## Tasks

A task is what gets handed to a sandbox. Two shapes:

1. **A string.** A simple prompt, sufficient for straightforward or well-scoped work.
2. **A markdown file.** For more verbose or structured prompts that need context, constraints, or acceptance criteria.

In interactive mode, the task may be irrelevant — the human is steering.

## Derby Configuration

The derby is controlled by an officiant — the person designing and running the experiment. Initially this is a config file; the vision is for it to become a UI.

The officiant tunes knobs:

- **Replicas.** How many sandboxes to run with the same loadout, to account for non-determinism and build statistical confidence.
- **Model.** Which model powers the agent in each sandbox.
- **Resources.** Container image and resource limits — test a crew with constrained resources vs. maximal resources.
- **Crew.** The set of skills, agent definitions, and team definitions loaded into each sandbox.

## Why This Exists

Agent tooling — skills, agent definitions, team definitions — is only as good as the outcomes it produces. But today, improving that tooling is a manual, intuitive process. You write a skill, you use it, you notice it helped or didn't, you tweak it. That works at small scale but doesn't compound.

Sandbox Derby makes the feedback loop explicit and repeatable. Define a task. Run it with different toolkits. Measure the results. Publish the learnings. Repeat. Over time, the tools get better — not because someone guessed what to improve, but because the evidence showed what worked.

The sandbox is the workspace. The derby is the experiment. Together, they turn agent tooling from craft into engineering.

## Principles

**Isolation is structural.** A sandbox is a container. The agent's permissions are bounded by walls, not by instructions it may or may not follow. This means you can grant broad autonomy with confidence.

**The sandbox is agnostic.** It doesn't know what task it's been given, what tools it has, or whether it's part of a derby. It's a workspace. What happens inside is not its concern.

**Anthropic-first.** The primary focus is Claude agents for MVP. Anthropic leads in quality, and their philosophy aligns with the project's values. Other model providers may be supported in the future, but Claude is the first-class citizen.

**The derby is a consumer, not a feature of the sandbox.** The derby orchestrates sandboxes, but a sandbox doesn't need a derby to be useful. Running a single sandbox interactively is a first-class use case.

**Learnings leave the system as artifacts.** The derby doesn't apply its own findings. It publishes them — initially as markdown reports, eventually as issues or structured data — for external systems or humans to act on. This keeps the derby decoupled from whatever process consumes its output.

**Earn complexity.** Start with one sandbox and a task. Add crew when you want to test crew. Run a derby when you want to compare. Each layer is optional and additive.
