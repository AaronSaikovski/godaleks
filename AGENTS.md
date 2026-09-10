# Repository Guidelines

## Project Overview

GoDaleks recreates Macintosh Daleks/BSD Robots in Go with Ebitengine. It targets native desktop and browser WebAssembly, using an 800×600 viewport and a 50×33 grid of 16px cells. Gameplay combines robot-chase collisions, scrap hazards, teleportation, a sonic screwdriver, and Last Stand.

## Architecture & Data Flow

- `main.go` → `cmd.NewGame()` → `ebiten.RunGame()`. All gameplay lives in one `cmd` package; files separate subsystems, not service layers.
- `cmd/types.go`: `Game` owns entities, resources, timers, effects, spatial caches, render buffers, and concrete image/audio dependencies. `cmd/root.go` implements `Update`, `Draw`, and fixed-size `Layout`.
- `Update` processes mouse input, advances animations/effects, then dispatches keyboard input and game state. Simulation uses `1/ebiten.TPS()` with a 60-TPS fallback; some input/HUD timers use wall-clock time.
- Normal movement (`cmd/movement.go`) updates logical `GridPos` immediately, interpolates `VisualPos` with smootherstep, and resolves collisions after movement finishes. Rendering reads visual positions. Last Stand has a separate continuous accelerating movement/collision path; preserve that distinction.
- `cmd/collision.go` resolves player death before enemy/scrap collisions, scoring, and level progression. Normal collision handling protects emperors until ordinary Daleks are gone; check the separate Last Stand path when changing combat rules.
- Collision audio follows actual removals: normal turns play one crash cue for the combined Dalek/scrap removal batch; emperor defeat retains its separate cue. Last Stand has its own collision-audio path. Check both modes, scrap hits, and mixed simultaneous collisions when changing collision feedback.
- State flow: menu → playing → level complete → playing, or game over/win. Level completion pauses 1.5 seconds; completing level 10 wins. Reset behavior is state-dependent.
- Browser flow: `index.html` loads `wasm_exec.js`, instantiates `godaleks.wasm`, and starts Go. Keep these files together; the HTML uses relative URLs.

## Key Directories

- `cmd/`: gameplay and rendering; `cmd/assets/` contains embedded PNG/WAV inputs.
- `scripts/serve/`: standalone Go static server for local WASM, including HTTP tests; no build watcher.
- `.github/workflows/`: native checks, Pages deployment, and tagged native releases.
- `web/`, `bin/`, `dist/`: generated local browser bundle, executable, and release outputs; ignored by Git.
- `images/`: README screenshots, separate from embedded gameplay assets.

## Development Commands

Run from the repository root; `Taskfile.yml` is authoritative.

| Command | Purpose |
| --- | --- |
| `task run` | Run native game (`go run ./main.go`) |
| `task build` | Native build to `bin/godaleks.exe`; suffix does not force Windows |
| `task test` | All tests: `go test -v ./...` |
| `task vet` | `go vet ./...` |
| `task lint` | **Modifies files:** `go fmt ./...`, `go mod tidy -v`, `go fix ./...` |
| `task staticcheck` / `task seccheck` | Run Staticcheck / govulncheck through `go run ...@latest` |
| `task wasm:build` | Build/copy browser bundle into `web/` |
| `task wasm:serve` | Serve existing bundle at `http://127.0.0.1:8080` |
| `task wasm` | Build, then serve browser bundle |

Use `go mod download` to fetch pinned dependencies; `task deps` also upgrades them with `go get -u ./...`. `task release` invokes mutating lint first. `task clean` clears the Go build cache and deletes `bin/` and `dist/`.

## Code Conventions & Common Patterns

