package daleks

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 50
	gridHeight   = 35 // Reduced from 37 to 35 to ensure sprites stay in bounds
	cellSize     = 16
)

type Position struct {
	X, Y int
}

type FloatPosition struct {
	X, Y float64
}

type GameState int

var gameImages *DalekGameImages

//r soundPlayer *SoundPlayer

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
	StateWin
)

type Dalek struct {
	GridPos   Position      // Current grid position
	VisualPos FloatPosition // Interpolated visual position
	TargetPos FloatPosition // Target visual position
	IsMoving  bool          // Whether currently animating
	MoveTimer float64       // Animation timer
}

type Game struct {
	state           GameState
	player          Position
	daleks          []Dalek // Changed from []Position to []Dalek
	scraps          []Position
	level           int
	score           int
	teleports       int
	safeTeleports   int
	screwdrivers    int
	lastStands      int
	gameOverMessage string
	lastMoveTime    time.Time

	playerImage *ebiten.Image
	dalekImage  *ebiten.Image
	scrapImage  *ebiten.Image
	// Movement animation settings
	moveAnimationDuration float64 // Duration for Dalek movement animation
	daleksMoving          bool    // Whether daleks are currently moving
	// Teleportation animation
	teleportAnimation bool
	teleportTimer     float64
	teleportOldPos    Position
	teleportNewPos    Position
	// Sonic screwdriver animation
	screwdriverAnimation  bool
	screwdriverTimer      float64
	screwdriverTargets    []Position
	isLastStandActive     bool
	showGrid              bool
	gridToggleMessage     string
	gridToggleMessageTime time.Time
	// Last Stand smooth movement
	lastStandSpeed        float64 // Speed in cells per second during Last Stand
	lastStandAcceleration float64 // Acceleration multiplier per second
	lastStandMaxSpeed     float64 // Maximum speed cap
	// Mouse support
	lastClickTime time.Time
	soundPlayer   *SoundPlayer
}

func init() {
	gameImages = loadImages()

	// ...existing code...
	// soundPlayer, err := NewSoundPlayer()
	// if err != nil {
	// 	// Handle error appropriately for your game
	// 	panic(err)
	// }

}

// Loads images
func loadImages() *DalekGameImages {
	gameImages := &DalekGameImages{}
	gameImages.LoadImages()
	return gameImages
}

// Helper function to calculate centered sprite position
func getCenteredSpritePosition(gridX, gridY, offsetX, offsetY int, spriteImage *ebiten.Image) (float64, float64) {
	// Get sprite dimensions
	spriteBounds := spriteImage.Bounds()
	spriteWidth := spriteBounds.Dx()
	spriteHeight := spriteBounds.Dy()

	// Calculate center position within the grid cell
	cellCenterX := float64(offsetX + gridX*cellSize + cellSize/2)
	cellCenterY := float64(offsetY + gridY*cellSize + cellSize/2)

	// Subtract half sprite size to center it
	x := cellCenterX - float64(spriteWidth)/2
	y := cellCenterY - float64(spriteHeight)/2

	return x, y
}

