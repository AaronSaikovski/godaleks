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
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) movePlayer(dx, dy int) {
	if g.state != StatePlaying || g.daleksMoving || g.isLastStandActive {
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

	if g.state != StatePlaying || g.daleksMoving || g.isLastStandActive {
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

	if g.state != StatePlaying || g.screwdrivers <= 0 || g.daleksMoving || g.isLastStandActive {
		return
	}

	g.screwdrivers--

	// Count normal daleks to determine if emperor is vulnerable
	normalDalekCount := g.countNormalDaleks()
	emperorIsVulnerable := normalDalekCount == 0

	// Find all daleks adjacent to player (including diagonally)
	daleksToDestroy := make([]int, 0)
	g.screwdriverTargets = make([]Position, 0)

	for i, dalek := range g.daleks {
		// Skip emperor if it's still invulnerable
		if dalek.IsEmperor && !emperorIsVulnerable {
			continue
		}

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

	// Pre-calculate alpha
	alpha := uint8(255 * (1.0 - progress))

	// Ultra-optimized sparkle effect - minimal particles
	numParticles := 6 // Reduced from 8
	radius := float64(cellSize) * (1.0 + progress*2.0)
	angleStep := 2.0 * 3.14159 / float64(numParticles)

	for i := 0; i < numParticles; i++ {
		angle := float64(i)*angleStep + progress*6.28
		cosA := math.Cos(angle)
		sinA := math.Sin(angle)

		px := int(x + radius*0.5*cosA)
		py := int(y + radius*0.5*sinA)

		// Draw single pixel sparkle for maximum performance
		if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
			screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, alpha})
		}
	}

	// Optimized flash - draw ring only with fewer points
	flashAlpha := uint8(255 * (1.0 - progress) * 0.8)
	flashRadius := int(float64(cellSize/2) * (1.0 - progress*0.5))

	// Draw circular outline with larger angle steps (fewer points)
	for angle := 0.0; angle < 6.28; angle += 0.4 {
		px := int(x + float64(flashRadius)*math.Cos(angle))
		py := int(y + float64(flashRadius)*math.Sin(angle))
		if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
			screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, flashAlpha})
		}
	}
}

func (g *Game) drawScrewdriverEffect(screen *ebiten.Image, pos Position, progress float64, offsetX, offsetY int) {
	x := float64(offsetX + pos.X*cellSize + cellSize/2)
	y := float64(offsetY + pos.Y*cellSize + cellSize/2)

	// Sonar-like expanding ring effect
	maxRadius := float64(cellSize) * 1.5

	// Draw multiple concentric rings for sonar effect
	numRings := 3
	for ring := 0; ring < numRings; ring++ {
		// Stagger the rings so they appear at different times
		ringProgress := progress - float64(ring)*0.15
		if ringProgress < 0 {
			continue
		}
		if ringProgress > 1.0 {
			ringProgress = 1.0
		}

		// Expanding ring radius
		ringRadius := int(maxRadius * ringProgress)

		// Fade out as ring expands
		ringAlpha := uint8(255 * (1.0 - ringProgress) * 0.8)

		// Draw the ring with moderate density for smooth sonar appearance
		for angle := 0.0; angle < 6.28; angle += 0.3 {
			px := int(x + float64(ringRadius)*math.Cos(angle))
			py := int(y + float64(ringRadius)*math.Sin(angle))
			if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
				screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, ringAlpha})
			}
		}
	}

	// Central pulse that fades quickly
	if progress < 0.3 {
		centralAlpha := uint8(255 * (1.0 - progress/0.3))
		centralRadius := int(float64(cellSize/4) * (1.0 + progress*2.0))

		for angle := 0.0; angle < 6.28; angle += 0.2 {
			px := int(x + float64(centralRadius)*math.Cos(angle))
			py := int(y + float64(centralRadius)*math.Sin(angle))
			if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
				screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, centralAlpha})
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

// Draw collision explosion effect (optimized)
func (g *Game) drawCollisionEffect(screen *ebiten.Image, pos FloatPosition, progress float64, offsetX, offsetY int) {
	x := float64(offsetX) + pos.X*float64(cellSize) + float64(cellSize)/2
	y := float64(offsetY) + pos.Y*float64(cellSize) + float64(cellSize)/2

	// Early exit if effect is off-screen
	maxRadius := float64(cellSize) * 1.5
	if x+maxRadius < 0 || x-maxRadius > float64(screenWidth) ||
		y+maxRadius < 0 || y-maxRadius > float64(screenHeight) {
		return
	}

	// Pre-calculate alpha values (avoid repeated calculations)
	alpha := uint8(255 * (1.0 - progress))

	// Ultra-optimized particle rendering - minimal particles
	numParticles := 8 // Reduced from 12 for better performance
	radius := maxRadius * progress
	angleStep := 2.0 * 3.14159 / float64(numParticles)

	// Draw particles with single pixel operations
	for i := 0; i < numParticles; i++ {
		angle := float64(i) * angleStep
		cosA := math.Cos(angle)
		sinA := math.Sin(angle)

		px := int(x + radius*cosA)
		py := int(y + radius*sinA)

		// Single pixel per particle for maximum performance
		if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
			screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, alpha})
		}
	}

	// Optimized flash with fewer points
	if progress < 0.5 {
		flashAlpha := uint8(255 * (1.0 - progress*2))
		flashRadius := int(float64(cellSize/3) * (1.0 + progress))

		// Draw only edge pixels with larger angle steps
		for angle := 0.0; angle < 6.28; angle += 0.5 {
			px := int(x + float64(flashRadius)*math.Cos(angle))
			py := int(y + float64(flashRadius)*math.Sin(angle))
			if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
				screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, flashAlpha})
			}
		}
	}

	// Optimized shock wave with minimal points
	if progress > 0.2 {
		waveProgress := (progress - 0.2) / 0.8
		waveRadius := int(maxRadius * waveProgress)
		waveAlpha := uint8(255 * (1.0 - waveProgress))

		// Draw ring with large angle steps for maximum performance
		for angle := 0.0; angle < 6.28; angle += 0.4 {
			px := int(x + float64(waveRadius)*math.Cos(angle))
			py := int(y + float64(waveRadius)*math.Sin(angle))
			if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
				screen.Set(px, py, color.RGBA{0x00, 0x00, 0x00, waveAlpha})
			}
		}
	}
}

// Update collision animations
func (g *Game) updateCollisionEffects(deltaTime float64) {
	activeEffects := make([]CollisionEffect, 0, len(g.collisionEffects))

	for i := range g.collisionEffects {
		effect := &g.collisionEffects[i]
		effect.Timer += deltaTime

		// Keep effect if not finished
		if effect.Timer < effect.Duration {
			activeEffects = append(activeEffects, *effect)
		}
	}

	g.collisionEffects = activeEffects
}
