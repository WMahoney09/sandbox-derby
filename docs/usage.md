# Usage

## Prerequisites

- Docker and Docker Compose
- An Anthropic API key
- (Optional) A GitHub token for repo access

## Setup

Copy the example env file and add your keys:

```
cp .env.example .env
# Edit .env with your ANTHROPIC_API_KEY and GITHUB_TOKEN
```

## Drive Mode (Interactive)

Build the image and start a sandbox:

```
docker compose up -d
```

Attach to the sandbox:

```
docker compose exec sandbox bash
```

You're now inside the sandbox. Run `claude` to start an interactive Claude session.

To change the loadout, edit the volume mount in `docker-compose.yml`:

```yaml
volumes:
  - ./loadouts/example:/home/agent/loadout:ro   # swap "example" for your loadout
```

To run bare (no loadout):

```yaml
volumes:
  - ./loadouts/bare:/home/agent/loadout:ro
```

Stop the sandbox:

```
docker compose down
```

## Coast Mode (Autonomous)

Run a sandbox in coast mode with a course and target repo:

```
docker compose run --rm \
  -e TARGET_REPO=https://github.com/org/repo.git \
  -v ./courses/example.md:/home/agent/course/course.md:ro \
  coast
```

The agent will clone the repo, execute the course, and exit. Inspect results via `docker cp` or by mounting an output volume.

## Derby

See `examples/derby.yaml.example` for the configuration format. Run a derby with:

```
go run ./cmd/derby run examples/derby.yaml.example
```