// createScrapImage creates a Dalek debris pile sprite like in the classic game
func createScrapImage() *ebiten.Image {
	img := ebiten.NewImage(cellSize-2, cellSize-2)

	size := cellSize - 2
	centerX := size / 2
	centerY := size / 2

	// Create Dalek debris - looks like scattered Dalek parts and metal fragments
	// Main debris cluster in center
	debrisPositions := []Position{
		// Central cluster
		{centerX, centerY},
		{centerX - 1, centerY},
		{centerX + 1, centerY},
		{centerX, centerY - 1},
		{centerX, centerY + 1},
		{centerX - 1, centerY - 1},
		{centerX + 1, centerY + 1},

		// Scattered pieces around the main cluster
		{centerX - 2, centerY},
		{centerX + 2, centerY},
		{centerX, centerY - 2},
		{centerX, centerY + 2},
		{centerX - 1, centerY + 2},
		{centerX + 1, centerY - 2},

		// Outer scattered debris
		{centerX - 3, centerY + 1},
		{centerX + 3, centerY - 1},
		{centerX - 2, centerY - 2},
		{centerX + 2, centerY + 2},

		// Small fragments
		{centerX - 4, centerY},
		{centerX + 4, centerY},
		{centerX, centerY - 3},
		{centerX, centerY + 3},
		{centerX - 3, centerY - 1},
		{centerX + 3, centerY + 1},
	}

	// Draw the debris pieces
	for _, pos := range debrisPositions {
		if pos.X >= 0 && pos.X < size && pos.Y >= 0 && pos.Y < size {
			img.Set(pos.X, pos.Y, color.Black)
		}
	}

	// Add some additional random scattered bits to make it look more chaotic
	additionalDebris := []Position{
		{centerX - 1, centerY + 3},
		{centerX + 1, centerY - 3},
		{centerX - 4, centerY + 1},
		{centerX + 4, centerY - 1},
		{centerX - 2, centerY + 3},
		{centerX + 2, centerY - 3},
	}

	for _, pos := range additionalDebris {
		if pos.X >= 0 && pos.X < size && pos.Y >= 0 && pos.Y < size {
			img.Set(pos.X, pos.Y, color.Black)
		}
	}

	return img
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

// resetGame resets the game to initial state and starts from level 1
func (g *Game) resetGame() {
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
	g.lastStandSpeed = 2.0
	g.daleks = nil
	g.scraps = nil
	g.gameOverMessage = ""
	g.startLevel()
}

func (g *Game) startLevel() {
	// Clear the board and reset all states
	g.scraps = nil
	g.daleksMoving = false
	g.isLastStandActive = false
	g.lastStandSpeed = 2.0
	g.teleportAnimation = false
	g.teleportTimer = 0
	g.screwdriverAnimation = false
	g.screwdriverTimer = 0
	g.screwdriverTargets = nil
	g.lastStands = 1

	// Place player randomly
	g.player = Position{
		X: rand.Intn(gridWidth),
		Y: rand.Intn(gridHeight),
	}

	// Place daleks (5 + level number)
	dalekCount := 5 + g.level
	g.daleks = make([]Dalek, 0, dalekCount)

	for len(g.daleks) < dalekCount {
		pos := Position{
			X: rand.Intn(gridWidth),
			Y: rand.Intn(gridHeight),
		}

		// Don't place dalek on player or too close
		if g.distance(pos, g.player) > 3 && !g.positionOccupied(pos) {
			dalek := Dalek{
				GridPos:   pos,
				VisualPos: FloatPosition{X: float64(pos.X), Y: float64(pos.Y)},
				TargetPos: FloatPosition{X: float64(pos.X), Y: float64(pos.Y)},
				IsMoving:  false,
				MoveTimer: 0,
			}
			g.daleks = append(g.daleks, dalek)
		}
	}

	g.state = StatePlaying
	g.soundPlayer.Play("gamestart")
}

func (g *Game) distance(a, b Position) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return dx*dx + dy*dy // Using squared distance for efficiency
}

func (g *Game) positionOccupied(pos Position) bool {
	for _, dalek := range g.daleks {
		if dalek.GridPos == pos {
			return true
		}
	}
	for _, scrap := range g.scraps {
		if scrap == pos {
			return true
		}
	}
	return false
}

// Convert screen coordinates to grid coordinates
func (g *Game) screenToGrid(screenX, screenY int) (int, int, bool) {
	offsetX := (screenWidth - gridWidth*cellSize) / 2
	offsetY := 50

	gridX := (screenX - offsetX) / cellSize
	gridY := (screenY - offsetY) / cellSize

	// Check if within grid bounds
	if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
		return gridX, gridY, true
	}
	return 0, 0, false
}

// Handle mouse click for player movement
func (g *Game) handleMouseClick(x, y int) {
	// Prevent too rapid clicking
	if time.Since(g.lastClickTime) < 100*time.Millisecond {
		return
	}
	g.lastClickTime = time.Now()

	// Convert to grid coordinates
	gridX, gridY, valid := g.screenToGrid(x, y)
	if !valid {
		return
	}

	targetPos := Position{X: gridX, Y: gridY}

	// Check if clicking on current player position (stay in place)
	if targetPos == g.player {
		if !g.isLastStandActive {
			g.moveDaleks()
		}
		return
	}

	// Calculate movement direction
	dx := 0
	dy := 0

	if targetPos.X > g.player.X {
		dx = 1
	} else if targetPos.X < g.player.X {
		dx = -1
	}

	if targetPos.Y > g.player.Y {
		dy = 1
	} else if targetPos.Y < g.player.Y {
		dy = -1
	}

	// Only allow one-step moves (adjacent cells including diagonal)
	if abs(dx) <= 1 && abs(dy) <= 1 {
		g.movePlayer(dx, dy)
	}
}