- Use `gofmt`, exported Go names for public types/constructors, and lower-camel internal helpers/fields. Gameplay changes generally belong in `*Game` receiver methods with early-return guards.
- State management is callback-driven, with concrete dependencies constructed in `NewGame`; there is no application goroutine/channel architecture or DI framework.
- Preserve reusable maps/slices, in-place filtering, precomputed trig tables, cached HUD strings, and shared draw options. Reset reused transforms/color state between sprites; avoid adding allocations in update/draw loops.
- Text uses a shared cached `text/v2.GoXFace` for the bitmap font; the drawing helper converts existing baseline coordinates to text/v2's top origin. Use `vector` primitives and `ColorScale.ScaleAlpha` for premultiplied-alpha fades.
- Keep logical and visual positions separate. Mark `scrapGridDirty` whenever scraps change so `ensureScrapGrid` remains correct. Grid-array lookups are cheap, but some helpers rebuild grids per call.
- Assets use `//go:embed` in `cmd/loadimages.go` and `cmd/loadsounds.go`; runtime asset paths are not external filesystem dependencies. Scrap sprites are generated in `cmd/sprites.go`.
- Existing error handling distinguishes fatal engine startup errors from warning-only asset/audio initialization. Audio decode errors wrap with `%w`; playback is nil-safe and creates fresh players for overlapping sounds. Image warnings do not imply a working fallback.
- Update `CHANGELOG.md` for code or structural changes. Prefer source/configuration over older README or task-history claims when they disagree.

## Important Files

- `cmd/constants.go`: dimensions, states, layout offsets, transition delay.
- `cmd/game.go`: resets, level spawning, occupancy and scrap-grid maintenance.
- `cmd/effect.go`: player actions and effects; `cmd/input.go`: mouse movement and direction arrows.
- `cmd/draw.go`, `cmd/menu.go`: gameplay/HUD and title rendering.
- Release version displays live in `main.go` (window title) and `cmd/menu.go` (splash title); keep them aligned with `Taskfile.yml`, README version references, and release notes.
- `go.mod`, `go.sum`: toolchain/dependencies; `Taskfile.yml`: local workflows.
- `.github/workflows/goreleaser.yml`: actual tag-release builds. `.goreleaser.yaml` configures the separate local snapshot command; they are not equivalent pipelines.

## Runtime/Tooling Preferences

- Current `go.mod` requires Go **1.27.1** and Ebitengine **v2.10.0**. Use Go modules and Task v3; no Node/Bun package toolchain is required. All native build, WASM deployment, and release Go setup steps use `go-version-file: go.mod`; update the module's Go directive to advance CI together.
- Linux native builds/tests need CGO and X11/OpenGL/ALSA development libraries; exact packages are in `.github/workflows/build.yml`. Tag releases use CGO for Linux/macOS and disable it for Windows. WASM targets `GOOS=js GOARCH=wasm`.
- Local WASM builds copy the checked-in `wasm_exec.js`; Pages CI copies it from the active Go SDK. Check compiler/shim compatibility when changing Go versions. CI additionally runs Binaryen optimization; local builds do not.
- The local server supplies WASM MIME and `no-store` headers. Rebuild and reload after changes; use HTTP rather than opening the HTML file directly.

## Testing & QA

- Standard Go `testing`: `cmd/cmd_test.go` covers math/grid helpers; `scripts/serve/main_test.go` covers directory validation and HTTP behavior. There is no configured coverage threshold; helper tests do not prove complete gameplay.
- Follow same-package tests, table-driven subtests, minimal `Game` literals, and server fixtures using `t.TempDir`, `t.Cleanup`, and `httptest`. Minimal game values avoid audio construction, but package initialization still loads Ebiten images.
- Focused examples: `go test ./cmd -run '^TestRebuildScrapGridClearsStaleEntries$' -count=1` and `go test -v ./scripts/serve`.
- Headless Linux: install native dependencies plus Xvfb, then `xvfb-run -a go test ./...`, as Pages CI does. A `-run` filter does not bypass graphics initialization. Server-package tests can run independently of Ebiten.
- For gameplay/render/input changes, run the actual native game or WASM browser build and exercise the affected behavior. Check movement timing, collision outcomes, state transitions, and visual feedback; the existing tests do not cover those end to end.
