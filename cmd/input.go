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
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (g *Game) drawMouseIndicator(screen *ebiten.Image) {
	if g.state != StatePlaying {
		return
	}

	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Convert to grid coordinates
	gridX, gridY, valid := g.screenToGrid(mouseX, mouseY)
	if !valid {
		return
	}

	targetPos := Position{X: gridX, Y: gridY}

	// Check if it's a valid move (adjacent to player)
	dx := abs(targetPos.X - g.player.X)
	dy := abs(targetPos.Y - g.player.Y)

	// Only show indicator for valid moves or current position
	if dx <= 1 && dy <= 1 {
		offsetX := (screenWidth - gridWidth*cellSize) / 2
		offsetY := 50

		x := float64(offsetX + gridX*cellSize)
		y := float64(offsetY + gridY*cellSize)

		// Choose color based on move type
		var indicatorColor color.Color
		if targetPos == g.player {
			indicatorColor = color.RGBA{0, 255, 0, 100} // Green for wait/current position
		} else {
			// Check if position is occupied by scrap
			occupied := false
			for _, scrap := range g.scraps {
				if scrap == targetPos {
					occupied = true
					break
				}
			}

			if occupied {
				indicatorColor = color.RGBA{255, 0, 0, 100} // Red for blocked
			} else {
				indicatorColor = color.RGBA{0, 0, 255, 100} // Blue for valid move
			}
		}

		// Draw semi-transparent overlay on the cell
		ebitenutil.DrawRect(screen, x, y, cellSize, cellSize, indicatorColor)

		// Draw border
		ebitenutil.DrawRect(screen, x, y, cellSize, 1, color.Black)
		ebitenutil.DrawRect(screen, x, y, 1, cellSize, color.Black)
		ebitenutil.DrawRect(screen, x+cellSize-1, y, 1, cellSize, color.Black)
		ebitenutil.DrawRect(screen, x, y+cellSize-1, cellSize, 1, color.Black)
	}
}