func (g *Game) movePlayer(dx, dy int) {
	if g.state != StatePlaying || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	// Prevent too rapid movement
	if time.Since(g.lastMoveTime) < 100*time.Millisecond {
		return
	}
	g.lastMoveTime = time.Now()

	newPos := Position{
		X: g.player.X + dx,
		Y: g.player.Y + dy,
	}

	// Check bounds
	if newPos.X < 0 || newPos.X >= gridWidth || newPos.Y < 0 || newPos.Y >= gridHeight {
		return
	}

	// Check if position is occupied by scrap
	for _, scrap := range g.scraps {
		if scrap == newPos {
			return
		}
	}

	g.player = newPos

	// In Last Stand mode, daleks move continuously, so no need to call moveDaleks
	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

func (g *Game) teleport(safe bool) {

	// ...existing code...
	g.soundPlayer.Play("teleport")

	if g.state != StatePlaying || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	if safe && g.safeTeleports <= 0 {
		return
	}
	if !safe && g.teleports <= 0 {
		return
	}

	// Store old position for animation
	g.teleportOldPos = g.player

	var newPos Position
	maxAttempts := 100

	if safe {
		// Safe teleport - find position with no daleks nearby
		for i := 0; i < maxAttempts; i++ {
			newPos = Position{
				X: rand.Intn(gridWidth),
				Y: rand.Intn(gridHeight),
			}

			if !g.positionOccupied(newPos) && g.isSafePosition(newPos) {
				break
			}
		}
		g.safeTeleports--
	} else {
		// Regular teleport - just find an empty spot
		for i := 0; i < maxAttempts; i++ {
			newPos = Position{
				X: rand.Intn(gridWidth),
				Y: rand.Intn(gridHeight),
			}

			if !g.positionOccupied(newPos) {
				break
			}
		}
		g.teleports--
	}

	// Start teleportation animation
	g.teleportNewPos = newPos
	g.teleportAnimation = true
	g.teleportTimer = 0
	g.player = newPos

	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

func (g *Game) isSafePosition(pos Position) bool {
	// Check if any dalek can reach this position in one move
	for _, dalek := range g.daleks {
		if g.distance(pos, dalek.GridPos) <= 2 { // Within one move
			return false
		}
	}
	return true
}

// func (g *Game) useScrewdriverOLD() {

// 	if g.state != StatePlaying || g.screwdrivers <= 0 || (g.daleksMoving && !g.isLastStandActive) {
// 		return
// 	}

// 	g.screwdrivers--

// 	// Find all daleks adjacent to player (including diagonally)
// 	daleksToDestroy := make([]int, 0)
// 	g.screwdriverTargets = make([]Position, 0)

// 	for i, dalek := range g.daleks {
// 		dx := abs(dalek.GridPos.X - g.player.X)
// 		dy := abs(dalek.GridPos.Y - g.player.Y)

// 		// Adjacent includes all 8 surrounding cells
// 		if dx <= 1 && dy <= 1 && (dx != 0 || dy != 0) {
// 			daleksToDestroy = append(daleksToDestroy, i)
// 			g.screwdriverTargets = append(g.screwdriverTargets, dalek.GridPos)
// 		}
// 	}

// 	// Start screwdriver animation if there are targets
// 	if len(g.screwdriverTargets) > 0 {
// 		g.screwdriverAnimation = true
// 		g.soundPlayer.Play("screwdriver")
// 		g.screwdriverTimer = 0
// 	}

// 	// Remove destroyed daleks and add scraps
// 	newDaleks := make([]Dalek, 0, len(g.daleks))
// 	for i, dalek := range g.daleks {
// 		destroyed := false
// 		for _, destroyIndex := range daleksToDestroy {
// 			if i == destroyIndex {
// 				destroyed = true
// 				g.score += 5 // Bonus points for screwdriver kill
// 				// Add debris pile at dalek's position
// 				g.scraps = append(g.scraps, dalek.GridPos)
// 				break
// 			}
// 		}
// 		if !destroyed {
// 			newDaleks = append(newDaleks, dalek)
// 		}
// 	}

// 	g.daleks = newDaleks

// 	// Move remaining daleks after screwdriver use (if not in Last Stand)
// 	if !g.isLastStandActive {
// 		g.moveDaleks()
// 	}
// }

// Updated version to not use debris field when using screwdriver
func (g *Game) useScrewdriver() {

	if g.state != StatePlaying || g.screwdrivers <= 0 || (g.daleksMoving && !g.isLastStandActive) {
		return
	}

	g.screwdrivers--

	// Find all daleks adjacent to player (including diagonally)
	daleksToDestroy := make([]int, 0)
	g.screwdriverTargets = make([]Position, 0)

	for i, dalek := range g.daleks {
		dx := abs(dalek.GridPos.X - g.player.X)
		dy := abs(dalek.GridPos.Y - g.player.Y)

		// Adjacent includes all 8 surrounding cells
		if dx <= 1 && dy <= 1 && (dx != 0 || dy != 0) {
			daleksToDestroy = append(daleksToDestroy, i)
			g.screwdriverTargets = append(g.screwdriverTargets, dalek.GridPos)
		}
	}

	// Start screwdriver animation if there are targets
	if len(g.screwdriverTargets) > 0 {
		g.screwdriverAnimation = true
		g.soundPlayer.Play("screwdriver")
		g.screwdriverTimer = 0
	}

	// Remove destroyed daleks but DON'T add scraps
	newDaleks := make([]Dalek, 0, len(g.daleks))
	for i, dalek := range g.daleks {
		destroyed := false
		for _, destroyIndex := range daleksToDestroy {
			if i == destroyIndex {
				destroyed = true
				g.score += 5 // Bonus points for screwdriver kill
				// REMOVED: Don't add debris pile at dalek's position
				// g.scraps = append(g.scraps, dalek.GridPos)
				break
			}
		}
		if !destroyed {
			newDaleks = append(newDaleks, dalek)
		}
	}

	g.daleks = newDaleks

	// Move remaining daleks after screwdriver use (if not in Last Stand)
	if !g.isLastStandActive {
		g.moveDaleks()
	}
}

func (g *Game) lastStand() {
	if g.state != StatePlaying || g.lastStands <= 0 || g.daleksMoving {
		return
	}

	g.lastStands--
	g.isLastStandActive = true
	g.lastStandSpeed = 2.0 // Reset speed to starting value
	g.daleksMoving = true  // Enable daleks movement for Last Stand
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) moveDaleks() {
	// Start movement animation for all daleks
	g.daleksMoving = true

	for i := range g.daleks {
		dalek := &g.daleks[i]

		// Calculate new grid position
		dx := 0
		dy := 0

		if dalek.GridPos.X < g.player.X {
			dx = 1
		} else if dalek.GridPos.X > g.player.X {
			dx = -1
		}

		if dalek.GridPos.Y < g.player.Y {
			dy = 1
		} else if dalek.GridPos.Y > g.player.Y {
			dy = -1
		}

		newGridPos := Position{
			X: dalek.GridPos.X + dx,
			Y: dalek.GridPos.Y + dy,
		}

		// Update dalek's positions for smooth animation
		dalek.GridPos = newGridPos
		dalek.TargetPos = FloatPosition{
			X: float64(newGridPos.X),
			Y: float64(newGridPos.Y),
		}
		dalek.IsMoving = true
		dalek.MoveTimer = 0
	}
}

func (g *Game) updateDalekAnimations(deltaTime float64) {
	// Ensure Last Stand is properly disabled if game is not in playing state
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	if g.isLastStandActive {
		// Smooth continuous movement during Last Stand
		g.updateLastStandMovement(deltaTime)
	} else {
		// Normal step-by-step movement
		g.updateNormalMovement(deltaTime)
	}
}

// Add this function to check for collision with a threshold
func (g *Game) checkCollisionWithThreshold(pos1, pos2 FloatPosition, threshold float64) bool {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	distSquared := dx*dx + dy*dy
	return distSquared < threshold*threshold
}

// func (g *Game) updateNormalMovementOLD(deltaTime float64) {
// 	allFinished := true

// 	for i := range g.daleks {
// 		dalek := &g.daleks[i]

// 		if dalek.IsMoving {
// 			dalek.MoveTimer += deltaTime

// 			// Calculate interpolation progress (0.0 to 1.0)
// 			progress := dalek.MoveTimer / g.moveAnimationDuration
// 			if progress >= 1.0 {
// 				progress = 1.0
// 				dalek.IsMoving = false
// 				dalek.MoveTimer = 0
// 			} else {
// 				allFinished = false
// 			}

// 			// Smoother easing function (cubic ease-in-out)
// 			var easedProgress float64
// 			if progress < 0.5 {
// 				easedProgress = 4 * progress * progress * progress
// 			} else {
// 				p := progress - 1
// 				easedProgress = 1 + 4*p*p*p
// 			}

// 			// Calculate current visual position
// 			startX := dalek.VisualPos.X
// 			startY := dalek.VisualPos.Y
// 			targetX := dalek.TargetPos.X
// 			targetY := dalek.TargetPos.Y

// 			// Interpolate position
// 			dalek.VisualPos.X = startX + (targetX-startX)*easedProgress
// 			dalek.VisualPos.Y = startY + (targetY-startY)*easedProgress

// 			// Ensure we end up exactly at target
// 			if !dalek.IsMoving {
// 				dalek.VisualPos = dalek.TargetPos
// 			}
// 		}
// 	}

// 	// Check if all daleks finished moving
// 	if allFinished {
// 		g.daleksMoving = false
// 		g.checkCollisions()
// 	}
// }

func (g *Game) updateNormalMovement(deltaTime float64) {
	allFinished := true

	for i := range g.daleks {
		dalek := &g.daleks[i]

		if dalek.IsMoving {
			dalek.MoveTimer += deltaTime

			// Calculate interpolation progress (0.0 to 1.0)
			progress := dalek.MoveTimer / g.moveAnimationDuration
			if progress >= 1.0 {
				progress = 1.0
				dalek.IsMoving = false
				dalek.MoveTimer = 0
			} else {
				allFinished = false
			}

			// Much smoother easing function - smoothstep interpolation
			// This gives a very smooth start and end with consistent middle speed
			var easedProgress float64
			easedProgress = progress * progress * (3.0 - 2.0*progress)

			// Alternative: Even smoother with smootherstep (quintic) - uncomment to try
			// easedProgress = progress * progress * progress * (progress*(progress*6.0-15.0)+10.0)

			// Calculate current visual position
			startX := dalek.VisualPos.X
			startY := dalek.VisualPos.Y
			targetX := dalek.TargetPos.X
			targetY := dalek.TargetPos.Y

			// Interpolate position
			dalek.VisualPos.X = startX + (targetX-startX)*easedProgress
			dalek.VisualPos.Y = startY + (targetY-startY)*easedProgress

			// Ensure we end up exactly at target
			if !dalek.IsMoving {
				dalek.VisualPos = dalek.TargetPos
			}
		}
	}

	// Check if all daleks finished moving
	if allFinished {
		g.daleksMoving = false
		g.checkCollisions()
	}
}

func (g *Game) updateLastStandMovement(deltaTime float64) {
	// Don't continue if game is over
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	// Accelerate the movement speed
	g.lastStandSpeed *= math.Pow(g.lastStandAcceleration, deltaTime)
	if g.lastStandSpeed > g.lastStandMaxSpeed {
		g.lastStandSpeed = g.lastStandMaxSpeed
	}

	anyMoving := false

	// Convert player position to FloatPosition for consistent comparison
	playerPos := FloatPosition{X: float64(g.player.X), Y: float64(g.player.Y)}
	collisionThreshold := 0.5 // Adjust this value to fine-tune collision detection

	// Check player-dalek collisions first
	for _, dalek := range g.daleks {
		if g.checkCollisionWithThreshold(playerPos, dalek.VisualPos, collisionThreshold) {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false
			g.daleksMoving = false
			return
		}
	}

	// Update dalek positions
	for i := range g.daleks {
		dalek := &g.daleks[i]

		dx := playerPos.X - dalek.VisualPos.X
		dy := playerPos.Y - dalek.VisualPos.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > 0.1 { // Still has distance to cover (increased threshold)
			// Normalize direction
			dx /= dist
			dy /= dist

			// Store old position for collision check
			oldPos := dalek.VisualPos

			// Move toward player at current speed
			moveDistance := g.lastStandSpeed * deltaTime
			dalek.VisualPos.X += dx * moveDistance
			dalek.VisualPos.Y += dy * moveDistance

			// Clamp to grid bounds
			dalek.VisualPos.X = math.Max(0, math.Min(float64(gridWidth-1), dalek.VisualPos.X))
			dalek.VisualPos.Y = math.Max(0, math.Min(float64(gridHeight-1), dalek.VisualPos.Y))

			// Update grid position for collision detection
			dalek.GridPos.X = int(math.Round(dalek.VisualPos.X))
			dalek.GridPos.Y = int(math.Round(dalek.VisualPos.Y))

			anyMoving = true

			// Check for collisions with scraps
			for _, scrap := range g.scraps {
				scrapPos := FloatPosition{X: float64(scrap.X), Y: float64(scrap.Y)}
				if g.checkCollisionWithThreshold(dalek.VisualPos, scrapPos, collisionThreshold) {
					dalek.VisualPos = oldPos // Prevent moving through scraps
					g.daleks = append(g.daleks[:i], g.daleks[i+1:]...)
					g.score += 2
					g.soundPlayer.Play("crash")
					// Don't add duplicate scraps
					if !g.positionOccupied(Position{X: int(oldPos.X), Y: int(oldPos.Y)}) {
						g.scraps = append(g.scraps, Position{X: int(oldPos.X), Y: int(oldPos.Y)})
					}
					return
				}
			}

			// Check for collisions with other daleks
			for j := i + 1; j < len(g.daleks); j++ {
				if g.checkCollisionWithThreshold(dalek.VisualPos, g.daleks[j].VisualPos, collisionThreshold) {
					collisionPos := Position{
						X: int((dalek.VisualPos.X + g.daleks[j].VisualPos.X) / 2),
						Y: int((dalek.VisualPos.Y + g.daleks[j].VisualPos.Y) / 2),
					}
					g.daleks = append(g.daleks[:i], g.daleks[i+1:]...)
					g.daleks = append(g.daleks[:j-1], g.daleks[j:]...)
					g.score += 4 // 2 points per dalek
					g.soundPlayer.Play("crash")
					if !g.positionOccupied(collisionPos) {
						g.scraps = append(g.scraps, collisionPos)
					}
					return
				}
			}
		}
	}

	// Check collisions every few frames for better performance
	// but still frequent enough for good responsiveness
	if int(g.lastStandSpeed*deltaTime*60)%2 == 0 {
		g.checkCollisions()
	}

	// If player died during collision check, end Last Stand immediately
	if g.state != StatePlaying {
		g.isLastStandActive = false
		g.daleksMoving = false
		return
	}

	// End Last Stand if no daleks are moving or they're all gone
	if !anyMoving || len(g.daleks) == 0 {
		g.isLastStandActive = false
		g.daleksMoving = false
		if g.state == StatePlaying && len(g.daleks) == 0 {
			g.score += 50 // Bonus for surviving Last Stand
		}
	}
}

func (g *Game) checkCollisions() {
	// Early exit if game is already over
	if g.state != StatePlaying {
		return
	}

	// Check player-dalek collision FIRST (most important check)
	for _, dalek := range g.daleks {
		if g.player == dalek.GridPos {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false // End Last Stand immediately
			g.daleksMoving = false
			return
		}
	}

	// Check dalek-dalek and dalek-scrap collisions
	newDaleks := make([]Dalek, 0, len(g.daleks))
	collidedPositions := make(map[Position]bool)

	// First pass: check for collisions with scraps
	for _, dalek := range g.daleks {
		collided := false

		// Check collision with scraps
		for _, scrap := range g.scraps {
			if dalek.GridPos == scrap {
				collided = true
				g.soundPlayer.Play("crash")
				g.score += 2
				collidedPositions[dalek.GridPos] = true
				break
			}
		}

		if !collided {
			newDaleks = append(newDaleks, dalek)
		}
	}

	// Second pass: check for dalek-dalek collisions
	finalDaleks := make([]Dalek, 0, len(newDaleks))

	for i, dalek := range newDaleks {
		collided := false

		// Check if this dalek collides with any other dalek
		for j, other := range newDaleks {
			if i != j && dalek.GridPos.X == other.GridPos.X && dalek.GridPos.Y == other.GridPos.Y {
				collided = true
				g.score += 2
				collidedPositions[dalek.GridPos] = true
				g.soundPlayer.Play("crash")
				break
			}
		}

		if !collided {
			finalDaleks = append(finalDaleks, dalek)
		}

	}

	// Add debris piles for all collided positions
	for pos := range collidedPositions {
		// Check if scrap already exists at this position
		scrapExists := false
		for _, scrap := range g.scraps {
			if scrap == pos {
				scrapExists = true
				break
			}
		}
		if !scrapExists {
			g.scraps = append(g.scraps, pos)
		}
	}

	g.daleks = finalDaleks

	// Check player-dalek collision again after updating dalek positions
	// (in case daleks moved onto player during collision resolution)
	for _, dalek := range g.daleks {
		if g.player == dalek.GridPos {
			g.state = StateGameOver
			g.soundPlayer.Play("gameover")
			g.gameOverMessage = "Game Over! You were caught by a Dalek!"
			g.isLastStandActive = false // End Last Stand immediately
			g.daleksMoving = false
			return
		}
	}

	// Check if level is complete
	if len(g.daleks) == 0 {
		g.score += g.level * 10
		g.level++
		g.teleports += 2
		// if g.level%3 == 0 { // Bonus screwdriver every 3 levels
		// 	g.screwdrivers++
		// }
		g.screwdrivers += 2 // Increase screwdrivers by 2 every level
		if g.level%5 == 0 { // Bonus last stand every 5 levels
			g.lastStands++
		}
		if g.level > 10 {
			g.state = StateWin
			g.gameOverMessage = "Congratulations! You survived all levels!"
			g.soundPlayer.Play("gameover")
		} else {
			g.startLevel()
		}
	}
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

func (g *Game) drawMenu(screen *ebiten.Image) {
	title := "GODALEKS - alpha v0.04"
	text.Draw(screen, title, basicfont.Face7x13, screenWidth/2-len(title)*3, 100, color.Black)

	instructions := []string{
		"Use arrow keys or mouse to move",
		"Q, E, Z, C for diagonal movement",
		"N To start a new game",
		"SPACE or . to wait",
		"T to teleport randomly",
		"R to teleport safely",
		"S to use sonic screwdriver",
		"L for Last Stand (all daleks rush you)",
		"G to turn game grid On/Off",
		"",
		"MOUSE: Click adjacent cell to move there",
		"Click on player to wait in place",
		"",
		"Avoid the Daleks!",
		"Make them crash into each other!",
		"Sonic Screwdriver destroys adjacent Daleks!",
		"Last Stand forces all daleks to move!",
		"",
		"Press SPACE or click to start",
	}

	for i, line := range instructions {
		text.Draw(screen, line, basicfont.Face7x13, 50, 200+i*20, color.Black)
	}
}

func (g *Game) drawMouseIndicator(screen *ebiten.Image) {
	if g.state != StatePlaying {
		return
	}

	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Convert to grid coordinates
	gridX, gridY, valid := g.screenToGrid(mouseX, mouseY)
	if !valid {
		return
	}

	targetPos := Position{X: gridX, Y: gridY}

	// Check if it's a valid move (adjacent to player)
	dx := abs(targetPos.X - g.player.X)
	dy := abs(targetPos.Y - g.player.Y)

	// Only show indicator for valid moves or current position
	if dx <= 1 && dy <= 1 {
		offsetX := (screenWidth - gridWidth*cellSize) / 2
		offsetY := 50

		x := float64(offsetX + gridX*cellSize)
		y := float64(offsetY + gridY*cellSize)

		// Choose color based on move type
		var indicatorColor color.Color
		if targetPos == g.player {
			indicatorColor = color.RGBA{0, 255, 0, 100} // Green for wait/current position
		} else {
			// Check if position is occupied by scrap
			occupied := false
			for _, scrap := range g.scraps {
				if scrap == targetPos {
					occupied = true
					break
				}
			}

			if occupied {
				indicatorColor = color.RGBA{255, 0, 0, 100} // Red for blocked
			} else {
				indicatorColor = color.RGBA{0, 0, 255, 100} // Blue for valid move
			}
		}

		// Draw semi-transparent overlay on the cell
		ebitenutil.DrawRect(screen, x, y, cellSize, cellSize, indicatorColor)

		// Draw border
		ebitenutil.DrawRect(screen, x, y, cellSize, 1, color.Black)
		ebitenutil.DrawRect(screen, x, y, 1, cellSize, color.Black)
		ebitenutil.DrawRect(screen, x+cellSize-1, y, 1, cellSize, color.Black)
		ebitenutil.DrawRect(screen, x, y+cellSize-1, cellSize, 1, color.Black)
	}
}

func (g *Game) drawTeleportEffect(screen *ebiten.Image, pos Position, progress float64, offsetX, offsetY int) {
	x := float64(offsetX + pos.X*cellSize + cellSize/2)
	y := float64(offsetY + pos.Y*cellSize + cellSize/2)

	// Create sparkle/energy effect with black particles for classic Mac style
	numParticles := 8
	radius := float64(cellSize) * (1.0 + progress*2.0) // Expanding circle

	for i := 0; i < numParticles; i++ {
		angle := float64(i)*2.0*3.14159/float64(numParticles) + progress*6.28 // Rotating
		px := x + radius*0.5*math.Cos(angle)
		py := y + radius*0.5*math.Sin(angle)

		// Fade out over time
		alpha := uint8(255 * (1.0 - progress))
		sparkleColor := color.RGBA{0x00, 0x00, 0x00, alpha} // Black sparkles

		// Draw sparkle particles
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if int(px)+dx >= 0 && int(px)+dx < screenWidth &&
					int(py)+dy >= 0 && int(py)+dy < screenHeight {
					screen.Set(int(px)+dx, int(py)+dy, sparkleColor)
				}
			}
		}
	}

	// Central flash effect with black
	flashAlpha := uint8(255 * (1.0 - progress) * 0.8)
	flashColor := color.RGBA{0x00, 0x00, 0x00, flashAlpha} // Black flash

	flashRadius := int(float64(cellSize/2) * (1.0 - progress*0.5))
	for dx := -flashRadius; dx <= flashRadius; dx++ {
		for dy := -flashRadius; dy <= flashRadius; dy++ {
			if dx*dx+dy*dy <= flashRadius*flashRadius {
				px := int(x) + dx
				py := int(y) + dy
				if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
					screen.Set(px, py, flashColor)
				}
			}
		}
	}
}

