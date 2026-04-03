# Sandbox Derby
## 🚧 WIP 🏗️ 

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

A **loadout** is the variable under test. The set of augmentations loaded into a sandbox: skills, agent definitions, team definitions. A sandbox can run bare (no loadout at all) or fully loaded.

A **course** is the task given to a sandbox. A markdown file with instructions. In drive mode the course may be irrelevant — the human is steering.

The goal: turn agent tooling improvement from guesswork into evidence.

## Quick Start

```
git clone https://github.com/WMahoney09/sandbox-derby.git
cd sandbox-derby
./setup.sh        # installs CLI, builds image, creates .env
# edit .env with your API keys
derby drive        # interactive sandbox
```

## Modes

- **Drive** — Interactive. Enter a sandbox and work with the agent directly.
- **Coast** — Autonomous. Hand the sandbox a task and let it run to completion.
- **Scrimmage** — Comparative. Launch multiple sandboxes with varying configurations and get a report comparing outcomes.

```
derby drive --loadout https://github.com/org/skills.git
derby coast --course ./courses/task.md --repo https://github.com/org/repo.git --skip-permissions
derby scrimmage examples/derby.yaml.example
```

## Loadouts

A loadout mirrors the `~/.claude/` directory structure and can come from three sources:

**Bare** — no loadout. The sandbox runs with vanilla Claude. This is the default and serves as a baseline for comparison.

**Local path** — a directory on the host. Point it at a skill library you're developing, or merge multiple libraries into one directory to test a combined loadout.

```
derby drive --loadout ./loadouts/example
derby drive --loadout /path/to/my/skills
```

**Remote (git URL)** — a repository cloned into the sandbox at startup. Test a skill library without downloading it locally, or compare two remote libraries head-to-head.

```
derby drive --loadout https://github.com/org/skills-a.git
derby coast --loadout https://github.com/org/skills-b.git --course ./courses/task.md --repo https://github.com/org/repo.git
```

In a derby, mix freely — test a local work-in-progress loadout against the published remote version in the same run.

## Why?

Agent tooling — skills, agent definitions, team definitions — is only as good as the outcomes it produces. Today, improving that tooling is manual and intuitive. You write a skill, use it, notice whether it helped, and tweak it. That works at small scale but doesn't compound.

Sandbox Derby makes the feedback loop explicit and repeatable. Define a task. Run it with different tool configurations. Compare results. Publish learnings. Repeat.

## Status

The CLI (`derby drive`, `derby coast`, `derby scrimmage`) is functional. Drive, coast, and scrimmage modes work with local and remote loadouts. See [vision-statement.md](vision-statement.md) for the full design and [docs/usage.md](docs/usage.md) for detailed usage.

## Tech

- **Go** CLI for orchestration
- **Docker** for sandbox isolation
- **Claude agents** as the first-class citizen (Anthropic-first)
- Base image: `debian:bookworm-slim` + Claude Code CLI, gh, Node.js, Python, git

## License

TBD
