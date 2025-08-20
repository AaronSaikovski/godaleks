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

package daleks

import (
	"math"
)

func (g *Game) moveDaleks() {
	// Start movement animation for all daleks
	g.daleksMoving = true

	for i := range g.daleks {
		dalek := &g.daleks[i]

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

			// Much smoother easing function - smoothstep interpolation
			// This gives a very smooth start and end with consistent middle speed
			var easedProgress float64
			easedProgress = progress * progress * (3.0 - 2.0*progress)

			// Alternative: Even smoother with smootherstep (quintic) - uncomment to try
			// easedProgress = progress * progress * progress * (progress*(progress*6.0-15.0)+10.0)

			// Calculate current visual position
			startX := dalek.VisualPos.X
			startY := dalek.VisualPos.Y
			targetX := dalek.TargetPos.X
			targetY := dalek.TargetPos.Y

			// Interpolate position
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

	// Accelerate the movement speed
	g.lastStandSpeed *= math.Pow(g.lastStandAcceleration, deltaTime)
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

			// Store old position for collision check
			oldPos := dalek.VisualPos

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

			// Check for collisions with scraps
			for _, scrap := range g.scraps {
				scrapPos := FloatPosition{X: float64(scrap.X), Y: float64(scrap.Y)}
				if g.checkCollisionWithThreshold(dalek.VisualPos, scrapPos, collisionThreshold) {
					dalek.VisualPos = oldPos // Prevent moving through scraps
					g.daleks = append(g.daleks[:i], g.daleks[i+1:]...)
					g.score += 2
					g.soundPlayer.Play("crash")
					// Don't add duplicate scraps
					if !g.positionOccupied(Position{X: int(oldPos.X), Y: int(oldPos.Y)}) {
						g.scraps = append(g.scraps, Position{X: int(oldPos.X), Y: int(oldPos.Y)})
					}
					return
				}
			}

			// Check for collisions with other daleks
			for j := i + 1; j < len(g.daleks); j++ {
				if g.checkCollisionWithThreshold(dalek.VisualPos, g.daleks[j].VisualPos, collisionThreshold) {
					collisionPos := Position{
						X: int((dalek.VisualPos.X + g.daleks[j].VisualPos.X) / 2),
						Y: int((dalek.VisualPos.Y + g.daleks[j].VisualPos.Y) / 2),
					}
					g.daleks = append(g.daleks[:i], g.daleks[i+1:]...)
					g.daleks = append(g.daleks[:j-1], g.daleks[j:]...)
					g.score += 4 // 2 points per dalek
					g.soundPlayer.Play("crash")
					if !g.positionOccupied(collisionPos) {
						g.scraps = append(g.scraps, collisionPos)
					}
					return
				}
			}
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
