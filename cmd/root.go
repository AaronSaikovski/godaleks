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
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

var gameImages *DalekGameImages

//var soundPlayer *SoundPlayer

func init() {
	gameImages = loadImages()
}

func NewGame() *Game {
	// Note: rand.Seed is deprecated in Go 1.20+
	// The global random source is automatically seeded now

	soundPlayer, err := NewSoundPlayer()
	if err != nil {
		// Log error but continue with nil soundPlayer (graceful degradation)
		fmt.Printf("Warning: Failed to initialize sound player: %v\n", err)
	}

	now := time.Now()
	g := &Game{
		state:         StateMenu,
		level:         1,
		teleports:     10,
		safeTeleports: 3,
		screwdrivers:  2,
		lastStands:    1,
		lastMoveTime:  now,

		//playerImage:           createPlayerImage(),
		//dalekImage:            createDalekImage(),

		playerImage:  gameImages.Human,
		dalekImage:   gameImages.Dalek,
		emperorImage: gameImages.DalekEmperor,

		scrapImage:             createScrapImage(),
		moveAnimationDuration:  0.7, // Smooth eased glide per Dalek step
		daleksMoving:           false,
		showGrid:               false, // Default OFF
		emperorWarningMessage:  "",
		emperorWarningTimer:    0,
		emperorWarningDuration: 3.0, // Show warning for 3 seconds
		// Last Stand smooth movement settings
		lastStandSpeed:        2.0, // Slower, more controlled initial speed
		lastStandAcceleration: 1.1, // Very gentle acceleration for ultra-smooth ramping
		lastStandMaxSpeed:     8.0, // Lower max speed for slower, smoother gameplay
		lastClickTime:         now,
		soundPlayer:           soundPlayer,

		// Pre-allocate reusable buffers
		dalekRemoveMap:     make(map[int]bool, 20),
		scrapMap:           make(map[Position]bool, 50),
		daleksByPosition:   make(map[Position][]int, 20),
		newScrapsBuf:       make([]Position, 0, 5),
		survivingDaleksBuf: make([]Dalek, 0, 30),
	}

	// Cache sprite dimensions for performance (avoid Bounds() calls in draw loop)
	if g.playerImage != nil {
		bounds := g.playerImage.Bounds()
		g.playerHalfWidth = float64(bounds.Dx()) / 2
		g.playerHalfHeight = float64(bounds.Dy()) / 2
	}
	if g.dalekImage != nil {
		bounds := g.dalekImage.Bounds()
		g.dalekHalfWidth = float64(bounds.Dx()) / 2
		g.dalekHalfHeight = float64(bounds.Dy()) / 2
	}
	if g.scrapImage != nil {
		bounds := g.scrapImage.Bounds()
		g.scrapHalfWidth = float64(bounds.Dx()) / 2
		g.scrapHalfHeight = float64(bounds.Dy()) / 2
	}

	// Pre-compute trig lookup table for particle effects
	for i := range trigTableSize {
		angle := float64(i) * 2.0 * math.Pi / float64(trigTableSize)
		g.sinTable[i] = math.Sin(angle)
		g.cosTable[i] = math.Cos(angle)
	}

	g.cachedGridStatus = "Grid: OFF"

	return g
}

