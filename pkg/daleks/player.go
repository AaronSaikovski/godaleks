package daleks

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Player - Player struct
type Player struct {
	image         *ebiten.Image
	xPos, yPos    float64
	speed         float64
	isAlive       bool
	active        bool
	PlayerType    PlayerType
	NewGame       bool
	isTeleporting bool
}

// GetPlayerImageWidth - Returns the Player Image width
func (P *Player) GetPlayerImageWidth() int {
	return P.image.Bounds().Dx()
}

// GetPlayerImageHeight - Returns the Player Image height
func (P *Player) GetPlayerImageHeight() int {
	return P.image.Bounds().Dy()
}

// Move - Moves the hero around the grid
func (Player *Player) Move() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		Player.isTeleporting = false
		Player.yPos -= Player.speed

		//Move the Robot
		MoveRobot(Player, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		Player.isTeleporting = false
		Player.yPos += Player.speed

		//Move the Robot
		MoveRobot(Player, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		Player.isTeleporting = false
		Player.xPos -= Player.speed

		//Move the Robot
		MoveRobot(Player, Robots)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		Player.isTeleporting = false
		Player.xPos += Player.speed

		//Move the Robot
		MoveRobot(Player, Robots)
	}
	// //Laststand
	// if ebiten.IsKeyPressed(ebiten.KeyL) {
	// 	LastStand(HeroPlayer)
	// 	HeroPlayer.active = false
	// 	//fmt.Print("L pressed")
	// }
	//Sonic screwdriver
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		Player.isTeleporting = false
		fmt.Print("S pressed")
	}
	//Teleport
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		//TeleportHero(HeroPlayer)
		//Player.Teleport()
		//Player.active = true
		Player.isTeleporting = true
	}
	// New game
	if ebiten.IsKeyPressed(ebiten.KeyN) {
		Player.isTeleporting = false
		Player.NewGame = true
		//fmt.Print("Newgame")
	}
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

// DrawHero - Draws our Hero
func (Player *Player) DrawHero(screen *ebiten.Image) {

	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(Player.xPos, Player.yPos)
	screen.DrawImage(Player.image, playerOp)
}

// Teleport - Teleports the hero to a random place on the game grid
// func (Player *Player) Teleport() {
// 	//Player.active = false
// 	Player.isTeleporting = true
// 	// xHeroTeleport, yHeroTeleport := RandomPlayerStartPosition(Player)
// 	// Player.xPos = xHeroTeleport
// 	// Player.yPos = yHeroTeleport

// 	//Player.image.Clear()
// }

// DrawRobots - Draws a robot player(s)
func (Player *Player) DrawRobots(screen *ebiten.Image, Robots []*Player) {

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
