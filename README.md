# Sandbox Derby 🚧 WIP 🏗️ 

```

                             ██
          ░░▒▒▓▓█████████████  █████████████████████
  ░  ░  ░▒▒▓▓█████████████████████████████████████████
   ░  ░ ░▒▓▓████████████████████████████████████████████
  ░  ░  ▒▓█ _____ ██████████████████████████ _____ █████
   ░  ░    /     \                          /     \
  ░  ░    |   ●   |                        |   ●   |
           \_____/                          \_____/

  ___________________   ____________________________  __
  __  ___/__    |__  | / /__  __ \__  __ )_  __ \_  |/ /
  _____ \__  /| |_   |/ /__  / / /_  __  |  / / /_    /
  ____/ /_  ___ |  /|  / _  /_/ /_  /_/ // /_/ /_    |
  /____/ /_/  |_/_/ |_/  /_____/ /_____/ \____/ /_/|_|

         ___________________________________  __
         ___  __ \__  ____/__  __ \__  __ ) \/ /
         __  / / /_  __/  __  /_/ /_  __  |_  /
         _  /_/ /_  /___  _  _, _/_  /_/ /_  /
         /_____/ /_____/  /_/ |_| /_____/ /_/

              Test your Skills... & your Agents

```

Run agents in clean, isolated, configurable containers — and learn from what they do.

## What Is This?

Sandbox Derby is a tool for testing and comparing agent tooling empirically. It gives you two things:

1. **Sandboxes** — Docker containers with a Claude agent inside. You configure what tools the agent has (its _loadout_), hand it a task (its _course_), and let it work. Isolation is structural: the container walls are real, so agents can operate with broad autonomy safely.

2. **Derbies** — structured comparisons across sandboxes. Run the same task with different toolkits, or the same toolkit across different tasks, and compare results. The derby captures outputs and distills findings into actionable reports.

A **loadout** — The loadout is the variable under test. The set of augmentations loaded into a sandbox: skills, agent definitions, team definitions. A sandbox can run bare (loadout at all) or fully loaded.

A **course** is the task given to a sandbox. Either a prompt string or a markdown file. In drive mode the course may be irrelevant — the human is steering.

The goal: turn agent tooling improvement from guesswork into evidence.

## Why?

Agent tooling — skills, agent definitions, team definitions — is only as good as the outcomes it produces. Today, improving that tooling is manual and intuitive. You write a skill, use it, notice whether it helped, and tweak it. That works at small scale but doesn't compound.

Sandbox Derby makes the feedback loop explicit and repeatable. Define a task. Run it with different tool configurations. Compare results. Publish learnings. Repeat.

## Modes

- **Drive** — Interactive. Enter a sandbox and work with the agent directly.
- **Coast** — Autonomous. Hand the sandbox a task and let it run to completion.
- **Derby** — Comparative. Launch multiple sandboxes with varying configurations and get a report comparing outcomes.

## Status

Early development. The project is in the PoC phase — rough but functional is the target. See [vision-statement.md](vision-statement.md) for the full design and [docs/workstreams/sandbox-derby-poc/problem-statement.md](docs/workstreams/sandbox-derby-poc/problem-statement.md) for scope and success criteria.

## Tech

- **Go** with Charm for the TUI (MVP phase)
- **Docker** for sandbox isolation
- **Claude agents** as the first-class citizen (Anthropic-first)
- Base image adapted from kubesat (`debian:bookworm-slim` + Claude Code CLI, gh, Node.js, Python, git)

## License

TBD