func (g *Game) Update() error {
	// Ebiten calls Update at a fixed tick rate (SetTPS, default 60), decoupled
	// from the display's render rate. Deriving deltaTime from the wall clock is
	// therefore uneven: when Ebiten runs several catch-up ticks back-to-back in a
	// single rendered frame, each measures a near-zero elapsed time and then the
	// next frame measures a full frame's worth. That jitter is what makes the
	// interpolated Dalek movement look jerky. Because the tick rate is fixed, a
	// constant timestep advances the animation by an even amount every tick,
	// giving perfectly smooth motion (Ebiten still keeps game-time correct by
	// running extra catch-up ticks when the machine falls behind).
	tps := ebiten.TPS()
	if tps <= 0 {
		tps = 60 // Fall back to the default when TPS is set to SyncWithFPS.
	}
	deltaTime := 1.0 / float64(tps)

	// Handle mouse input for player movement
	if g.state == StatePlaying && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleMouseClick(x, y)
	}

	// Update Dalek animations (handles both normal and Last Stand movement)
	if g.daleksMoving || g.isLastStandActive {
		g.updateDalekAnimations(deltaTime)

		// Double-check that game hasn't ended during Dalek movement
		if g.state != StatePlaying {
			g.isLastStandActive = false
			g.daleksMoving = false
		}
	}

	// Update teleport animation
	if g.teleportAnimation {
		g.teleportTimer += deltaTime
		if g.teleportTimer >= 0.5 {
			g.teleportAnimation = false
			g.teleportTimer = 0
		}
	}

	// Update screwdriver animation
	if g.screwdriverAnimation {
		g.screwdriverTimer += deltaTime
		if g.screwdriverTimer >= 0.8 {
			g.screwdriverAnimation = false
			g.screwdriverTimer = 0
			g.screwdriverTargets = nil
		}
	}

	// Update collision effects
	g.updateCollisionEffects(deltaTime)

	// Update emperor warning message timer
	if g.emperorWarningMessage != "" {
		g.emperorWarningTimer += deltaTime
		if g.emperorWarningTimer >= g.emperorWarningDuration {
			g.emperorWarningMessage = ""
			g.emperorWarningTimer = 0
		}
	}

	switch g.state {
	case StateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.startLevel()
		}

	case StatePlaying:

		// New Game - press N to start a new game
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.resetGame()
			return nil
		}

		// Toggle grid
		if inpututil.IsKeyJustPressed(ebiten.KeyG) {
			g.showGrid = !g.showGrid

			if g.showGrid {
				g.gridToggleMessage = "Grid ON"
				g.cachedGridStatus = "Grid: ON"
			} else {
				g.gridToggleMessage = "Grid OFF"
				g.cachedGridStatus = "Grid: OFF"
			}
			g.gridToggleMessageTime = time.Now()
		}

		// Allow movement during Last Stand, but not during normal dalek movement
		if !g.daleksMoving || g.isLastStandActive {

			// Movement and actions

			// //UP
			// if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			// 	g.movePlayer(0, -1)
			// }
			// //Down
			// if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			// 	g.movePlayer(0, 1)
			// }
			// //Left
			// if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			// 	g.movePlayer(-1, 0)
			// }
			// //Right
			// if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			// 	g.movePlayer(1, 0)
			// }
			//UP
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
				g.movePlayer(0, -1)
			}
			//Down
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
				g.movePlayer(0, 1)
			}
			//Left
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
				g.movePlayer(-1, 0)
			}
			//Right
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				g.movePlayer(1, 0)
			}

			// Diagonal movement
			if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
				g.movePlayer(-1, -1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyE) {
				g.movePlayer(1, -1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
				g.movePlayer(-1, 1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyC) {
				g.movePlayer(1, 1)
			}

			// Stay in place
			if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				if !g.isLastStandActive {
					g.moveDaleks()
				}
			}

			// Teleport
			if inpututil.IsKeyJustPressed(ebiten.KeyT) {
				g.teleport(false)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyR) {
				g.teleport(true)
			}

			// Sonic screwdriver
			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				g.useScrewdriver()
			}

			// Last stand
			if inpututil.IsKeyJustPressed(ebiten.KeyL) {
				g.lastStand()
			}
		}

	case StateLevelComplete:
		g.levelCompleteTimer += deltaTime
		if g.levelCompleteTimer >= levelTransitionDelay {
			g.startLevel()
		}

	case StateGameOver, StateWin:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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
			// Reset Last Stand speed settings
			g.lastStandSpeed = 2.0
			// Clear cached strings
			g.cachedFinalScore = ""
			g.cachedLevelMsg = ""
			g.cachedLevelNextMsg = ""
			g.cachedGridStatus = "Grid: OFF"
			// Clear any remaining game state
			g.daleks = nil
			g.scraps = nil
			g.gameOverMessage = ""
			g.state = StateMenu
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// White background
	screen.Fill(color.White)

	switch g.state {
	case StateMenu:
		g.drawMenu(screen)
	case StatePlaying:
		g.drawGame(screen)
		g.drawHUD(screen)
		g.drawPlayerArrows(screen)
		g.drawMouseIndicator(screen)
	case StateLevelComplete:
		g.drawGame(screen)
		g.drawHUD(screen)
		g.drawLevelComplete(screen)
	case StateGameOver, StateWin:
		g.drawGame(screen)
		g.drawHUD(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
