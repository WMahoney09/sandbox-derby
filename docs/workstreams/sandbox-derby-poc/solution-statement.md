# Solution Statement

## Candidates

### Candidate A: Compose + Shell Scripts

Everything runs through Docker Compose. Shell scripts orchestrate derbies. Zero Go code in PoC. Volume-mounted loadouts. Drive via `docker compose exec`, coast via `docker compose run`, derby via a shell script that launches N Compose runs and collects output.

**Strengths:** Fastest path to working PoC. Familiar tooling. Zero rebuilds during loadout iteration.
**Tradeoffs:** Shell scripts for concurrent derby orchestration get messy fast. Two paradigm shifts to MVP (shell → Go, possibly Compose → Docker SDK). Compose scaling doesn't natively support different loadouts per sandbox.
**LOE:** 2 (PoC only — rewrite cost for MVP is additional)

### Candidate B: Go CLI from Day One

Go CLI using Docker SDK handles everything — image builds, container lifecycle, TTY for drive, concurrent containers for derby. No Compose at all. Single paradigm from PoC through MVP.

**Strengths:** No throwaway work. Full programmatic control. Natural path to Charm TUI. Go concurrency model fits derby orchestration perfectly.
**Tradeoffs:** Significantly more upfront work. TTY passthrough via Docker SDK is non-trivial. Building a container orchestrator is a known source of accidental complexity. Harder to debug — custom code sits between operator and Docker.
**LOE:** 4

### Candidate C: Compose for Sandbox, Go for Derby (selected)

Compose owns sandbox lifecycle — drive and coast work immediately with zero custom code. Go enters only for derby orchestration: reading config, launching N Compose runs with varying loadouts/courses/replicas, collecting outputs, generating markdown reports. Go grows into the TUI for MVP, eventually absorbing drive/coast.

**Strengths:** Drive mode works almost immediately. Go earns its place — only appears where Compose can't handle the job (concurrent orchestration, reporting). Natural evolution path: Go starts small, absorbs full surface for MVP. Supports both mount-time and build-time loadouts without forcing an early decision. Compose remains as a debugging escape hatch.
**Tradeoffs:** Two tools in the stack during PoC (Compose + Go). Boundary between Compose territory and Go territory shifts as project matures. Compose invocation from Go is less clean than pure Docker SDK.
**LOE:** 3

## Next Step
Recommendation: Plan
Confidence: high
Rationale: Approach C was the clear choice — it earns complexity incrementally and gives each tool (Compose, Go) a well-defined role that evolves naturally toward MVP.
