# Repository Guidelines

## Project Overview

GoDaleks recreates Macintosh Daleks/BSD Robots in Go with Ebitengine for native desktop and browser WebAssembly. Players survive robot-chase collisions using scrap hazards, teleportation, a sonic screwdriver, and Last Stand. The logical viewport is fixed at 800×600, with a 50×33 grid of 16px cells.

## Architecture & Data Flow

- `main.go` constructs `cmd.NewGame()` and calls `ebiten.RunGame()`. Despite its name, `cmd/` is one gameplay package, not a CLI or service layer.
- `Game` in `cmd/types.go` owns entities, inventory, timers, effects, spatial caches, render buffers, and concrete image/audio dependencies. `cmd/root.go` implements Ebitengine's `Update`, `Draw`, and fixed-size `Layout` callbacks.
- `Update` handles mouse input, advances movement/effects, then dispatches state-specific keyboard input. Simulation uses `1/ebiten.TPS()` with a 60-TPS fallback; some input/HUD timers use wall-clock time. Rendering also updates presentation caches.
- Normal turns move one grid cell, update logical `GridPos` immediately, interpolate `VisualPos` with smootherstep, then resolve collisions when animation finishes. Last Stand instead uses continuous accelerating movement and its own collision passes. Both paths use the level speed multiplier `1 + 0.2 * max(0, level-10)`: unchanged through level 10, 1.2× at 11 through 2× at 15. Normal turns shorten animation without adding grid steps; Last Stand scales movement speed. Check **both paths** when changing movement, combat, or audio.
- `cmd/collision.go` checks player death before enemy/scrap removals, scoring, and progression. Normal collisions protect emperors until ordinary Daleks are gone; Last Stand's direct removal loops do not apply the same immunity. Crash audio follows actual removals; `exterminate` denotes emperor-caused player death, not emperor defeat.
- State flow: menu → playing → level complete → next level, or game over/win. Level transitions pause 1.5 seconds; completing level 15 wins. Restart/reset logic is state-dependent.
- `index.html` launches the game on demand in a `srcdoc` iframe, isolating Ebitengine's canvas/keyboard handling from the landing page. Inside it, `wasm_exec.js` supplies `Go`, then `godaleks.wasm` is fetched, instantiated, and run. Preserve iframe isolation and Tab navigation; keep all three files together because URLs are relative.

## Key Directories

- `cmd/`, `cmd/assets/`: gameplay/rendering and embedded PNG/WAV inputs.
- `scripts/serve/`: standalone Go static server and HTTP handler tests; no build watcher or live reload.
- `.github/workflows/`: native checks, Pages deployment, and tagged native releases.
- `web/`, `bin/`, `dist/`: ignored local browser bundle, executable, and release outputs. Build browser previews into `web/`, not over root deployment artifacts.
- `images/`: README screenshots, separate from game assets. `tasks/` holds historical notes; verify them against current code.

## Development Commands

Run from the repository root; `Taskfile.yml` is authoritative.

| Command | Purpose |
| --- | --- |
| `task run` | Run native game (`go run ./main.go`) |
| `task build` | Build `bin/godaleks.exe`; the suffix does not select Windows |
| `task test` | Run all tests (`go test -v ./...`) |
| `task lint` | **Mutates files:** `go fmt ./...`, `go mod tidy -v`, `go fix ./...` |
| `task vet` | Run `go vet ./...` |
| `task staticcheck` / `task seccheck` | Run Staticcheck / govulncheck using `go run ...@latest` |
| `task wasm:build` | Compile WASM and copy HTML/runtime JS into `web/` |
| `task wasm:serve` | Serve existing `web/` at `http://127.0.0.1:8080` |
| `task wasm` | Build the browser bundle, then serve it |

Use `go mod download` to fetch pinned dependencies; `task deps` also **upgrades** them with `go get -u ./...`. `task release` runs mutating lint and builds locally; it does not publish. `task goreleaser` creates a local snapshot. `task clean` clears the Go build cache and deletes `bin/` and `dist/` on macOS/Linux; Windows uses `task clean_windows`.

## Code Conventions & Common Patterns

