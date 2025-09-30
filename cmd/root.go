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
	"math/rand"
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
	rand.Seed(time.Now().UnixNano())

	soundPlayer, err := NewSoundPlayer()
	if err != nil {
		// Handle error appropriately for your game
		panic(err)
	}

	g := &Game{
		state:         StateMenu,
		level:         1,
		teleports:     10,
		safeTeleports: 3,
		screwdrivers:  2,
		lastStands:    1,
		lastMoveTime:  time.Now(),

		//playerImage:           createPlayerImage(),
		//dalekImage:            createDalekImage(),

		playerImage: gameImages.Human,
		dalekImage:  gameImages.Dalek,

		scrapImage:            createScrapImage(),
		moveAnimationDuration: 0.8, // Changed from 0.6 to 0.8 for slower movement
		daleksMoving:          false,
		showGrid:              false, // Default OFF
		// Last Stand smooth movement settings
		lastStandSpeed:        1.5,  // Changed from 2.0 to 1.5 for slower initial speed
		lastStandAcceleration: 1.2,  // Changed from 1.5 to 1.2 for gentler acceleration
		lastStandMaxSpeed:     15.0, // Changed from 20.0 to 15.0 for lower max speed
		lastClickTime:         time.Now(),
		soundPlayer:           soundPlayer,
	}

	return g
}

func (g *Game) Update() error {
	deltaTime := 1.0 / 60.0 // Assuming 60 FPS

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
			} else {
				g.gridToggleMessage = "Grid OFF"
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

		// Debug info - add this temporarily to see Last Stand status
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			fmt.Printf("Last Stand Debug - Active: %v, Moving: %v, Speed: %.2f, Daleks: %d, Last Stands Available: %d\n",
				g.isLastStandActive, g.daleksMoving, g.lastStandSpeed, len(g.daleks), g.lastStands)
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
		g.drawMouseIndicator(screen)
	case StateGameOver, StateWin:
		g.drawGame(screen)
		g.drawHUD(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
