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
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// arrowColor is the colour of the player's direction markers. Declared once as a
// color.Color so passing it to drawArrow doesn't box a fresh interface value per
// frame. Black, matching the original Daleks arrow markers.
var arrowColor color.Color = color.RGBA{0, 0, 0, 255}

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

// drawPlayerArrows draws the classic Daleks direction indicators: a bold black
// arrow in each of the 8 cells around the human player, but only for moves that
// are actually valid (in-bounds and not blocked by scrap). The arrows are hidden
// while the Daleks are mid-move and redrawn once settled, so — like the original
// — they disappear on a click and reappear around the player's new position.
func (g *Game) drawPlayerArrows(screen *ebiten.Image) {
	// No movement choices to show while Daleks are moving or during Last Stand.
	if g.state != StatePlaying || g.isLastStandActive || g.daleksMoving {
		return
	}

	// Ensure the scrap grid is current, then read it directly below — avoids the
	// full grid rebuild that isScrapAt() would trigger on each of the (up to) 8
	// lookups, and ensureScrapGrid only does work when the scraps actually changed.
	g.ensureScrapGrid()

	pcx := float64(gridOffsetX) + (float64(g.player.X)+0.5)*float64(cellSize)
	pcy := float64(gridOffsetY) + (float64(g.player.Y)+0.5)*float64(cellSize)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			target := Position{X: g.player.X + dx, Y: g.player.Y + dy}

			// Valid move = inside the grid and not a scrap heap (wall).
			if target.X < 0 || target.X >= gridWidth || target.Y < 0 || target.Y >= gridHeight {
				continue
			}
			if g.scrapGrid[target.X][target.Y] {
				continue
			}

			ncx := pcx + float64(dx)*float64(cellSize)
			ncy := pcy + float64(dy)*float64(cellSize)
			drawArrow(screen, ncx, ncy, dx, dy, arrowColor)
		}
	}
}

// drawArrow renders a bold arrow centred at (cx,cy) pointing along (dx,dy),
// matching the original Daleks direction markers: a straight shaft with a clear
// arrowhead. Drawn with anti-aliased strokes so diagonals stay clean.
func drawArrow(screen *ebiten.Image, cx, cy float64, dx, dy int, clr color.Color) {
	const strokeWidth = 2.0

	// Unit direction (diagonals normalised so all arrows share one length).
	length := 1.0
	if dx != 0 && dy != 0 {
		length = math.Sqrt2
	}
	ux := float64(dx) / length
	uy := float64(dy) / length
	// Perpendicular, for the arrowhead barbs.
	px := -uy
	py := ux

	shaft := float64(cellSize) * 0.36 // half the shaft length
	barb := float64(cellSize) * 0.32  // arrowhead length
	spread := 0.85                    // how wide the arrowhead opens

	tipX := cx + ux*shaft
	tipY := cy + uy*shaft
	backX := cx - ux*shaft
	backY := cy - uy*shaft

	b1x := tipX - ux*barb + px*barb*spread
	b1y := tipY - uy*barb + py*barb*spread
	b2x := tipX - ux*barb - px*barb*spread
	b2y := tipY - uy*barb - py*barb*spread

	vector.StrokeLine(screen, float32(backX), float32(backY), float32(tipX), float32(tipY), strokeWidth, clr, true)
	vector.StrokeLine(screen, float32(tipX), float32(tipY), float32(b1x), float32(b1y), strokeWidth, clr, true)
	vector.StrokeLine(screen, float32(tipX), float32(tipY), float32(b2x), float32(b2y), strokeWidth, clr, true)
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

	// Don't highlight the player's own cell — no overlay for the wait-in-place move.
	if targetPos == g.player {
		return
	}

	// Only show the indicator for adjacent cells (valid one-step moves).
	dx := abs(targetPos.X - g.player.X)
	dy := abs(targetPos.Y - g.player.Y)
	if dx > 1 || dy > 1 {
		return
	}

	x := float64(gridOffsetX + gridX*cellSize)
	y := float64(gridOffsetY + gridY*cellSize)

	// Red if blocked by scrap, blue for a valid move.
	g.ensureScrapGrid()
	var indicatorColor color.Color
	if g.scrapGrid[targetPos.X][targetPos.Y] {
		indicatorColor = color.RGBA{255, 0, 0, 100} // Red for blocked
	} else {
		indicatorColor = color.RGBA{0, 0, 255, 100} // Blue for valid move
	}

	// Draw semi-transparent overlay on the cell
	ebitenutil.DrawRect(screen, x, y, cellSize, cellSize, indicatorColor)

	// Draw border
	ebitenutil.DrawRect(screen, x, y, cellSize, 1, color.Black)
	ebitenutil.DrawRect(screen, x, y, 1, cellSize, color.Black)
	ebitenutil.DrawRect(screen, x+cellSize-1, y, 1, cellSize, color.Black)
	ebitenutil.DrawRect(screen, x, y+cellSize-1, cellSize, 1, color.Black)
}
