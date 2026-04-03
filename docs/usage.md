# Usage

## Prerequisites

- Docker and Docker Compose
- Go 1.23+
- An Anthropic API key
- (Optional) A GitHub token for repo access

## Setup

Run the setup script from a fresh clone:

```
./setup.sh
```

This installs the `derby` CLI, creates `.env` from the example, and builds the Docker image. Edit `.env` with your API keys before running any commands.

## Loadouts

A loadout is the set of tools, skills, and configuration loaded into a sandbox. It mirrors the `~/.claude/` directory structure. Loadouts can come from three sources:

**Bare (no loadout).** The default. The sandbox runs with no skills or custom configuration — a clean baseline.

```
derby drive
derby coast --course ./courses/example.md --repo https://github.com/org/repo.git
```

**Local path.** A directory on the host machine. Useful for testing loadouts you're actively developing, or for combining multiple skill libraries into a single directory.

```
derby drive --loadout ./loadouts/example
derby drive --loadout /path/to/my/skills
```

**Remote (git URL).** A git repository cloned into the sandbox at startup. Useful for testing a skill library without downloading it locally, or for comparing two remote libraries against each other.

```
derby drive --loadout https://github.com/org/skills-repo.git
```

In a derby, each entry can use a different loadout source — mix local paths and git URLs freely:

```yaml
entries:
  - name: current-skills
    loadout: /path/to/my/skills
    course: ./courses/backlog.md

  - name: candidate-skills
    loadout: https://github.com/org/new-skills.git
    course: ./courses/backlog.md
```

This lets you test one skill library against another, or test a local modification against the published version, without changing anything on the host.

**Combining loadouts.** To test a union of multiple skill libraries, merge them into a single local directory and point the loadout there. The loadout directory mirrors `~/.claude/`, so skills go in a `skills/` subdirectory, agent guidance in `CLAUDE.md`, and settings in `settings.json`.

## Drive Mode (Interactive)

Start an interactive sandbox session:

```
derby drive
derby drive --loadout ./loadouts/example
derby drive --loadout https://github.com/org/skills.git
```

This drops you into a bash shell inside the container. Run `claude` to start a Claude session. When you exit, the container is cleaned up.

## Coast Mode (Autonomous)

Run a sandbox autonomously against a course:

```
derby coast --course ./courses/example.md --repo https://github.com/org/repo.git --skip-permissions
derby coast --course ./courses/example.md --repo https://github.com/org/repo.git --loadout https://github.com/org/skills.git --skip-permissions
```

The `--repo` is the workspace the agent works in (cloned at startup). The `--course` is the task (a markdown file). Use `--skip-permissions` for autonomous execution without permission prompts.

## Derby (Comparative)

Run multiple sandboxes with varying configurations and get a report:

```
derby run examples/derby.yaml.example
```

Each sandbox gets a unique ID (Sandbox #1, #2, etc.) for correlation in the report. See `examples/` for configuration format. Key fields per entry:

- `loadout` — local path or git URL
- `course` — path to a markdown course file
- `skip_permissions` — boolean, per entry
- `replicas` — how many sandboxes to run with this configuration

The report is written to `derby-results/` with results keyed by sandbox ID.
