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

	// Check dalek-dalek and dalek-scrap collisions
	newDaleks := make([]Dalek, 0, len(g.daleks))
	collidedPositions := make(map[Position]bool)

	// First pass: check for collisions with scraps
	for _, dalek := range g.daleks {
		collided := false

		// Check collision with scraps
		for _, scrap := range g.scraps {
			if dalek.GridPos == scrap {
				collided = true
				g.soundPlayer.Play("crash")
				g.score += 2
				collidedPositions[dalek.GridPos] = true
				break
			}
		}

		if !collided {
			newDaleks = append(newDaleks, dalek)
		}
	}

	// Second pass: check for dalek-dalek collisions
	finalDaleks := make([]Dalek, 0, len(newDaleks))

	for i, dalek := range newDaleks {
		collided := false

		// Check if this dalek collides with any other dalek
		for j, other := range newDaleks {
			if i != j && dalek.GridPos.X == other.GridPos.X && dalek.GridPos.Y == other.GridPos.Y {
				collided = true
				g.score += 2
				collidedPositions[dalek.GridPos] = true
				g.soundPlayer.Play("crash")
				break
			}
		}

		if !collided {
			finalDaleks = append(finalDaleks, dalek)
		}

	}

	// Add debris piles for all collided positions
	for pos := range collidedPositions {
		// Check if scrap already exists at this position
		scrapExists := false
		for _, scrap := range g.scraps {
			if scrap == pos {
				scrapExists = true
				break
			}
		}
		if !scrapExists {
			g.scraps = append(g.scraps, pos)
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
