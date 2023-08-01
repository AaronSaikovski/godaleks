package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	spriteWidth  = 32
	spriteHeight = 32
	playerSpeed  = 1.0
)

// Create our empty vars
var (
	img        *ebiten.Image
	err        error
	HeroImage  *ebiten.Image
	RobotImage *ebiten.Image
	HeroPlayer Hero
	Robots     Robot
)

// type Player interface {
// 	image      *ebiten.Image
// 	xPos, yPos float64
// 	speed      float64
// 	isAlive    bool

// }

// Create the player class
type Hero struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
	isAlive    bool
}

type Robot struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
	isAlive    bool
}

//var img *ebiten.Image

func init() {
	var err error
	HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	if err != nil {
		log.Fatal(err)
	}

	//HeroPlayer = Hero{HeroImage, screenWidth / 2.0, screenHeight / 2.0, playerSpeed, true}

	// Start the hero in a random starting position
	xHeroStart, yHeroStart := randomPlayerStartPosition()
	HeroPlayer = Hero{HeroImage, xHeroStart, yHeroStart, playerSpeed, true}

	RobotImage, _, err = ebitenutil.NewImageFromFile("./assets/images/robot.png")
	if err != nil {
		log.Fatal(err)
	}

	xRobotStart, yRoboStart := randomPlayerStartPosition()
	//Robots = Robot{RobotImage, screenWidth / 2.0, screenHeight / 2.0, playerSpeed, true}
	Robots = Robot{RobotImage, xRobotStart, yRoboStart, playerSpeed, true}

}

type Game struct{}

// CheckPlayerBounds - Ensures the players stay within the game grid
func CheckHeroBoundary(HeroPlayer *Hero) {
	// Check if sprite goes off the left or right edge
	if HeroPlayer.xPos < 0 {
		HeroPlayer.xPos = 0
	} else if HeroPlayer.xPos > screenWidth-spriteWidth {
		HeroPlayer.xPos = screenWidth - spriteWidth
	}

	// Check if sprite goes off the top or bottom edge
	if HeroPlayer.yPos < 0 {
		HeroPlayer.yPos = 0
	} else if HeroPlayer.yPos > screenHeight-spriteHeight {
		HeroPlayer.yPos = screenHeight - spriteHeight
	}

}

// MoveHero - Moves the hero around the grid
func MoveHero(HeroPlayer *Hero) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		HeroPlayer.yPos -= HeroPlayer.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		HeroPlayer.yPos += HeroPlayer.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		HeroPlayer.xPos -= HeroPlayer.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		HeroPlayer.xPos += HeroPlayer.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		fmt.Print("L pressed")
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		fmt.Print("S pressed")
	}
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		fmt.Print("T pressed")
	}
	if ebiten.IsKeyPressed(ebiten.KeyN) {
		fmt.Print("N pressed")
	}
}

// randomPlayerStartPosition - Places our player(s) in a random coordinate on the grid
func randomPlayerStartPosition() (xPos, yPos float64) {
	// Retrieve the window size
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate the maximum X and Y coordinates for the sprite to stay within the window
	maxX := float64(windowWidth - spriteWidth)
	maxY := float64(windowHeight - spriteHeight)

	// Generate random X and Y coordinates within the window bounds
	rand.Seed(time.Now().UnixNano())

	// Return random X & Y coords
	return rand.Float64() * maxX, rand.Float64() * maxY
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	//Ensure the Hero doesnt go off the game
	CheckHeroBoundary(&HeroPlayer)

	// Move the hero
	MoveHero(&HeroPlayer)

	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	HeroPlayer.yPos -= HeroPlayer.speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	HeroPlayer.yPos += HeroPlayer.speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyLeft) {
	// 	HeroPlayer.xPos -= HeroPlayer.speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyRight) {
	// 	HeroPlayer.xPos += HeroPlayer.speed
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyL) {
	// 	fmt.Print("L pressed")
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyS) {
	// 	fmt.Print("S pressed")
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyT) {
	// 	fmt.Print("T pressed")
	// }
	return nil
}

// DrawHero -  Draws our hero
func DrawHero(screen *ebiten.Image) {
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(HeroPlayer.xPos, HeroPlayer.yPos)
	screen.DrawImage(HeroPlayer.image, playerOp)
}

// DrawEvilRobot - Draws an evil robot
func DrawEvilRobot(screen *ebiten.Image) {
	evilRobotOp := &ebiten.DrawImageOptions{}
	evilRobotOp.GeoM.Translate(Robots.xPos, Robots.yPos)
	screen.DrawImage(Robots.image, evilRobotOp)
}

func (g *Game) Draw(screen *ebiten.Image) {

	//set background
	screen.Fill(color.RGBA{255, 255, 255, 0})

	// Draw our hero
	DrawHero(screen)

	// Draw some bad robots
	DrawEvilRobot(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoRobots")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
