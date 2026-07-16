// MIT License

// Copyright (c) 2025 - Aaron Saikovski <asaikovski@outlook.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"math"
)

func (g *Game) moveDaleks() {
	// Start movement animation for all daleks
	g.daleksMoving = true

	for i := range g.daleks {
		dalek := &g.daleks[i]

		// Store current visual position as the starting point for animation
		startPos := FloatPosition{
			X: float64(dalek.GridPos.X),
			Y: float64(dalek.GridPos.Y),
		}

		// Calculate new grid position
		dx := 0
		dy := 0

		if dalek.GridPos.X < g.player.X {
			dx = 1
		} else if dalek.GridPos.X > g.player.X {
			dx = -1
		}

		if dalek.GridPos.Y < g.player.Y {
			dy = 1
		} else if dalek.GridPos.Y > g.player.Y {
			dy = -1
		}

		newGridPos := Position{
			X: dalek.GridPos.X + dx,
			Y: dalek.GridPos.Y + dy,
		}

		// Update dalek's positions for smooth animation
		dalek.GridPos = newGridPos
		dalek.StartPos = startPos  // Store the starting position
		dalek.VisualPos = startPos // Start animation from current position
		dalek.TargetPos = FloatPosition{
			X: float64(newGridPos.X),
			Y: float64(newGridPos.Y),
		}
		dalek.IsMoving = true
		dalek.MoveTimer = 0
	}
}

func (g *Game) updateNormalMovement(deltaTime float64) {
	allFinished := true

	for i := range g.daleks {
		dalek := &g.daleks[i]

		if dalek.IsMoving {
			dalek.MoveTimer += deltaTime

			// Calculate interpolation progress (0.0 to 1.0)
			progress := dalek.MoveTimer / g.moveAnimationDuration
			if progress >= 1.0 {
				progress = 1.0
				dalek.IsMoving = false
				dalek.MoveTimer = 0
			} else {
				allFinished = false
			}

			// Ultra-smooth easing function - smootherstep (quintic) interpolation
			// This gives the smoothest possible interpolation with no visible jerk
			easedProgress := progress * progress * progress * (progress*(progress*6.0-15.0) + 10.0)

			// Interpolate position using stored start position for consistency
			startX := dalek.StartPos.X
			startY := dalek.StartPos.Y
			targetX := dalek.TargetPos.X
			targetY := dalek.TargetPos.Y

			// Calculate smooth interpolated position
			dalek.VisualPos.X = startX + (targetX-startX)*easedProgress
			dalek.VisualPos.Y = startY + (targetY-startY)*easedProgress

			// Ensure we end up exactly at target
			if !dalek.IsMoving {
				dalek.VisualPos = dalek.TargetPos
			}
		}
	}

	// Check if all daleks finished moving
	if allFinished {
		g.daleksMoving = false
		g.checkCollisions()
	}
}

