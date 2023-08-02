package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// PlayerType - Custom Player datatype
type PlayerType string

const (
	screenWidth  = 800
	screenHeight = 600
	spriteWidth  = 32
	spriteHeight = 32
	playerSpeed  = 0.5
	startRobots  = 5

	// Define our player types
	HumanPlayer PlayerType = "Human"
	RobotPlayer PlayerType = "Robot"
)

// Create our empty vars
var (
	img        *ebiten.Image
	err        error
	HeroImage  *ebiten.Image
	RobotImage *ebiten.Image
	HeroPlayer Player
	Robots     []*Player

	GameScore int
	GameEnded bool
)

// Game - Game struct
type Game struct {
	isGameOver bool
}

// Player - Player struct
type Player struct {
	image      *ebiten.Image
	xPos, yPos float64
	speed      float64
	isAlive    bool
	active     bool
	PlayerType PlayerType
}

// StartNewGame - Starts a new game and resets everything
func StartNewGame() {

	//Reset everything
	GameScore = 0
	GameEnded = false
	HeroImage = nil
	HeroPlayer = Player{}
	RobotImage = nil
	Robots = nil

	var err error
	HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	if err != nil {
		log.Fatal(err)
	}

	// Start the hero in a random starting position
	xHeroStart, yHeroStart := randomPlayerStartPosition()
	HeroPlayer = Player{HeroImage, xHeroStart, yHeroStart, playerSpeed, true, true, HumanPlayer}

	//Setup the Robots slice
	for i := 0; i < startRobots; i++ {
		strRobotImg := "./assets/images/robot0" + strconv.Itoa(i+1) + ".png"
		RobotImage, _, err = ebitenutil.NewImageFromFile(strRobotImg)
		if err != nil {
			log.Fatal(err)
		}

		xRobotStart, yRoboStart := randomPlayerStartPosition()
		newRobot := &Player{RobotImage, xRobotStart, yRoboStart, playerSpeed, true, true, RobotPlayer}
		Robots = append(Robots, newRobot)
	}

}

// Reset - Resets the game
func (g *Game) Reset() {

	// Clear the screen with a white color again after the reset
	ebiten.SetScreenTransparent(false)

	g.isGameOver = false
	StartNewGame()
}

// init - Start a new game
func init() {
	StartNewGame()
}

// CheckHeroBoundary - Ensures the players stay within the game grid
func CheckHeroBoundary(HeroPlayer *Player) {
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

// ArePlayersColliding - Are the two sprites colliding?
func ArePlayersColliding(Player1, Player2 *Player) bool {
	return Player1.xPos < Player2.xPos+float64(Player2.image.Bounds().Dx()) &&
		Player1.xPos+float64(Player1.image.Bounds().Dx()) > Player2.xPos &&
		Player1.yPos < Player2.yPos+float64(Player2.image.Bounds().Dy()) &&
		Player1.yPos+float64(Player1.image.Bounds().Dy()) > Player2.yPos
}

// TeleportHero - Teleports the hero to a random place on the game grid
func TeleportHero(HeroPlayer *Player) {
	HeroPlayer.active = false
	xHeroTeleport, yHeroTeleport := randomPlayerStartPosition()
	HeroPlayer.xPos = xHeroTeleport
	HeroPlayer.yPos = yHeroTeleport
}

// MoveHero - Moves the hero around the grid
func MoveHero(HeroPlayer *Player) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		HeroPlayer.yPos -= HeroPlayer.speed

		//Move the Robot
		MoveRobot(HeroPlayer, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		HeroPlayer.yPos += HeroPlayer.speed

		//Move the Robot
		MoveRobot(HeroPlayer, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		HeroPlayer.xPos -= HeroPlayer.speed

		//Move the Robot
		MoveRobot(HeroPlayer, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		HeroPlayer.xPos += HeroPlayer.speed

		//Move the Robot
		MoveRobot(HeroPlayer, Robots)
	}
	//Laststand
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		fmt.Print("L pressed")
	}
	//Sonic screwdriver
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		fmt.Print("S pressed")
	}
	//Teleport
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		TeleportHero(HeroPlayer)
		HeroPlayer.active = true
	}
	// New game
	// if ebiten.IsKeyPressed(ebiten.KeyN) {
	// 	fmt.Print("Newgame")
	// 	StartNewGame()
	// }
}

