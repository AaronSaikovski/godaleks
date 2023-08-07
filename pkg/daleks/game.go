package daleks

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// PlayerType - Custom Player datatype
type PlayerType string

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	PlayerSpeed  = 0.5
	StartRobots  = 5

	// Define our player types
	HumanPlayer PlayerType = "Human"
	RobotPlayer PlayerType = "Robot"
)

// Create our empty vars
var (
	//backgroundColor = color.RGBA{0, 0, 0, 0}
	//backgroundColor = color.RGBA{255, 255, 255, 0}
	backgroundColor = color.RGBA{50, 100, 50, 50}
	err             error
	HeroImage       *ebiten.Image
	RobotImage      *ebiten.Image
	HeroPlayer      Player
	Robots          []*Player
)

// Game - Game struct
type Game struct {
	input *Input
	board *Board
}

// NewGame - starts a new game
func NewGame() *Game {
	return &Game{
		input: NewInput(),
		board: NewBoard(),
	}
}

// Layout - define the size of the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Draw - Render the screen
func (g *Game) Draw(screen *ebiten.Image) {

	// set background
	screen.Fill(backgroundColor)

	if g.board.gameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Game Over. Score: %d", g.board.points))
	} else {
		// width := ScreenHeight / boardRows

		// Draw our hero
		g.board.player.DrawHero(screen)

		//Draw the robots
		g.board.player.DrawRobots(screen, Robots)

		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.board.points))
	}

}

// Update - update the logical state
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	return g.board.Update(g.input)
}