func (g *Game) updateLastStandMovement(deltaTime float64) {
	// Don't continue if game is over
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	// Accelerate the movement speed. math.Pow(1.1, dt) ≈ 1 + 0.1*dt at 60 FPS
	// (accel is very close to 1.0 and dt is ~0.016), so a linear approximation
	// is rounding-accurate and avoids a pow() call per frame.
	g.lastStandSpeed *= 1.0 + (g.lastStandAcceleration-1.0)*deltaTime
	if g.lastStandSpeed > g.lastStandMaxSpeed {
		g.lastStandSpeed = g.lastStandMaxSpeed
	}

	anyMoving := false

	// Convert player position to FloatPosition for consistent comparison
	playerPos := FloatPosition{X: float64(g.player.X), Y: float64(g.player.Y)}
	collisionThreshold := 0.5 // Adjust this value to fine-tune collision detection

	// Check player-dalek collisions first
	for _, dalek := range g.daleks {
		if g.checkCollisionWithThreshold(playerPos, dalek.VisualPos, collisionThreshold) {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false
			g.daleksMoving = false
			return
		}
	}

	// Update dalek positions
	for i := range g.daleks {
		dalek := &g.daleks[i]

		dx := playerPos.X - dalek.VisualPos.X
		dy := playerPos.Y - dalek.VisualPos.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > 0.1 { // Still has distance to cover (increased threshold)
			// Normalize direction
			dx /= dist
			dy /= dist

			// Move toward player at current speed
			moveDistance := g.lastStandSpeed * deltaTime
			dalek.VisualPos.X += dx * moveDistance
			dalek.VisualPos.Y += dy * moveDistance

			// Clamp to grid bounds
			dalek.VisualPos.X = math.Max(0, math.Min(float64(gridWidth-1), dalek.VisualPos.X))
			dalek.VisualPos.Y = math.Max(0, math.Min(float64(gridHeight-1), dalek.VisualPos.Y))

			// Update grid position for collision detection
			dalek.GridPos.X = int(math.Round(dalek.VisualPos.X))
			dalek.GridPos.Y = int(math.Round(dalek.VisualPos.Y))

			anyMoving = true
		}
	}

	// Handle collisions after movement (separate pass to avoid slice modification during iteration)
	// Reuse pre-allocated map
	for k := range g.dalekRemoveMap {
		delete(g.dalekRemoveMap, k)
	}
	g.newScrapsBuf = g.newScrapsBuf[:0] // Reuse pre-allocated buffer
	hasAnyCollision := false            // Track if any collision occurred

	// Rebuild the scrap grid once per frame so the O(1) lookup is fresh.
	g.rebuildScrapGrid()

	// Check for collisions with scraps using the O(1) scrapGrid lookup.
	// Dalek's visual position is rounded to the nearest grid cell; if that
	// cell contains a scrap, it's a collision.
	for i, dalek := range g.daleks {
		if g.dalekRemoveMap[i] {
			continue
		}
		gx := int(math.Round(dalek.VisualPos.X))
		gy := int(math.Round(dalek.VisualPos.Y))
		if gx < 0 || gx >= gridWidth || gy < 0 || gy >= gridHeight {
			continue
		}
		if g.scrapGrid[gx][gy] {
			g.dalekRemoveMap[i] = true
			g.score += 2
			hasAnyCollision = true
		}
	}

	// Check for dalek-dalek collisions
	for i := 0; i < len(g.daleks); i++ {
		if g.dalekRemoveMap[i] {
			continue
		}
		for j := i + 1; j < len(g.daleks); j++ {
			if g.dalekRemoveMap[j] {
				continue
			}
			if g.checkCollisionWithThreshold(g.daleks[i].VisualPos, g.daleks[j].VisualPos, collisionThreshold) {
				g.dalekRemoveMap[i] = true
				g.dalekRemoveMap[j] = true
				g.score += 4 // 2 points per dalek
				hasAnyCollision = true
				collisionPos := Position{
					X: int((g.daleks[i].VisualPos.X + g.daleks[j].VisualPos.X) / 2),
					Y: int((g.daleks[i].VisualPos.Y + g.daleks[j].VisualPos.Y) / 2),
				}
				if !g.positionOccupied(collisionPos) {
					g.newScrapsBuf = append(g.newScrapsBuf, collisionPos)
				}
				break
			}
		}
	}

	// Play crash sound for any collision (dalek-dalek OR dalek-scrap)
	if hasAnyCollision && g.soundPlayer != nil {
		g.soundPlayer.Play("crash")
	}

	// Remove collided daleks using reusable buffer
	if len(g.dalekRemoveMap) > 0 {
		g.survivingDaleksBuf = g.survivingDaleksBuf[:0]
		for i, dalek := range g.daleks {
			if !g.dalekRemoveMap[i] {
				g.survivingDaleksBuf = append(g.survivingDaleksBuf, dalek)
			}
		}
		g.daleks = append(g.daleks[:0], g.survivingDaleksBuf...)
	}

	// Add new scraps and collision effects
	if len(g.newScrapsBuf) > 0 {
		g.scraps = append(g.scraps, g.newScrapsBuf...)

		// Add collision explosion effect for each new scrap
		for _, scrapPos := range g.newScrapsBuf {
			g.collisionEffects = append(g.collisionEffects, CollisionEffect{
				Pos:      FloatPosition{X: float64(scrapPos.X), Y: float64(scrapPos.Y)},
				Timer:    0,
				Duration: 0.6, // 0.6 second explosion animation
			})
		}
	}

	// Check collisions every few frames for better performance
	// but still frequent enough for good responsiveness
	if int(g.lastStandSpeed*deltaTime*60)%2 == 0 {
		g.checkCollisions()
	}

	// If player died during collision check, end Last Stand immediately
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	// End Last Stand if no daleks are moving or they're all gone
	if !anyMoving || len(g.daleks) == 0 {
		g.isLastStandActive = false
		g.daleksMoving = false
		if g.state == StatePlaying && len(g.daleks) == 0 {
			g.score += 50 // Bonus for surviving Last Stand
		}
	}
}
