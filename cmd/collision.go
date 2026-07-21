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

import "fmt"

// Add this function to check for collision with a threshold
func (g *Game) checkCollisionWithThreshold(pos1, pos2 FloatPosition, threshold float64) bool {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	distSquared := dx*dx + dy*dy
	return distSquared < threshold*threshold
}

func (g *Game) countNormalDaleks() int {
	count := 0
	for _, dalek := range g.daleks {
		if !dalek.IsEmperor {
			count++
		}
	}
	return count
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

	// Reuse pre-allocated scrap map
	for k := range g.scrapMap {
		delete(g.scrapMap, k)
	}
	for _, scrap := range g.scraps {
		g.scrapMap[scrap] = true
	}

	// Reuse pre-allocated maps for collision detection
	for k := range g.daleksByPosition {
		delete(g.daleksByPosition, k)
	}
	for k := range g.dalekRemoveMap {
		delete(g.dalekRemoveMap, k)
	}

	// Map daleks by their position
	for i, dalek := range g.daleks {
		if g.scrapMap[dalek.GridPos] {
			// Dalek hit scrap - mark for removal (emperor immune to scraps)
			if !dalek.IsEmperor {
				g.score += 2
				g.dalekRemoveMap[i] = true
			}
		} else {
			// Add to position map for collision detection
			g.daleksByPosition[dalek.GridPos] = append(g.daleksByPosition[dalek.GridPos], i)
		}
	}

	// Identify dalek-dalek collisions
	playedCrashSound := false
	for _, indices := range g.daleksByPosition {
		if len(indices) > 1 {
			// Multiple daleks at this position - check for emperor
			hasEmperor := false
			for _, idx := range indices {
				if g.daleks[idx].IsEmperor {
					hasEmperor = true
					break
				}
			}

			if hasEmperor {
				// Emperor at this position - only normal daleks die
				for _, idx := range indices {
					if !g.daleks[idx].IsEmperor {
						g.dalekRemoveMap[idx] = true
						g.score += 2 // 2 points per normal dalek destroyed by emperor
					}
				}
			} else {
				// Only normal daleks at this position - all die
				for _, idx := range indices {
					g.dalekRemoveMap[idx] = true
					g.score += 2 // 2 points per dalek
				}
			}

			if !playedCrashSound && g.soundPlayer != nil {
				g.soundPlayer.Play("crash")
				playedCrashSound = true
			}
		}
	}

	// Filter daleks in-place using reusable buffer
	g.survivingDaleksBuf = g.survivingDaleksBuf[:0]
	for i, dalek := range g.daleks {
		if !g.dalekRemoveMap[i] {
			g.survivingDaleksBuf = append(g.survivingDaleksBuf, dalek)
		}
	}

	// Add new scraps for dalek collisions (avoid duplicates using scrapMap)
	for pos, indices := range g.daleksByPosition {
		if len(indices) > 1 && !g.scrapMap[pos] {
			g.scraps = append(g.scraps, pos)
			g.scrapGridDirty = true
			g.scrapMap[pos] = true // Update map to prevent future duplicates

			// Add collision explosion effect
			g.collisionEffects = append(g.collisionEffects, CollisionEffect{
				Pos:      FloatPosition{X: float64(pos.X), Y: float64(pos.Y)},
				Timer:    0,
				Duration: 0.6, // 0.6 second explosion animation
			})
		}
	}

	// Swap daleks with surviving buffer
	g.daleks = append(g.daleks[:0], g.survivingDaleksBuf...)

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

	// If all normal daleks are defeated, the emperor dies too
	normalDalekCount := g.countNormalDaleks()
	if normalDalekCount == 0 && len(g.daleks) > 0 {
		// Check if there's an emperor alive
		emperorFound := false
		emperorIndex := -1
		for i, dalek := range g.daleks {
			if dalek.IsEmperor {
				emperorFound = true
				emperorIndex = i
				break
			}
		}

		// Remove the emperor since all normal daleks are dead
		if emperorFound {
			g.score += 50                                // Bonus points for defeating the emperor
			emperorPos := g.daleks[emperorIndex].GridPos // Store position before removal

			g.survivingDaleksBuf = g.survivingDaleksBuf[:0]
			for i, dalek := range g.daleks {
				if i != emperorIndex {
					g.survivingDaleksBuf = append(g.survivingDaleksBuf, dalek)
				}
			}
			g.daleks = append(g.daleks[:0], g.survivingDaleksBuf...)

			// Add explosion effect at emperor's position
			g.collisionEffects = append(g.collisionEffects, CollisionEffect{
				Pos:      FloatPosition{X: float64(emperorPos.X), Y: float64(emperorPos.Y)},
				Timer:    0,
				Duration: 0.6,
			})
			if g.soundPlayer != nil {
				g.soundPlayer.Play("crash") // Play crash sound for emperor defeat
			}
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
			g.state = StateLevelComplete
			g.levelCompleteTimer = 0
			g.cachedLevelMsg = fmt.Sprintf("Level %d Complete!", g.level-1)
			g.cachedLevelNextMsg = fmt.Sprintf("Starting Level %d...", g.level)
		}
	}
}
