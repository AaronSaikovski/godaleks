# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoDaleks is a Go recreation of the classic Macintosh game "Daleks" (inspired by BSD UNIX "Robots") using the **Ebiten v2.9** game engine. Compiles to native desktop (Windows/Linux/macOS) and WebAssembly. Runs at 800x600 on a 50x33 grid with 16px cells.

## Build System

Uses [Taskfile](https://taskfile.dev/) — all commands in `Taskfile.yml`:

```bash
task build          # Debug build → bin/godaleks.exe
task run            # go run ./main.go
task release        # Optimized: -ldflags="-s -w" → bin/godaleks.exe
task test           # go test -v ./test/...
task lint           # go fmt ./... && go mod tidy
task vet            # go vet ./...
task staticcheck    # staticcheck ./...
task seccheck       # govulncheck ./...
task deps           # go mod tidy + download + update
task clean          # Remove bin/ and dist/
task goreleaser     # goreleaser release --snapshot --clean
```

## Go-Specific Notes

- **Go version**: 1.26+ (see `go.mod`)
- **Module path**: `github.com/AaronSaikovski/godaleks`
- **CGO**: Required for Linux (GLFW/X11/ALSA). Not required for Windows/macOS (Ebitengine uses `purego`). Not applicable for WASM.
- **WASM build**: `GOOS=js GOARCH=wasm go build -ldflags="-s -w" -trimpath -o godaleks.wasm ./main.go`
- **Embedded assets**: All images (`cmd/assets/*.png`) and sounds (`cmd/assets/*.wav`) use `//go:embed` directives in `cmd/loadimages.go` and `cmd/loadsounds.go`. No external files needed at runtime.
- **No test files exist yet** — `task test` runs against `./test/...` but that directory doesn't exist.

## Releasing

Releases are triggered by pushing a version tag:

```bash
git tag v1.2.1
git push origin v1.2.1
```

This triggers `.github/workflows/goreleaser.yml` which builds:
- **Linux amd64** on ubuntu-latest (CGO_ENABLED=1, requires X11/ALSA dev libs)
- **Windows amd64** on ubuntu-latest (CGO_ENABLED=0)
- **macOS amd64/arm64** on macos-latest (CGO_ENABLED=1)

To re-tag after fixing a release issue:
```bash
git tag -d v1.2.1 && git push origin :refs/tags/v1.2.1
git tag v1.2.1 && git push origin v1.2.1
```

## CI Workflows

- **build.yml** — Runs on push to `main`. Installs Linux X11/ALSA dev deps, runs `go vet`, `go fmt` check, native build.
- **deploy-wasm.yml** — Runs on push to `main`. Builds WASM binary, deploys to GitHub Pages. Runs `go vet` under `GOOS=js GOARCH=wasm` to avoid native dep requirements.
- **goreleaser.yml** — Runs on version tag push (`v*`). Multi-job: Linux/Windows on Ubuntu, macOS on macOS runner, then merges artifacts into a GitHub release.

## Architecture

All game logic is in the `cmd` package. `main.go` is the Ebiten bootstrap (`NewGame` + `RunGame`).

**Game loop**: `Game` struct (`cmd/types.go`) holds all state. `root.go` implements Ebiten's `Game` interface (`Update`/`Draw`/`Layout`). Delta-time is clamped to 30-240 FPS range.

**State machine** (`constants.go`): `StateMenu → StatePlaying → StateLevelComplete (1.5s delay) → StatePlaying` or `→ StateGameOver/StateWin`.

**Key subsystems** (each in own file under `cmd/`):
- `collision.go` — Grid-based collision with emperor invulnerability rules. Reuses pre-allocated maps (`scrapMap`, `daleksByPosition`, `dalekRemoveMap`) to avoid per-frame allocations.
- `movement.go` — Dalek AI toward player + smootherstep easing + Last Stand acceleration. Reuses `survivingDaleksBuf` and `newScrapsBuf`.
- `effect.go` — Player actions: teleport, sonic screwdriver, last stand. Also visual particle effects using pre-computed trig tables.
- `draw.go` — Rendering with cached strings (HUD, game-over, Last Stand). Uses `gridOffsetX`/`gridOffsetY` compile-time constants.
- `game.go` — Level management, entity spawning. `occupancyGrid [50][33]bool` and `scrapGrid [50][33]bool` for O(1) position lookups.
- `loadsounds.go` — Audio via Ebiten's audio system. Creates fresh `audio.Player` per playback to allow overlapping sounds.
- `input.go` — Mouse click → grid coordinate conversion and move validation.
- `menu.go` — Title screen rendering.

**Entity model**: `Dalek` struct has `GridPos` (logical) and `VisualPos` (interpolated for animation). `IsEmperor` flag controls boss behavior (invulnerable until all normal daleks defeated).

**Performance patterns used throughout**:
- Pre-allocated reusable buffers on `Game` struct (maps, slices) to eliminate per-frame GC pressure
- Grid arrays instead of maps for O(1) spatial lookups
- In-place slice filtering instead of allocating new slices
- Cached `fmt.Sprintf` results, only reformatted when values change
- Pre-computed sin/cos lookup table for particle effects
