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
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func (g *Game) drawGame(screen *ebiten.Image) {
	offsetX := (screenWidth - gridWidth*cellSize) / 2
	offsetY := 50

	// Draw grid only if enabled
	if g.showGrid {
		for x := 0; x <= gridWidth; x++ {
			ebitenutil.DrawLine(screen,
				float64(offsetX+x*cellSize), float64(offsetY),
				float64(offsetX+x*cellSize), float64(offsetY+gridHeight*cellSize),
				color.Black)
		}
		for y := 0; y <= gridHeight; y++ {
			ebitenutil.DrawLine(screen,
				float64(offsetX), float64(offsetY+y*cellSize),
				float64(offsetX+gridWidth*cellSize), float64(offsetY+y*cellSize),
				color.Black)
		}
	}

	// Draw scraps (centered) - use cached dimensions and reusable draw options
	for _, scrap := range g.scraps {
		cellCenterX := float64(offsetX) + float64(scrap.X*cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + float64(scrap.Y*cellSize) + float64(cellSize)/2

		x := cellCenterX - g.scrapHalfWidth
		y := cellCenterY - g.scrapHalfHeight

		g.drawOp.GeoM.Reset()
		g.drawOp.GeoM.Translate(x, y)
		screen.DrawImage(g.scrapImage, &g.drawOp)
	}

	// Draw daleks using smooth interpolated positions - use cached dimensions
	for _, dalek := range g.daleks {
		// Use visual position for smooth movement, but calculate centered position
		cellCenterX := float64(offsetX) + dalek.VisualPos.X*float64(cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + dalek.VisualPos.Y*float64(cellSize) + float64(cellSize)/2

		x := cellCenterX - g.dalekHalfWidth
		y := cellCenterY - g.dalekHalfHeight

		g.drawOp.GeoM.Reset()
		g.drawOp.GeoM.Translate(x, y)
		screen.DrawImage(g.dalekImage, &g.drawOp)
	}

	// Draw player with teleportation effects (centered)
	if g.teleportAnimation {
		progress := g.teleportTimer / 0.5 // 0.5 second animation

		// Draw disappearing effect at old position
		if progress < 0.5 {
			g.drawTeleportEffect(screen, g.teleportOldPos, progress*2, offsetX, offsetY)
		}

		// Draw appearing effect at new position
		if progress > 0.3 {
			appearProgress := (progress - 0.3) / 0.7
			g.drawTeleportEffect(screen, g.teleportNewPos, 1.0-appearProgress, offsetX, offsetY)
		}

		// Draw player with fade effect (centered) - use cached dimensions
		cellCenterX := float64(offsetX) + float64(g.player.X*cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + float64(g.player.Y*cellSize) + float64(cellSize)/2

		x := cellCenterX - g.playerHalfWidth
		y := cellCenterY - g.playerHalfHeight

		g.drawOp.GeoM.Reset()
		g.drawOp.GeoM.Translate(x, y)

		g.drawOp.ColorM.Reset()
		// Fade in player at new position
		if progress > 0.3 {
			fadeAlpha := (progress - 0.3) / 0.7
			g.drawOp.ColorM.Scale(1, 1, 1, fadeAlpha)
		} else {
			g.drawOp.ColorM.Scale(1, 1, 1, 0) // Invisible during first part
		}

		screen.DrawImage(g.playerImage, &g.drawOp)
		g.drawOp.ColorM.Reset() // Reset color for next draw
	} else {
		// Normal player drawing (centered) - use cached dimensions
		cellCenterX := float64(offsetX) + float64(g.player.X*cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + float64(g.player.Y*cellSize) + float64(cellSize)/2

		x := cellCenterX - g.playerHalfWidth
		y := cellCenterY - g.playerHalfHeight

		g.drawOp.GeoM.Reset()
		g.drawOp.GeoM.Translate(x, y)
		screen.DrawImage(g.playerImage, &g.drawOp)
	}

	// Draw screwdriver effects
	if g.screwdriverAnimation {
		progress := g.screwdriverTimer / 0.8 // 0.8 second animation

		for _, target := range g.screwdriverTargets {
			g.drawScrewdriverEffect(screen, target, progress, offsetX, offsetY)
		}
	}

	// Draw collision explosion effects
	for _, effect := range g.collisionEffects {
		progress := effect.Timer / effect.Duration
		g.drawCollisionEffect(screen, effect.Pos, progress, offsetX, offsetY)
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Semi-transparent overlay
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 128})

	// Game over message
	text.Draw(screen, g.gameOverMessage, basicfont.Face7x13,
		screenWidth/2-len(g.gameOverMessage)*3, screenHeight/2-20, color.White)

	finalScore := fmt.Sprintf("Final Score: %d", g.score)
	text.Draw(screen, finalScore, basicfont.Face7x13,
		screenWidth/2-len(finalScore)*3, screenHeight/2+10, color.White)

	restart := "Press SPACE or click to restart"
	text.Draw(screen, restart, basicfont.Face7x13,
		screenWidth/2-len(restart)*3, screenHeight/2+40, color.White)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	// Cache status text and only update every 100ms to reduce fmt.Sprintf calls
	now := time.Now()
	if g.hudStatusText == "" || now.Sub(g.lastHUDUpdate) > 100*time.Millisecond {
		g.hudStatusText = fmt.Sprintf("Level: %d  Score: %d  Teleports: %d  Safe: %d  Screwdrivers: %d  Last Stands: %d  Daleks: %d",
			g.level, g.score, g.teleports, g.safeTeleports, g.screwdrivers, g.lastStands, len(g.daleks))
		g.lastHUDUpdate = now
	}
	text.Draw(screen, g.hudStatusText, basicfont.Face7x13, 10, 20, color.Black)

	// Grid indicator
	gridStatus := "Grid: OFF"
	if g.showGrid {
		gridStatus = "Grid: ON"
	}
	text.Draw(screen, gridStatus, basicfont.Face7x13, 10, 40, color.Black)

	// Last Stand indicator
	if g.isLastStandActive {
		lastStandMsg := fmt.Sprintf("LAST STAND ACTIVE! Speed: %.1f", g.lastStandSpeed)
		text.Draw(screen, lastStandMsg, basicfont.Face7x13, 10, screenHeight-30, color.Black)
	}

	// Temporary center-screen notification for grid toggle
	if g.gridToggleMessage != "" && time.Since(g.gridToggleMessageTime) < 1500*time.Millisecond {
		msg := g.gridToggleMessage
		x := screenWidth/2 - len(msg)*3
		y := 60
		text.Draw(screen, msg, basicfont.Face7x13, x, y, color.Black)
	}
}
