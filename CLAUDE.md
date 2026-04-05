# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoDaleks is a Go recreation of the classic Macintosh game "Daleks" (inspired by BSD UNIX "Robots"). It uses the **Ebiten v2** game engine and compiles to both native desktop and WebAssembly. The game runs at 800x600 on a 50x33 grid with 16px cells.

## Build System

This project uses [Taskfile](https://taskfile.dev/) (not Make). All commands are in `Taskfile.yml`.

```bash
task build          # Debug build → bin/godaleks.exe
task run            # Build and run
task release        # Optimized build with ldflags stripping
task test           # Run tests: go test -v ./test/...
task lint           # go fmt + go mod tidy
task vet            # go vet ./...
task staticcheck    # staticcheck ./...
task seccheck       # govulncheck ./...
task deps           # go mod tidy + download + update
task clean          # Remove bin/ and dist/
task goreleaser     # Build with goreleaser (snapshot)
```

Direct Go commands also work: `go run ./main.go`, `go build -o bin/godaleks.exe ./main.go`.

## Architecture

All game logic lives in the `cmd` package. `main.go` is just the Ebiten bootstrap (window setup + `RunGame`).

**Game loop**: `Game` struct in `cmd/types.go` holds all state. `root.go` implements the Ebiten `Game` interface (`Update`/`Draw`/`Layout`). The `Update` method is the main game loop with delta-time clamping (30-240 FPS range).

**State machine**: `GameState` enum in `constants.go` — `StateMenu → StatePlaying → StateGameOver/StateWin`.

**Key subsystems** (each in its own file in `cmd/`):
- `collision.go` — Grid-based collision detection with emperor invulnerability rules
- `movement.go` — Dalek AI movement toward player + smooth easing interpolation + Last Stand acceleration mode
- `effect.go` — Player actions: teleport (normal/safe), sonic screwdriver, last stand
- `draw.go` — Rendering: grid, sprites, HUD, collision effects
- `input.go` — Mouse click handling
- `sprites.go` — Programmatic sprite generation (scrap heap image)
- `loadimages.go` — Embedded PNG sprite loading (player, dalek, emperor)
- `loadsounds.go` — Embedded audio loading via Ebiten's audio system
- `menu.go` — Title screen and game-over screen rendering

**Entity model**: `Dalek` struct has both `GridPos` (logical) and `VisualPos` (interpolated for animation). The `IsEmperor` flag controls boss behavior (invulnerable until all normal daleks defeated).

**Performance**: The `Game` struct pre-allocates reusable maps (`positionMap`, `dalekRemoveMap`, `collisionPosMap`) and caches sprite dimensions to avoid per-frame allocations.

## WebAssembly

`godaleks.wasm` is a pre-built WASM binary served by `index.html` + `wasm_exec.js`. The WASM build uses the standard Go WASM target.

## CI/CD

GitHub Actions workflows in `.github/workflows/`: `build.yml` for CI builds, `goreleaser.yml` for releases. GoReleaser config targets Linux, Windows, and macOS.

## Workflow Orchestration

### 1. Plan Node Default

- Enter plan mode for ANY non-trivial task (3+ steps or architectural decisions)
- If something goes sideways, STOP and re-plan immediately – don't keep pushing
- Use plan mode for verification steps, not just building
- Write detailed specs upfront to reduce ambiguity

### 2. Subagent Strategy

- Use subagents liberally to keep main context window clean
- Offload research, exploration, and parallel analysis to subagents
- For complex problems, throw more compute at it via subagents
- One task per subagent for focused execution

### 3. Self-Improvement Loop

- After ANY correction from the user: update `tasks/lessons.md` with the pattern
- Write rules for yourself that prevent the same mistake
- Ruthlessly iterate on these lessons until mistake rate drops
- Review lessons at session start for relevant project

### 4. Verification Before Done

- Never mark a task complete without proving it works
- Diff behavior between main and your changes when relevant
- Ask yourself: "Would a staff engineer approve this?"
- Run tests, check logs, demonstrate correctness

### 5. Demand Elegance (Balanced)

- For non-trivial changes: pause and ask "is there a more elegant way?"
- If a fix feels hacky: "Knowing everything I know now, implement the elegant solution"
- Skip this for simple, obvious fixes – don't over-engineer
- Challenge your own work before presenting it

### 6. Autonomous Bug Fixing

- When given a bug report: just fix it. Don't ask for hand-holding
- Point at logs, errors, failing tests – then resolve them
- Zero context switching required from the user
- Go fix failing CI tests without being told how

## Task Management

1. **Plan First**: Write plan to `tasks/todo.md` with checkable items
2. **Verify Plan**: Check in before starting implementation
3. **Track Progress**: Mark items complete as you go
4. **Explain Changes**: High-level summary at each step
5. **Document Results**: Add review section to `tasks/todo.md`
6. **Capture Lessons**: Update `tasks/lessons.md` after corrections

## Core Principles

- **Simplicity First**: Make every change as simple as possible. Impact minimal code.
- **No Laziness**: Find root causes. No temporary fixes. Senior developer standards.
- **Minimal Impact**: Changes should only touch what's necessary. Avoid introducing bugs.

