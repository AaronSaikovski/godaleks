# Lessons

## Game-feel decisions need a real play-test before committing hard
- **Context**: User chose "instant Dalek movement" to match the original video; after
  testing they found it too jerky and asked to restore smooth movement.
- **Rule**: When a change alters *game feel* (movement timing, easing, animation), treat the
  chosen option as provisional until the user has actually run it. Keep the reverted-from code
  easy to restore (small, contained diffs), and don't delete the alternative implementation's
  logic until the new behaviour is confirmed in-game. Offer a quick "keep both / make it a
  constant" path when the tradeoff is subjective.

## Match visual style to the reference explicitly (colour + weight)
- **Context**: First pass drew thin 1px blue chevrons; user wanted bold arrows, then black.
- **Rule**: For pixel-art / retro UI, zoom into the reference asset (sips upscale) before
  drawing. Confirm **colour** and **stroke weight** against the source — don't assume a colour
  from an in-game screenshot is the intended one. Use `vector.StrokeLine` (anti-aliased, width)
  for bold clean strokes rather than 1px `ebitenutil.DrawLine`.
