package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	playerSpeed  = 0.5
)

// Create our empty vars
var (
	err            error
	HeroImage      *ebiten.Image
	EvilRobotImage *ebiten.Image
	HeroPlayer     Hero
	EvilRobots     Robot
)

// Create the player class
type Hero struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
}

type Robot struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
}

var img *ebiten.Image

func init() {
	var err error
	HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	if err != nil {
		log.Fatal(err)
	}

	HeroPlayer = Hero{HeroImage, screenWidth / 2.0, screenHeight / 2.0, playerSpeed*2}

	EvilRobotImage, _, err = ebitenutil.NewImageFromFile("./assets/images/robot.png")
	if err != nil {
		log.Fatal(err)
	}

	EvilRobots = Robot{EvilRobotImage, screenWidth / 2.0, screenHeight / 2.0, playerSpeed}

}

type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
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
	evilRobotOp.GeoM.Translate(EvilRobots.xPos, EvilRobots.yPos)
	screen.DrawImage(EvilRobots.image, evilRobotOp)
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
