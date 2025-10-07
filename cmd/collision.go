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

// Add this function to check for collision with a threshold
func (g *Game) checkCollisionWithThreshold(pos1, pos2 FloatPosition, threshold float64) bool {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	distSquared := dx*dx + dy*dy
	return distSquared < threshold*threshold
}

func (g *Game) checkCollisions() {
	// Early exit if game is already over
	if g.state != StatePlaying {
		return
	}

	// Check player-dalek collision FIRST (most important check)
	for _, dalek := range g.daleks {
		if g.player == dalek.GridPos {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false // End Last Stand immediately
			g.daleksMoving = false
			return
		}
	}

	// Pre-allocate maps with estimated capacity
	scrapMap := make(map[Position]bool, len(g.scraps))
	for _, scrap := range g.scraps {
		scrapMap[scrap] = true
	}

	// Check dalek-dalek and dalek-scrap collisions
	collidedPositions := make(map[Position]bool, len(g.daleks)/4) // Estimate ~25% collisions
	positionCounts := make(map[Position]int, len(g.daleks))

	// Single pass: count daleks at each position and check scrap collisions
	for _, dalek := range g.daleks {
		if scrapMap[dalek.GridPos] {
			// Dalek hit scrap
			g.score += 2
			collidedPositions[dalek.GridPos] = true
		} else {
			positionCounts[dalek.GridPos]++
		}
	}

	// Identify dalek-dalek collisions (positions with count > 1)
	playedCrashSound := false
	for pos, count := range positionCounts {
		if count > 1 {
			g.score += 2 * count // 2 points per dalek
			collidedPositions[pos] = true
			if !playedCrashSound {
				g.soundPlayer.Play("crash")
				playedCrashSound = true
			}
		}
	}

	// Build final dalek list excluding collided ones
	finalDaleks := make([]Dalek, 0, len(g.daleks)-len(collidedPositions))
	for _, dalek := range g.daleks {
		if !collidedPositions[dalek.GridPos] {
			finalDaleks = append(finalDaleks, dalek)
		}
	}

	// Add new scraps for collided positions (avoid duplicates using scrapMap)
	for pos := range collidedPositions {
		if !scrapMap[pos] {
			g.scraps = append(g.scraps, pos)
			scrapMap[pos] = true // Update map to prevent future duplicates
		}
	}

	g.daleks = finalDaleks

	// Check player-dalek collision again after updating dalek positions
	// (in case daleks moved onto player during collision resolution)
	for _, dalek := range g.daleks {
		if g.player == dalek.GridPos {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false // End Last Stand immediately
			g.daleksMoving = false
			return
		}
	}

	// Check if level is complete
	if len(g.daleks) == 0 {
		g.score += g.level * 10
		g.level++
		g.teleports += 2

		g.screwdrivers += 2 // Increase screwdrivers by 2 every level
		if g.level%5 == 0 { // Bonus last stand every 5 levels
			g.lastStands++
		}
		if g.level > 10 {
			g.state = StateWin
			g.gameOverMessage = "Congratulations! You survived all levels!"
			g.soundPlayer.Play("gameover")
		} else {
			g.startLevel()
		}
	}
}