- Use `gofmt`, exported Go names for public types/constructors, lower-camel helpers/fields, and `*Game` receiver methods with early-return guards. State is callback-driven; there is no application goroutine/channel architecture or DI framework. Images initialize at package load; `NewGame` constructs audio and reusable state.
- Keep logical and visual positions separate. Whenever scraps change, set `scrapGridDirty = true` so `ensureScrapGrid` remains correct. `positionOccupied` and `isScrapAt` rebuild grids per call; do not assume the whole lookup is constant-time. The `distance` helper returns **squared** distance.
- Preserve reusable maps/slices, in-place filtering, precomputed trig tables, cached HUD strings, and shared draw options. Reset reused `GeoM`/`ColorScale` state between sprites; avoid allocations in update/draw loops.
- Text uses a cached `text/v2.GoXFace`; its helper converts baseline coordinates to top-origin drawing. Use `vector` primitives and `ColorScale.ScaleAlpha` for premultiplied-alpha fades.
- Assets use `//go:embed` in `cmd/loadimages.go` and `cmd/loadsounds.go`; scrap sprites are generated in `cmd/sprites.go`. Game assets need no external runtime asset directory.
- Engine startup errors are fatal; image/audio initialization errors warn. Audio decode errors wrap with `%w`; playback is nil-safe and creates fresh players for overlap. Image warnings do **not** provide a working fallback.
- Update `CHANGELOG.md` for code or structural changes. Prefer current source/configuration over older README or task-history claims.

## Important Files

- `cmd/constants.go`: dimensions, states, layout offsets, transition delay.
- `cmd/game.go`: resets, spawning, occupancy/scrap caches; `cmd/movement.go` and `cmd/collision.go`: movement and combat paths.
- `cmd/effect.go`, `cmd/input.go`: player actions/effects and mouse controls; `cmd/draw.go`, `cmd/menu.go`: gameplay/HUD and title rendering.
- `go.mod`, `go.sum`: toolchain/dependencies; `Taskfile.yml`: local commands.
- Version bumps must synchronize literals in `main.go` and `cmd/menu.go` with `Taskfile.yml`, README references, and release notes. `task generate` currently has no generation directive to update them.
- `.github/workflows/goreleaser.yml` manually builds/publishes `v*` tag releases; it does **not** use `.goreleaser.yaml`, which configures the separate local snapshot command.

## Runtime/Tooling Preferences

- `go.mod` requires **Go 1.27.1** and **Ebitengine v2.10.0**. Use Go modules and Task v3; no Node/Bun package toolchain is required. README's older Go prerequisite is stale. All CI Go setup steps use `go-version-file: go.mod`.
- Linux native builds/tests require CGO and X11/OpenGL/ALSA development libraries; see `.github/workflows/build.yml` for packages. Tagged releases enable CGO for Linux/macOS and disable it for Windows; local GoReleaser configuration differs.
- WASM targets `GOOS=js GOARCH=wasm`. Local builds copy checked-in `wasm_exec.js`; Pages CI copies it from the active Go SDK and additionally runs Binaryen optimization. Check compiler/runtime-shim compatibility after Go upgrades.
- The local server supplies `application/wasm` and `Cache-Control: no-store`. Use HTTP, not `file://`; rebuild and reload after source or root HTML/runtime changes. Serving alone does not rebuild or recopy files.

## Testing & QA

- **Mandatory after every repository change, including documentation:** run `task lint` first, then `task vet`, `task staticcheck`, and `task seccheck` against its resulting state. Fix failures and rerun affected checks. Report actual results; unavailable tools, download failures, and incomplete scans are blockers, not passes. Preserve security findings rather than suppressing them.
- Tests use Go's standard `testing` package: `cmd/cmd_test.go` covers math/grid/cache helpers; `scripts/serve/main_test.go` covers directory validation, MIME/cache headers, traversal, and missing files. No coverage threshold is configured; helper/formula tests do not prove gameplay.
- Follow same-package tests, table-driven subtests, and minimal `Game` literals. Server tests use `t.TempDir`, `t.Cleanup`, and `httptest` without opening a listening port. Minimal game fixtures avoid audio construction but still trigger Ebiten package initialization.
- Focused examples: `go test ./cmd -run '^TestRebuildScrapGridClearsStaleEntries$' -count=1` and `go test -v ./scripts/serve`.
- Headless Linux needs native libraries plus Xvfb: `xvfb-run -a go test ./...`. A `-run` filter does not bypass graphics initialization. Server-package tests run independently of Ebiten.
- Pages CI runs tests; native build CI runs vet, formatting checks, and build, not tests. For gameplay/render/input changes, also run the actual native game or WASM browser build and exercise movement, collisions, transitions, and visual feedback; existing tests do not cover these end to end.
