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
	"math/rand"
)

// resetGame resets the game to initial state and starts from level 1
func (g *Game) resetGame() {
	g.level = 1
	g.score = 0
	g.teleports = 10
	g.safeTeleports = 3
	g.screwdrivers = 2
	g.lastStands = 1
	g.teleportAnimation = false
	g.teleportTimer = 0
	g.screwdriverAnimation = false
	g.screwdriverTimer = 0
	g.screwdriverTargets = nil
	g.daleksMoving = false
	g.isLastStandActive = false
	g.lastStandSpeed = 2.0
	g.daleks = nil
	g.scraps = nil
	g.gameOverMessage = ""
	g.startLevel()
}

func (g *Game) startLevel() {
	// Clear the board and reset all states
	g.scraps = nil
	g.daleksMoving = false
	g.isLastStandActive = false
	g.lastStandSpeed = 2.0
	g.teleportAnimation = false
	g.teleportTimer = 0
	g.screwdriverAnimation = false
	g.screwdriverTimer = 0
	g.screwdriverTargets = nil
	g.lastStands = 1

	// Place player randomly
	g.player = Position{
		X: rand.Intn(gridWidth),
		Y: rand.Intn(gridHeight),
	}

	// Place daleks (5 + level number)
	dalekCount := 5 + g.level
	g.daleks = make([]Dalek, 0, dalekCount)

	for len(g.daleks) < dalekCount {
		pos := Position{
			X: rand.Intn(gridWidth),
			Y: rand.Intn(gridHeight),
		}

		// Don't place dalek on player or too close
		if g.distance(pos, g.player) > 3 && !g.positionOccupied(pos) {
			dalek := Dalek{
				GridPos:   pos,
				VisualPos: FloatPosition{X: float64(pos.X), Y: float64(pos.Y)},
				TargetPos: FloatPosition{X: float64(pos.X), Y: float64(pos.Y)},
				IsMoving:  false,
				MoveTimer: 0,
			}
			g.daleks = append(g.daleks, dalek)
		}
	}

	g.state = StatePlaying
	g.soundPlayer.Play("gamestart")
}

func (g *Game) distance(a, b Position) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return dx*dx + dy*dy // Using squared distance for efficiency
}

func (g *Game) positionOccupied(pos Position) bool {
	for _, dalek := range g.daleks {
		if dalek.GridPos == pos {
			return true
		}
	}
	for _, scrap := range g.scraps {
		if scrap == pos {
			return true
		}
	}
	return false
}

// Convert screen coordinates to grid coordinates
func (g *Game) screenToGrid(screenX, screenY int) (int, int, bool) {
	offsetX := (screenWidth - gridWidth*cellSize) / 2
	offsetY := 50

	gridX := (screenX - offsetX) / cellSize
	gridY := (screenY - offsetY) / cellSize

	// Check if within grid bounds
	if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
		return gridX, gridY, true
	}
	return 0, 0, false
}

func (g *Game) isSafePosition(pos Position) bool {
	// Check if any dalek can reach this position in one move
	for _, dalek := range g.daleks {
		if g.distance(pos, dalek.GridPos) <= 2 { // Within one move
			return false
		}
	}
	return true
}

func (g *Game) lastStand() {
	if g.state != StatePlaying || g.lastStands <= 0 || g.daleksMoving {
		return
	}

	g.lastStands--
	g.isLastStandActive = true
	g.lastStandSpeed = 2.0 // Reset speed to starting value
	g.daleksMoving = true  // Enable daleks movement for Last Stand
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
