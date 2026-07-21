# CHANGELOG

## v1.2.2. (2026-07-16)
- Added the classic direction arrows: a bold black arrow in each of the 8 cells around the human player, shown only for valid moves (in-bounds and not blocked by scrap). Arrows hide while the Daleks move and reappear around the player's new position, matching the original.
- Fixed jerky Dalek movement — `Update` now uses a fixed timestep (`1/TPS`) instead of a wall-clock deltaTime, giving even smootherstep-eased animation progress.
- Removed the "Speed: X.X" counter from the Last Stand HUD indicator (now just "LAST STAND ACTIVE!").
- Security: bumped the `go` directive to 1.26.5 for the patched `crypto/tls` stdlib (GO-2026-5856).
- Cleanups: removed dead `lastUpdateTime`/`cachedLastStand*` state and resolved two staticcheck findings (S1021, SA4006).
- Investigated Dalek/human grid centering — confirmed sprites are already pixel-perfect centered; no change needed.

## v0.0.4. (2025-08-12)
- Bug fixes and performance improvements.
- When using sonic screwdriver, not leaving behind a debris field.
- Smoother Dalek movement, less jerky.

## v0.0.3. (2025-08-11)
- Bug fixes and performance improvements.
- Increase screwdrivers by 2 every level

## v0.0.2. (2025-08-11)
- Added New game control.
- Added Sound effects.
- Bug fixes and performance improvements.

## v0.0.1. (2025-08-08)
- initial alpha release version
