package daleks

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var gameImages *DalekGameImages

func init() {
	gameImages = loadImages()
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		state:         StateMenu,
		level:         1,
		teleports:     10,
		safeTeleports: 3,
		screwdrivers:  2,
		lastStands:    1,
		lastMoveTime:  time.Now(),
		playerImage:   gameImages.Human,
		dalekImage:    gameImages.Dalek,

		scrapImage:            createScrapImage(),
		moveAnimationDuration: 0.6, // Duration for normal movement
		daleksMoving:          false,
		showGrid:              false, // Default OFF
		// Last Stand smooth movement settings
		lastStandSpeed:        2.0,  // Start speed in cells per second
		lastStandAcceleration: 1.5,  // Speed multiplier per second
		lastStandMaxSpeed:     20.0, // Maximum speed cap
		lastClickTime:         time.Now(),
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
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyK) {
				g.movePlayer(0, -1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
				g.movePlayer(0, 1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyH) {
				g.movePlayer(-1, 0)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
				g.movePlayer(1, 0)
			}

			// Diagonal movement
			if inpututil.IsKeyJustPressed(ebiten.KeyY) {
				g.movePlayer(-1, -1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyU) {
				g.movePlayer(1, -1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyB) {
				g.movePlayer(-1, 1)
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyN) {
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
