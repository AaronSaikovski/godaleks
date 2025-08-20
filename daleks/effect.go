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

package daleks

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) movePlayer(dx, dy int) {
	if g.state != StatePlaying || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	// Prevent too rapid movement
	if time.Since(g.lastMoveTime) < 100*time.Millisecond {
		return
	}
	g.lastMoveTime = time.Now()

	newPos := Position{
		X: g.player.X + dx,
		Y: g.player.Y + dy,
	}

	// Check bounds
	if newPos.X < 0 || newPos.X >= gridWidth || newPos.Y < 0 || newPos.Y >= gridHeight {
		return
	}

	// Check if position is occupied by scrap
	for _, scrap := range g.scraps {
		if scrap == newPos {
			return
		}
	}

	g.player = newPos

	// In Last Stand mode, daleks move continuously, so no need to call moveDaleks
	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

func (g *Game) teleport(safe bool) {

	// ...existing code...
	g.soundPlayer.Play("teleport")

	if g.state != StatePlaying || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	if safe && g.safeTeleports <= 0 {
		return
	}
	if !safe && g.teleports <= 0 {
		return
	}

	// Store old position for animation
	g.teleportOldPos = g.player

	var newPos Position
	maxAttempts := 100

	if safe {
		// Safe teleport - find position with no daleks nearby
		for i := 0; i < maxAttempts; i++ {
			newPos = Position{
				X: rand.Intn(gridWidth),
				Y: rand.Intn(gridHeight),
			}

			if !g.positionOccupied(newPos) && g.isSafePosition(newPos) {
				break
			}
		}
		g.safeTeleports--
	} else {
		// Regular teleport - just find an empty spot
		for i := 0; i < maxAttempts; i++ {
			newPos = Position{
				X: rand.Intn(gridWidth),
				Y: rand.Intn(gridHeight),
			}

			if !g.positionOccupied(newPos) {
				break
			}
		}
		g.teleports--
	}

	// Start teleportation animation
	g.teleportNewPos = newPos
	g.teleportAnimation = true
	g.teleportTimer = 0
	g.player = newPos

	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

// Updated version to not use debris field when using screwdriver
func (g *Game) useScrewdriver() {

	if g.state != StatePlaying || g.screwdrivers <= 0 || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	g.screwdrivers--

	// Find all daleks adjacent to player (including diagonally)
	daleksToDestroy := make([]int, 0)
	g.screwdriverTargets = make([]Position, 0)

	for i, dalek := range g.daleks {
		dx := abs(dalek.GridPos.X - g.player.X)
		dy := abs(dalek.GridPos.Y - g.player.Y)

		// Adjacent includes all 8 surrounding cells
		if dx <= 1 && dy <= 1 && (dx != 0 || dy != 0) {
			daleksToDestroy = append(daleksToDestroy, i)
			g.screwdriverTargets = append(g.screwdriverTargets, dalek.GridPos)
		}
	}

	// Start screwdriver animation if there are targets
	if len(g.screwdriverTargets) > 0 {
		g.screwdriverAnimation = true
		g.soundPlayer.Play("screwdriver")
		g.screwdriverTimer = 0
	}

	// Remove destroyed daleks but DON'T add scraps
	newDaleks := make([]Dalek, 0, len(g.daleks))
	for i, dalek := range g.daleks {
		destroyed := false
		for _, destroyIndex := range daleksToDestroy {
			if i == destroyIndex {
				destroyed = true
				g.score += 5 // Bonus points for screwdriver kill
				// REMOVED: Don't add debris pile at dalek's position
				// g.scraps = append(g.scraps, dalek.GridPos)
				break
			}
		}
		if !destroyed {
			newDaleks = append(newDaleks, dalek)
		}
	}

	g.daleks = newDaleks

	// Move remaining daleks after screwdriver use (if not in Last Stand)
	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

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

func (g *Game) updateDalekAnimations(deltaTime float64) {
	// Ensure Last Stand is properly disabled if game is not in playing state
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	if g.isLastStandActive {
		// Smooth continuous movement during Last Stand
		g.updateLastStandMovement(deltaTime)
	} else {
		// Normal step-by-step movement
		g.updateNormalMovement(deltaTime)
	}
}
