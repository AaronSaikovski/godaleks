package daleks

import (
	"math/rand"
	"time"

	"github.com/AaronSaikovski/gorobots/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type Board struct {
	//rows     int
	//cols     int
	theDoctor *Player
	robots    []*Player
	points    int
	gameOver  bool
	//lastStand bool
	timer time.Time
}

func NewBoard() (*Board, error) {
	//rand.Seed(time.Now().UnixNano())
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))

	board := &Board{
		timer: time.Now(),
	}
	// Place our hero
	board.PositionHero()

	//Place the robots
	board.PositionRobots()

	return board, nil
}

// PositionHero - Puts our hero randomly on the board
func (b *Board) PositionHero() {

	// if theDoctor.isTeleporting {
	// 	HeroImage = nil
	// 	HeroPlayer = Player{}
	// 	b.player = nil
	// }

	// HeroImage, _, err = ebitenutil.NewImageFromFile("./assets/images/hero.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create a new hero and start the hero in a random starting position
	b.theDoctor = NewPlayer(assets.HeroImage, 0, 0, PlayerSpeed, true, true, HumanPlayer, false, false)
	xHeroStart, yHeroStart := randomPlayerStartPosition(b.theDoctor)
	b.theDoctor.xPos = xHeroStart
	b.theDoctor.yPos = yHeroStart
	b.theDoctor.isTeleporting = false

}

// PositionRobots - Place the robots on the board
func (b *Board) PositionRobots() {
	//Setup the Robots slice and add image and add random position
	for i := 0; i < StartRobots; i++ {
		//strRobotImg := "./assets/images/robot0" + strconv.Itoa(i+1) + ".png"
		//strRobotImg := "./assets/images/robot.png"
		//RobotImage, _, err = ebitenutil.NewImageFromFile(strRobotImg)
		// RobotImage, _, err = ebitenutil.NewImageFromFile(assets.Dalek)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Create a new robot struct and set start pos
		newRobot := NewPlayer(assets.Dalek, 0, 0, PlayerSpeed, true, true, RobotPlayer, false, false)
		xRobotStart, yRoboStart := randomPlayerStartPosition(newRobot)
		newRobot.xPos = xRobotStart
		newRobot.yPos = yRoboStart
		Robots = append(Robots, newRobot)
		b.robots = Robots
	}
}

// MAIN Update method
// Update  - updates the board state.
func (b *Board) Update(input *Input) error {

	if b.gameOver {
		return nil
	}

	// Ensure the Hero doesnt go off the game
	CheckPlayerBoundary(b.theDoctor)

	// Check if Robots are colliding
	CheckRobotsCollision(b.robots)

	//Check if robots are all alive
	if !CheckAllRobotsAlive(b.robots) {
		b.gameOver = true
	} else {
		b.gameOver = false
	}

	// Respond to moving the doctor - will move the robots too
	b.theDoctor.Move()

	// Check if Robots are colliding with player
	CheckHeroCollision(b.theDoctor, b.robots)

	// Is Hero alive?
	if !b.theDoctor.isAlive {
		b.gameOver = true
	} else {
		b.gameOver = false
	}

	return nil
}

// RandomPlayerStartPosition - Places our player(s) in a random coordinate on the grid
func randomPlayerStartPosition(Player *Player) (xPos, yPos float64) {
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
