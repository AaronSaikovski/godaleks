package daleks

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Board struct {
	//rows     int
	//cols     int
	player    *Player
	robots    []*Player
	points    int
	gameOver  bool
	lastStand bool
	timer     time.Time
}

func NewBoard() *Board {
	//rand.Seed(time.Now().UnixNano())
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))

	board := &Board{
		timer: time.Now(),
	}
	// Place our hero
	board.placeHero()

	//Place the robots
	board.placeRobots()

	return board
}

// placeHero - Puts our hero randomly on the board
func (b *Board) placeHero() {
	HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new hero and start the hero in a random starting position
	HeroPlayer = Player{HeroImage, 0, 0, PlayerSpeed, true, true, HumanPlayer, false, false}
	xHeroStart, yHeroStart := RandomPlayerStartPosition(&HeroPlayer)
	HeroPlayer.xPos = xHeroStart
	HeroPlayer.yPos = yHeroStart
	HeroPlayer.isTeleporting = false
	b.player = &HeroPlayer

}

// placeRobots - Place the robots on the board
func (b *Board) placeRobots() {
	//Setup the Robots slice and add image and add random position
	for i := 0; i < StartRobots; i++ {
		//strRobotImg := "./assets/images/robot0" + strconv.Itoa(i+1) + ".png"
		strRobotImg := "./assets/images/robot.png"
		RobotImage, _, err = ebitenutil.NewImageFromFile(strRobotImg)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new robot struct and set start pos
		newRobot := &Player{RobotImage, 0, 0, PlayerSpeed, true, true, RobotPlayer, false, false}
		xRobotStart, yRoboStart := RandomPlayerStartPosition(newRobot)
		newRobot.xPos = xRobotStart
		newRobot.yPos = yRoboStart
		Robots = append(Robots, newRobot)
		b.robots = Robots
	}
}

// Reset - Resets the game
func (b *Board) Reset() {

	// Clear the screen with a white color again after the reset
	//ebiten.SetScreenTransparent(false)

	//taken fom https://github.com/hajimehoshi/ebiten/commit/8e5ae8873878a32e27e0c87fb6b3fb9c7e0d4c0a
	// and https://github.com/hajimehoshi/ebiten/issues/2378
	// op := &ebiten.RunGameOptions{}
	// op.ScreenTransparent = false
	// if err := ebiten.RunGameWithOptions(g{}, op); err != nil {
	// 	log.Fatal(err)
	// }

	//b.gameOver = false
	//b.StartNewGame()
}

// StartNewGame - Starts a new game and resets everything
func (b *Board) StartNewGame() {

	// Reset everything
	b.gameOver = false
	b.points = 0

	// Clear the screen with a white color again after the reset
	//ebiten.SetScreenTransparent(false)

	// Place our hero
	//b.placeHero()

	//Place the robots
	//b.placeRobots()

	// GameScore = 0
	// GameEnded = false
	// HeroImage = nil
	// HeroPlayer = Player{}
	// RobotImage = nil
	// Robots = nil
	// HighScore = 0

	// var err error
	// HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Create a new hero and start the hero in a random starting position
	// HeroPlayer = Player{HeroImage, 0, 0, PlayerSpeed, true, true, HumanPlayer}
	// xHeroStart, yHeroStart := RandomPlayerStartPosition(&HeroPlayer)
	// HeroPlayer.xPos = xHeroStart
	// HeroPlayer.yPos = yHeroStart

	// //Setup the Robots slice and add image and add random position
	// for i := 0; i < StartRobots; i++ {
	// 	//strRobotImg := "./assets/images/robot0" + strconv.Itoa(i+1) + ".png"
	// 	strRobotImg := "./assets/images/robot.png"
	// 	RobotImage, _, err = ebitenutil.NewImageFromFile(strRobotImg)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// Create a new robot struct and set start pos
	// 	newRobot := &Player{RobotImage, 0, 0, PlayerSpeed, true, true, RobotPlayer}
	// 	xRobotStart, yRoboStart := RandomPlayerStartPosition(newRobot)
	// 	newRobot.xPos = xRobotStart
	// 	newRobot.yPos = yRoboStart
	// 	Robots = append(Robots, newRobot)
	// }

}

func (b *Board) Update(input *Input) error {
	if b.gameOver {
		return nil
	}

	// Ensure the Hero doesnt go off the game
	CheckPlayerBoundary(&HeroPlayer)

	//Teleport our here
	if HeroPlayer.isTeleporting {
		HeroPlayer.image.Clear()

		//delay redraw
		time.Sleep(2 * time.Second)

		//redraw the hero
		HeroPlayer.isTeleporting = false
		b.placeHero()
	}

	// Check if Robots are colliding
	CheckRobotsCollision(Robots)

	//Check if robots are all alive
	if !CheckAllRobotsAlive(Robots) {
		b.gameOver = true
	} else {
		b.gameOver = false
	}

	// Move the player
	HeroPlayer.Move()

	// Check if Robots are colliding
	CheckHeroCollision(&HeroPlayer, Robots)

	// Is Hero alive?
	if !HeroPlayer.isAlive {
		b.gameOver = true
	} else {
		b.gameOver = false
	}

	// New game
	// if ebiten.IsKeyPressed(ebiten.KeyN) {
	// 	//b.Reset()
	// 	//StartNewGame()
	// }

	return nil
}

// RandomPlayerStartPosition - Places our player(s) in a random coordinate on the grid
func RandomPlayerStartPosition(Player *Player) (xPos, yPos float64) {
	// Retrieve the window size
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate the maximum X and Y coordinates for the sprite to stay within the window
	maxX := float64(windowWidth - Player.GetPlayerImageWidth())
	maxY := float64(windowHeight - Player.GetPlayerImageHeight())

	//seed the randomiser
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Return random X & Y coords
	return rand.Float64() * maxX, rand.Float64() * maxY
}