func (g *Game) drawScrewdriverEffect(screen *ebiten.Image, pos Position, progress float64, offsetX, offsetY int) {
	x := float64(offsetX + pos.X*cellSize + cellSize/2)
	y := float64(offsetY + pos.Y*cellSize + cellSize/2)

	// Create electric/energy effect for sonic screwdriver
	numBolts := 6
	radius := float64(cellSize) * 0.8

	for i := 0; i < numBolts; i++ {
		angle := float64(i)*2.0*3.14159/float64(numBolts) + progress*12.56 // Fast rotation

		// Create zigzag lightning bolt effect
		for step := 0; step < 5; step++ {
			stepRadius := radius * float64(step) / 5.0
			zigzag := math.Sin(progress*20.0+float64(step)) * 3.0 // Zigzag offset

			boltX := x + stepRadius*math.Cos(angle) + zigzag*math.Cos(angle+1.57)
			boltY := y + stepRadius*math.Sin(angle) + zigzag*math.Sin(angle+1.57)

			// Fade based on progress and distance
			alpha := uint8(255 * (1.0 - progress) * (1.0 - float64(step)/5.0))
			boltColor := color.RGBA{0x00, 0x00, 0x00, alpha} // Black lightning

			if int(boltX) >= 0 && int(boltX) < screenWidth &&
				int(boltY) >= 0 && int(boltY) < screenHeight {
				screen.Set(int(boltX), int(boltY), boltColor)
				// Make bolts thicker
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						px := int(boltX) + dx
						py := int(boltY) + dy
						if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
							screen.Set(px, py, boltColor)
						}
					}
				}
			}
		}
	}

	// Central pulse effect
	pulseRadius := int(float64(cellSize/3) * (1.0 + progress*2.0))
	pulseAlpha := uint8(255 * (1.0 - progress) * 0.6)
	pulseColor := color.RGBA{0x00, 0x00, 0x00, pulseAlpha}

	for dx := -pulseRadius; dx <= pulseRadius; dx++ {
		for dy := -pulseRadius; dy <= pulseRadius; dy++ {
			if dx*dx+dy*dy <= pulseRadius*pulseRadius {
				px := int(x) + dx
				py := int(y) + dy
				if px >= 0 && px < screenWidth && py >= 0 && py < screenHeight {
					screen.Set(px, py, pulseColor)
				}
			}
		}
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	offsetX := (screenWidth - gridWidth*cellSize) / 2
	offsetY := 50

	// Draw grid only if enabled
	if g.showGrid {
		for x := 0; x <= gridWidth; x++ {
			ebitenutil.DrawLine(screen,
				float64(offsetX+x*cellSize), float64(offsetY),
				float64(offsetX+x*cellSize), float64(offsetY+gridHeight*cellSize),
				color.Black)
		}
		for y := 0; y <= gridHeight; y++ {
			ebitenutil.DrawLine(screen,
				float64(offsetX), float64(offsetY+y*cellSize),
				float64(offsetX+gridWidth*cellSize), float64(offsetY+y*cellSize),
				color.Black)
		}
	}

	// Draw scraps (centered)
	for _, scrap := range g.scraps {
		x, y := getCenteredSpritePosition(scrap.X, scrap.Y, offsetX, offsetY, g.scrapImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.scrapImage, op)
	}

	// Draw daleks using smooth interpolated positions (centered)
	for _, dalek := range g.daleks {
		// Use visual position for smooth movement, but calculate centered position
		cellCenterX := float64(offsetX) + dalek.VisualPos.X*float64(cellSize) + float64(cellSize)/2
		cellCenterY := float64(offsetY) + dalek.VisualPos.Y*float64(cellSize) + float64(cellSize)/2

		// Get sprite dimensions and center it
		spriteBounds := g.dalekImage.Bounds()
		spriteWidth := spriteBounds.Dx()
		spriteHeight := spriteBounds.Dy()

		x := cellCenterX - float64(spriteWidth)/2
		y := cellCenterY - float64(spriteHeight)/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.dalekImage, op)
	}

	// Draw player with teleportation effects (centered)
	if g.teleportAnimation {
		progress := g.teleportTimer / 0.5 // 0.5 second animation

		// Draw disappearing effect at old position
		if progress < 0.5 {
			g.drawTeleportEffect(screen, g.teleportOldPos, progress*2, offsetX, offsetY)
		}

		// Draw appearing effect at new position
		if progress > 0.3 {
			appearProgress := (progress - 0.3) / 0.7
			g.drawTeleportEffect(screen, g.teleportNewPos, 1.0-appearProgress, offsetX, offsetY)
		}

		// Draw player with fade effect (centered)
		x, y := getCenteredSpritePosition(g.player.X, g.player.Y, offsetX, offsetY, g.playerImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)

		// Fade in player at new position
		if progress > 0.3 {
			fadeAlpha := (progress - 0.3) / 0.7
			op.ColorM.Scale(1, 1, 1, fadeAlpha)
		} else {
			op.ColorM.Scale(1, 1, 1, 0) // Invisible during first part
		}

		screen.DrawImage(g.playerImage, op)
	} else {
		// Normal player drawing (centered)
		x, y := getCenteredSpritePosition(g.player.X, g.player.Y, offsetX, offsetY, g.playerImage)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.playerImage, op)
	}

	// Draw screwdriver effects
	if g.screwdriverAnimation {
		progress := g.screwdriverTimer / 0.8 // 0.8 second animation

		for _, target := range g.screwdriverTargets {
			g.drawScrewdriverEffect(screen, target, progress, offsetX, offsetY)
		}
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	// Semi-transparent overlay
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 128})

	// Game over message
	text.Draw(screen, g.gameOverMessage, basicfont.Face7x13,
		screenWidth/2-len(g.gameOverMessage)*3, screenHeight/2-20, color.White)

	finalScore := fmt.Sprintf("Final Score: %d", g.score)
	text.Draw(screen, finalScore, basicfont.Face7x13,
		screenWidth/2-len(finalScore)*3, screenHeight/2+10, color.White)

	restart := "Press SPACE or click to restart"
	text.Draw(screen, restart, basicfont.Face7x13,
		screenWidth/2-len(restart)*3, screenHeight/2+40, color.White)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	// Status information
	status := fmt.Sprintf("Level: %d  Score: %d  Teleports: %d  Safe: %d  Screwdrivers: %d  Last Stands: %d  Daleks: %d",
		g.level, g.score, g.teleports, g.safeTeleports, g.screwdrivers, g.lastStands, len(g.daleks))
	text.Draw(screen, status, basicfont.Face7x13, 10, 20, color.Black)

	// Grid indicator
	gridStatus := "Grid: OFF"
	if g.showGrid {
		gridStatus = "Grid: ON"
	}
	text.Draw(screen, gridStatus, basicfont.Face7x13, 10, 40, color.Black)

	// Last Stand indicator
	if g.isLastStandActive {
		lastStandMsg := fmt.Sprintf("LAST STAND ACTIVE! Speed: %.1f", g.lastStandSpeed)
		text.Draw(screen, lastStandMsg, basicfont.Face7x13, 10, screenHeight-30, color.Black)
	}

	// Temporary center-screen notification for grid toggle
	if g.gridToggleMessage != "" && time.Since(g.gridToggleMessageTime) < 1500*time.Millisecond {
		msg := g.gridToggleMessage
		x := screenWidth/2 - len(msg)*3
		y := 60
		text.Draw(screen, msg, basicfont.Face7x13, x, y, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