// MoveRobot - Moves the robot to chase the player
func MoveRobot(HeroPlayer *Player, RobotPlayer []*Player) {
	for index := range Robots {

		// Only move the robot if they are alive
		if RobotPlayer[index].isAlive {
			if RobotPlayer[index].xPos < HeroPlayer.xPos {
				RobotPlayer[index].xPos += RobotPlayer[index].speed
			} else {
				RobotPlayer[index].xPos -= RobotPlayer[index].speed
			}

			if RobotPlayer[index].yPos < HeroPlayer.yPos {
				RobotPlayer[index].yPos += RobotPlayer[index].speed
			} else {
				RobotPlayer[index].yPos -= RobotPlayer[index].speed
			}

		}

	}

}

// randomPlayerStartPosition - Places our player(s) in a random coordinate on the grid
func randomPlayerStartPosition() (xPos, yPos float64) {
	// Retrieve the window size
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate the maximum X and Y coordinates for the sprite to stay within the window
	maxX := float64(windowWidth - spriteWidth)
	maxY := float64(windowHeight - spriteHeight)

	//seed the randomiser
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Return random X & Y coords
	return rand.Float64() * maxX, rand.Float64() * maxY
}

func CheckHeroCollision(HeroPlayer *Player) bool {
	// Check for collisions among sprites
	for index := range Robots {
		if ArePlayersColliding(HeroPlayer, Robots[index]) {
			HeroPlayer.isAlive = false
			return true
		} else {
			HeroPlayer.isAlive = true
			return false
		}
	}
	return false
}

// CheckRobotsCollision - Check if the Robots are colliding with each other
func CheckRobotsCollision() {
	for i := 0; i < len(Robots); i++ {
		for j := i + 1; j < len(Robots); j++ {
			// e.g., perform actions like removing sprites, triggering events, etc.
			if ArePlayersColliding(Robots[i], Robots[j]) {
				// Handle collision between sprites[i] and sprites[j]
				Robots[i].isAlive = false
				Robots[j].isAlive = false
			}
		}
	}

}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {

	if !g.isGameOver {
		//Ensure the Hero doesnt go off the game
		CheckHeroBoundary(&HeroPlayer)

		// check if we have a collision between the player and a robot
		if CheckHeroCollision(&HeroPlayer) {
			fmt.Print("player collision")
			HeroPlayer.isAlive = false
			g.isGameOver = true
		} else {
			g.isGameOver = false
			HeroPlayer.isAlive = true
		}

		// Move the hero
		if HeroPlayer.isAlive {
			MoveHero(&HeroPlayer)
		}

	}

	// Check if Robots are colliding
	CheckRobotsCollision()

	// New game
	if ebiten.IsKeyPressed(ebiten.KeyN) {
		g.Reset()
		StartNewGame()
	}

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

// DrawRobot - Draws an robot player
func DrawRobot(screen *ebiten.Image) {

	//Setup the Robots slice
	for index := range Robots {

		robotOp := &ebiten.DrawImageOptions{}
		robotOp.GeoM.Translate(Robots[index].xPos, Robots[index].yPos)
		screen.DrawImage(Robots[index].image, robotOp)

		// Only draw the robot if they are alive
		// if Robots[index].isAlive {
		// 	robotOp := &ebiten.DrawImageOptions{}
		// 	robotOp.GeoM.Translate(Robots[index].xPos, Robots[index].yPos)
		// 	screen.DrawImage(Robots[index].image, robotOp)
		// }

	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	//set background
	screen.Fill(color.RGBA{255, 255, 255, 0})

	// Draw our hero
	DrawHero(screen)

	// Draw some bad robots
	DrawRobot(screen)

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
