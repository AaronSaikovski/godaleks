# Task: Instant Dalek movement + player direction arrows

## Decisions (from user)
- **Movement**: Match the original classic Mac Daleks — instant/snap, no interpolation.
- **Arrows**: Show 8-direction arrows around the player, but only for *valid* moves
  (in-bounds and not scrap). Must persist/reappear at the new position after each move.

## Plan

### A. Instant Dalek movement (remove interpolation)
- [ ] `movement.go` — rewrite `moveDaleks()` to set each Dalek's `GridPos`/`VisualPos`
      directly to the new cell, then call `checkCollisions()` immediately. No `daleksMoving` window.
- [ ] `movement.go` — delete now-dead `updateNormalMovement()`.
- [ ] `effect.go` — simplify `updateDalekAnimations()` to only drive Last Stand.
- [ ] Remove orphaned animation state: `Dalek.StartPos/TargetPos/IsMoving/MoveTimer` (`types.go`),
      `Game.moveAnimationDuration` (`types.go`), their inits (`game.go`, `root.go`).
- [ ] Keep `VisualPos` — still used for rendering and Last Stand sub-cell motion.

### B. Player direction arrows (valid-moves only)
- [ ] Add `drawPlayerArrows(screen)` — for each of 8 (dx,dy), if target cell is in-bounds
      and not scrap, draw an outward arrow (chevron via `ebitenutil.DrawLine`) in that neighbor cell.
- [ ] Call it from `Draw` in the `StatePlaying` case (drawn every frame → reappears after moves).
- [ ] Skip during Last Stand.

### C. Verify + document
- [ ] `task build`, `go vet`, `task test`, `task staticcheck`.
- [ ] Run the game to confirm instant movement + arrows.
- [ ] Update `CHANGELOG.md` and README changelog under v1.2.2.

## Course correction (after user testing)
- Instant movement felt **too jerky** → reverted to the smooth smootherstep interpolation
  (0.7s per step, fixed timestep). Restored `Dalek` animation fields, `moveAnimationDuration`,
  `moveDaleks()`, `updateNormalMovement()`, and the normal branch of `updateDalekAnimations()`.
- Arrows restyled to match the original: **bold** (2px anti-aliased strokes via `vector.StrokeLine`),
  and colour changed blue → **black** per request. Now also hidden while `daleksMoving` so they
  disappear on click and reappear once the move settles.

## Review

### Done
- Instant Dalek movement: `moveDaleks()` now snaps each Dalek one cell and calls
  `checkCollisions()` in the same step. Deleted `updateNormalMovement()`; simplified
  `updateDalekAnimations()` to only drive Last Stand. Removed orphaned animation state
  (`StartPos`/`TargetPos`/`IsMoving`/`MoveTimer`, `moveAnimationDuration`) and its inits.
- Player arrows: `drawPlayerArrows()` + `drawArrow()` in `input.go`, called from `Draw`
  (StatePlaying). Blue chevrons in each of the 8 neighbour cells, only for in-bounds,
  non-scrap targets; skipped during Last Stand. Redrawn every frame → reappear after moves.
- `VisualPos` retained (used by renderer and Last Stand sub-cell motion).
- CHANGELOG.md + README changelog updated under v1.2.2.

### Verification
- `go build ./...`, `go vet ./...`, `go test ./cmd/...` all pass.
- No references to removed symbols remain (grep clean).
- Game launches without panic. NOTE: could not visually confirm gameplay in this
  environment (screencapture blocked / no display access) — user is testing locally.

### Notes / follow-ups
- Last Stand is a GoDaleks-original power mode (not in classic Daleks); left as continuous
  chase movement by design.
- Arrows treat moving onto a Dalek cell as "valid" (legal but fatal), matching the original;
  only scrap/walls are hidden.
