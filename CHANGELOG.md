# CHANGELOG

## Unreleased
- Added Google Analytics (gtag.js, ID `G-1Q92MBTYT9`) to `index.html`, placed immediately after the `<head>` element for the WASM build's GitHub Pages deployment.
- SEO: added a full set of meta tags to `index.html` — `lang`, `charset`, responsive `viewport`, descriptive `<title>`, meta description/keywords/author/robots, `theme-color`, and a canonical URL.
- SEO: added Open Graph and Twitter Card tags plus JSON-LD `VideoGame` structured data for rich social/search previews (using the repo's `dalek.png` sprite as the preview image).
- SEO/accessibility: added a self-contained inline SVG favicon, a visually-hidden semantic `<h1>`/description header, a `<noscript>` fallback, and switched the game container to a `<main>` element.
- SEO: added `robots.txt` and `sitemap.xml`, and wired both into the GitHub Pages deploy workflow (`deploy-wasm.yml`).
- Hosting: set the custom domain to `https://godaleks.com/` — pointed all canonical/OG/sitemap/robots URLs at it and added a `CNAME` file (also copied during deploy).
- SEO: expanded the `index.html` keyword targeting for "daleks", "dalek game", "play daleks online" and related search terms (meta `keywords` + JSON-LD `keywords`).
- SEO: replaced the visually-hidden header with a real, visible `<article>` landing section (H1/H2s, "How to play", "About") — search engines discount hidden text, so genuine on-page content ranks far better for terms like "daleks game".
- SEO: enriched the JSON-LD `VideoGame` structured data with `alternateName` ("Daleks", "Daleks Game", "Robots Game"), `keywords`, `inLanguage`, `playMode`, and `applicationCategory: GameApplication` (removed the duplicate `applicationCategory`).
- SEO: added a `lastmod` (2026-07-22) to `sitemap.xml` so crawlers get a page-freshness signal.

## v1.2.2. (2026-07-16)
- Added the classic direction arrows: a bold black arrow in each of the 8 cells around the human player, shown only for valid moves (in-bounds and not blocked by scrap). Arrows hide while the Daleks move and reappear around the player's new position, matching the original.
- Removed the green highlight square that appeared over the player's own cell on mouse hover; the mouse indicator now only highlights adjacent move cells (blue = valid, red = blocked).
- Performance: the scrap grid now rebuilds only when scraps actually change (dirty flag + `ensureScrapGrid`) instead of on every `isScrapAt` call, so the per-frame arrow overlay does O(1) lookups with no per-frame rebuild. Added dirty-flag unit tests.
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
