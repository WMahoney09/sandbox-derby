# Sandbox Derby

## Vision

Sandbox Derby is a tool for running agents in clean, isolated, configurable containers — and learning from what they do.

A **sandbox** is a containerized workspace for a Claude agent. It starts clean. You decide what's on the workbench: which skills, which agent definitions, which team definitions — or none at all. The agent works inside the sandbox, and the sandbox walls are real. Permissions aren't a suggestion; they're structural. The agent can have broad autonomy because the isolation is physical, not behavioral.

A **derby** is what happens when you line up multiple sandboxes, give them the same course with different loadouts, and let them race. Same hill, different cars. At the finish line, you compare: which sandbox produced the best result? Which skills helped? Which got in the way? The derby synthesizes those learnings and publishes them — initially as markdown reports, eventually as issues filed against the repos that define the tools under test.

## Canon

These terms are the project's ontology. Use them consistently in code, docs, and conversation.

**Sandbox.** A Docker container with a Claude agent inside it. Configurable loadout. Mountable workspace. Two modes: **drive** (interactive — SSH in and steer the agent) or **coast** (autonomous — hand it a task and walk away). No orchestration opinions. No lifecycle management. It runs, it works, it stops.

**Loadout.** The set of augmentations loaded into a sandbox: skills, agent definitions, team definitions. A sandbox can run bare — no loadout at all — or fully loaded. The loadout is the variable under test.

**Course.** The task given to a sandbox. Either a prompt string or a markdown file. In drive mode, the course may be irrelevant — the human is steering.

**Derby.** A structured comparison. N sandboxes, varying loadouts and courses at the officiant's discretion. The derby runs them, captures their outputs, evaluates the results, and distills what worked into actionable findings.

**Officiant.** The person designing and running a derby. They configure the knobs: loadouts, courses, replicas, models, and resource limits.

**Drive.** Interactive mode. SSH into a sandbox and work with the agent directly.

**Coast.** Autonomous mode. Hand the sandbox a course and let it run to completion.

## Derby Patterns

1. **Constant loadout, varying course.** Does this loadout generalize across different tasks?
2. **Varying loadout, constant course.** Which loadout outperforms on this task?
3. **Constant loadout, constant course, many replicas.** What's the baseline performance? What common pitfalls emerge from non-determinism?

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
- **Resources.** Container image and resource limits — test a loadout with constrained resources vs. maximal resources.
- **Loadout.** The set of skills, agent definitions, and team definitions loaded into each sandbox.
- **Course.** The task assigned to each sandbox.

## Why This Exists

Agent tooling — skills, agent definitions, team definitions — is only as good as the outcomes it produces. But today, improving that tooling is a manual, intuitive process. You write a skill, you use it, you notice it helped or didn't, you tweak it. That works at small scale but doesn't compound.

Sandbox Derby makes the feedback loop explicit and repeatable. Define a course. Run it with different loadouts. Measure the results. Publish the learnings. Repeat. Over time, the tools get better — not because someone guessed what to improve, but because the evidence showed what worked.

The sandbox is the workspace. The derby is the experiment. Together, they turn agent tooling from craft into engineering.

## Principles

**Isolation is structural.** A sandbox is a container. The agent's permissions are bounded by walls, not by instructions it may or may not follow. This means you can grant broad autonomy with confidence.

**The sandbox is agnostic.** It doesn't know what task it's been given, what tools it has, or whether it's part of a derby. It's a workspace. What happens inside is not its concern.

**Anthropic-first.** The primary focus is Claude agents for MVP. Anthropic leads in quality, and their philosophy aligns with the project's values. Other model providers may be supported in the future, but Claude is the first-class citizen.

**The derby is a consumer, not a feature of the sandbox.** The derby orchestrates sandboxes, but a sandbox doesn't need a derby to be useful. Running a single sandbox interactively is a first-class use case.

**Learnings leave the system as artifacts.** The derby doesn't apply its own findings. It publishes them — initially as markdown reports, eventually as issues or structured data — for external systems or humans to act on. This keeps the derby decoupled from whatever process consumes its output.

**Earn complexity.** Start with one sandbox and a task. Add crew when you want to test crew. Run a derby when you want to compare. Each layer is optional and additive.
