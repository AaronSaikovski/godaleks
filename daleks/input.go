package daleks

import (
	"time"
)

// Handle mouse click for player movement
func (g *Game) handleMouseClick(x, y int) {
	// Prevent too rapid clicking
	if time.Since(g.lastClickTime) < 100*time.Millisecond {
		return
	}
	g.lastClickTime = time.Now()

	// Convert to grid coordinates
	gridX, gridY, valid := g.screenToGrid(x, y)
	if !valid {
		return
	}

	targetPos := Position{X: gridX, Y: gridY}

	// Check if clicking on current player position (stay in place)
	if targetPos == g.player {
		if !g.isLastStandActive {
			g.moveDaleks()
		}
		return
	}

	// Calculate movement direction
	dx := 0
	dy := 0

	if targetPos.X > g.player.X {
		dx = 1
	} else if targetPos.X < g.player.X {
		dx = -1
	}

	if targetPos.Y > g.player.Y {
		dy = 1
	} else if targetPos.Y < g.player.Y {
		dy = -1
	}

	// Only allow one-step moves (adjacent cells including diagonal)
	if abs(dx) <= 1 && abs(dy) <= 1 {
		g.movePlayer(dx, dy)
	}
}
