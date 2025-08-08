package daleks

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func (g *Game) drawTeleportEffect(screen *ebiten.Image, pos Position, progress float64, offsetX, offsetY int) {
	x := float64(offsetX + pos.X*cellSize + cellSize/2)
	y := float64(offsetY + pos.Y*cellSize + cellSize/2)

	// Create sparkle/energy effect with black particles for classic Mac style
	numParticles := 8
	radius := float64(cellSize) * (1.0 + progress*2.0) // Expanding circle

	for i := 0; i < numParticles; i++ {
		angle := float64(i)*2.0*3.14159/float64(numParticles) + progress*6.28 // Rotating
		px := x + radius*0.5*math.Cos(angle)
		py := y + radius*0.5*math.Sin(angle)

		// Fade out over time
		alpha := uint8(255 * (1.0 - progress))
		sparkleColor := color.RGBA{0x00, 0x00, 0x00, alpha} // Black sparkles

		// Draw sparkle particles
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if int(px)+dx >= 0 && int(px)+dx < screenWidth &&
					int(py)+dy >= 0 && int(py)+dy < screenHeight {
					screen.Set(int(px)+dx, int(py)+dy, sparkleColor)
				}
			}
		}
	}

	// Central flash effect with black
	flashAlpha := uint8(255 * (1.0 - progress) * 0.8)
	flashColor := color.RGBA{0x00, 0x00, 0x00, flashAlpha} // Black flash

	flashRadius := int(float64(cellSize/2) * (1.0 - progress*0.5))
	for dx := -flashRadius; dx <= flashRadius; dx++ {
		for dy := -flashRadius; dy <= flashRadius; dy++ {
			if dx*dx+dy*dy <= flashRadius*flashRadius {
				px := int(x) + dx
				py := int(y) + dy
				if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
					screen.Set(px, py, flashColor)
				}
			}
		}
	}
}

func (g *Game) drawScrewdriverEffect(screen *ebiten.Image, pos Position, progress float64, offsetX, offsetY int) {
	x := float64(offsetX + pos.X*cellSize + cellSize/2)
	y := float64(offsetY + pos.Y*cellSize + cellSize/2)

	// Create electric/energy effect for sonic screwdriver
	numBolts := 6
	radius := float64(cellSize) * 0.8

	for i := 0; i < numBolts; i++ {
		angle := float64(i)*2.0*3.14159/float64(numBolts) + progress*12.56 // Fast rotation

		// Create zigzag lightning bolt effect
		for step := 0; step < 5; step++ {
			stepRadius := radius * float64(step) / 5.0
			zigzag := math.Sin(progress*20.0+float64(step)) * 3.0 // Zigzag offset

			boltX := x + stepRadius*math.Cos(angle) + zigzag*math.Cos(angle+1.57)
			boltY := y + stepRadius*math.Sin(angle) + zigzag*math.Sin(angle+1.57)

			// Fade based on progress and distance
			alpha := uint8(255 * (1.0 - progress) * (1.0 - float64(step)/5.0))
			boltColor := color.RGBA{0x00, 0x00, 0x00, alpha} // Black lightning

			if int(boltX) >= 0 && int(boltX) < screenWidth &&
				int(boltY) >= 0 && int(boltY) < screenHeight {
				screen.Set(int(boltX), int(boltY), boltColor)
				// Make bolts thicker
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						px := int(boltX) + dx
						py := int(boltY) + dy
						if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
							screen.Set(px, py, boltColor)
						}
					}
				}
			}
		}
	}

	// Central pulse effect
	pulseRadius := int(float64(cellSize/3) * (1.0 + progress*2.0))
	pulseAlpha := uint8(255 * (1.0 - progress) * 0.6)
	pulseColor := color.RGBA{0x00, 0x00, 0x00, pulseAlpha}

	for dx := -pulseRadius; dx <= pulseRadius; dx++ {
		for dy := -pulseRadius; dy <= pulseRadius; dy++ {
			if dx*dx+dy*dy <= pulseRadius*pulseRadius {
				px := int(x) + dx
				py := int(y) + dy
				if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
					screen.Set(px, py, pulseColor)
				}
			}
		}
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
	// Status information
	status := fmt.Sprintf("Level: %d  Score: %d  Teleports: %d  Safe: %d  Screwdrivers: %d  Last Stands: %d  Daleks: %d",
		g.level, g.score, g.teleports, g.safeTeleports, g.screwdrivers, g.lastStands, len(g.daleks))
	text.Draw(screen, status, basicfont.Face7x13, 10, 20, color.Black)

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

func (g *Game) drawMenu(screen *ebiten.Image) {
	title := "DALEKS"
	text.Draw(screen, title, basicfont.Face7x13, screenWidth/2-len(title)*3, 100, color.Black)

	instructions := []string{
		"Use arrow keys or HJK to move",
		"Y, U, B, N for diagonal movement",
		"SPACE or . to wait",
		"T to teleport randomly",
		"R to teleport safely",
		"S to use sonic screwdriver",
		"L for Last Stand (all daleks rush you)",
		"G to turn game grid On/Off",
		"",
		"MOUSE: Click adjacent cell to move there",
		"Click on player to wait in place",
		"",
		"Avoid the Daleks!",
		"Make them crash into each other!",
		"Sonic Screwdriver destroys adjacent Daleks!",
		"Last Stand forces all daleks to move!",
		"",
		"Press SPACE or click to start",
	}

	for i, line := range instructions {
		text.Draw(screen, line, basicfont.Face7x13, 50, 200+i*20, color.Black)
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

	// Draw scraps (centered)
	for _, scrap := range g.scraps {
		x, y := getCenteredSpritePosition(scrap.X, scrap.Y, offsetX, offsetY, g.scrapImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.scrapImage, op)
	}

	// Draw daleks using smooth interpolated positions (centered)
	for _, dalek := range g.daleks {
		// Use visual position for smooth movement, but calculate centered position
		cellCenterX := float64(offsetX) + dalek.VisualPos.X*float64(cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + dalek.VisualPos.Y*float64(cellSize) + float64(cellSize)/2

		// Get sprite dimensions and center it
		spriteBounds := g.dalekImage.Bounds()
		spriteWidth := spriteBounds.Dx()
		spriteHeight := spriteBounds.Dy()

		x := cellCenterX - float64(spriteWidth)/2
		y := cellCenterY - float64(spriteHeight)/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.dalekImage, op)
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

		// Draw player with fade effect (centered)
		x, y := getCenteredSpritePosition(g.player.X, g.player.Y, offsetX, offsetY, g.playerImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)

		// Fade in player at new position
		if progress > 0.3 {
			fadeAlpha := (progress - 0.3) / 0.7
			op.ColorM.Scale(1, 1, 1, fadeAlpha)
		} else {
			op.ColorM.Scale(1, 1, 1, 0) // Invisible during first part
		}

		screen.DrawImage(g.playerImage, op)
	} else {
		// Normal player drawing (centered)
		x, y := getCenteredSpritePosition(g.player.X, g.player.Y, offsetX, offsetY, g.playerImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.playerImage, op)
	}

	// Draw screwdriver effects
	if g.screwdriverAnimation {
		progress := g.screwdriverTimer / 0.8 // 0.8 second animation

		for _, target := range g.screwdriverTargets {
			g.drawScrewdriverEffect(screen, target, progress, offsetX, offsetY)
		}
	}
}
