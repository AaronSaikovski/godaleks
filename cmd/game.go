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
	g.emperorWarningMessage = ""
	g.emperorWarningTimer = 0
	g.cachedFinalScore = ""
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
	g.emperorWarningMessage = ""
	g.emperorWarningTimer = 0

	// Place player randomly
	g.player = Position{
		X: rand.Intn(gridWidth),
		Y: rand.Intn(gridHeight),
	}

	// Place daleks (5 + level number)
	dalekCount := 5 + g.level
	maxDaleks := dalekCount

	// Determine if emperor should spawn (only from level 2 onwards, with 60% chance)
	spawnEmperor := g.level >= 2 && rand.Float64() < 0.6
	if spawnEmperor {
		maxDaleks++ // Add one more slot for the emperor
	}

	g.daleks = make([]Dalek, 0, maxDaleks)

	// Spawn the emperor first if it should appear this level
	if spawnEmperor {
		for {
			pos := Position{
				X: rand.Intn(gridWidth),
				Y: rand.Intn(gridHeight),
			}

			// Don't place emperor on player or too close
			if g.distance(pos, g.player) > 3 && !g.positionOccupied(pos) {
				floatPos := FloatPosition{X: float64(pos.X), Y: float64(pos.Y)}
				emperor := Dalek{
					GridPos:   pos,
					VisualPos: floatPos,
					StartPos:  floatPos,
					TargetPos: floatPos,
					IsMoving:  false,
					MoveTimer: 0,
					IsEmperor: true,
				}
				g.daleks = append(g.daleks, emperor)

				// Show emperor warning message
				g.emperorWarningMessage = "WARNING: DALEK EMPEROR DETECTED!"
				g.emperorWarningTimer = 0

				// Play emperor alert sound
				g.soundPlayer.Play("dalek_emperor")
				break
			}
		}
	}

	// Spawn normal daleks
	for len(g.daleks) < maxDaleks {
		pos := Position{
			X: rand.Intn(gridWidth),
			Y: rand.Intn(gridHeight),
		}

		// Don't place dalek on player or too close
		if g.distance(pos, g.player) > 3 && !g.positionOccupied(pos) {
			floatPos := FloatPosition{X: float64(pos.X), Y: float64(pos.Y)}
			dalek := Dalek{
				GridPos:   pos,
				VisualPos: floatPos,
				StartPos:  floatPos,
				TargetPos: floatPos,
				IsMoving:  false,
				MoveTimer: 0,
				IsEmperor: false,
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
	// Use grid array for O(1) lookups with zero GC pressure
	g.rebuildOccupancyGrid()
	return g.occupancyGrid[pos.X][pos.Y]
}

func (g *Game) rebuildOccupancyGrid() {
	// Clear grid
	g.occupancyGrid = [gridWidth][gridHeight]bool{}

	// Mark dalek positions
	for _, dalek := range g.daleks {
		if dalek.GridPos.X >= 0 && dalek.GridPos.X < gridWidth &&
			dalek.GridPos.Y >= 0 && dalek.GridPos.Y < gridHeight {
			g.occupancyGrid[dalek.GridPos.X][dalek.GridPos.Y] = true
		}
	}
	// Mark scrap positions
	for _, scrap := range g.scraps {
		if scrap.X >= 0 && scrap.X < gridWidth &&
			scrap.Y >= 0 && scrap.Y < gridHeight {
			g.occupancyGrid[scrap.X][scrap.Y] = true
		}
	}
}

func (g *Game) rebuildScrapGrid() {
	g.scrapGrid = [gridWidth][gridHeight]bool{}
	for _, scrap := range g.scraps {
		if scrap.X >= 0 && scrap.X < gridWidth &&
			scrap.Y >= 0 && scrap.Y < gridHeight {
			g.scrapGrid[scrap.X][scrap.Y] = true
		}
	}
}

func (g *Game) isScrapAt(pos Position) bool {
	if pos.X < 0 || pos.X >= gridWidth || pos.Y < 0 || pos.Y >= gridHeight {
		return false
	}
	g.rebuildScrapGrid()
	return g.scrapGrid[pos.X][pos.Y]
}

// Convert screen coordinates to grid coordinates
func (g *Game) screenToGrid(screenX, screenY int) (int, int, bool) {
	gridX := (screenX - gridOffsetX) / cellSize
	gridY := (screenY - gridOffsetY) / cellSize

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
