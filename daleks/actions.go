package daleks

import (
	"math"
	"math/rand"
	"time"
)

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
		g.screwdriverTimer = 0
	}

	// Remove destroyed daleks and add scraps
	newDaleks := make([]Dalek, 0, len(g.daleks))
	for i, dalek := range g.daleks {
		destroyed := false
		for _, destroyIndex := range daleksToDestroy {
			if i == destroyIndex {
				destroyed = true
				g.score += 5 // Bonus points for screwdriver kill
				// Add debris pile at dalek's position
				g.scraps = append(g.scraps, dalek.GridPos)
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

			// Smooth easing function (ease-out)
			easedProgress := 1.0 - (1.0-progress)*(1.0-progress)

			// Calculate current visual position
			startX := dalek.VisualPos.X
			startY := dalek.VisualPos.Y
			targetX := dalek.TargetPos.X
			targetY := dalek.TargetPos.Y

			// If this is the start of movement, set the start position
			if dalek.MoveTimer <= deltaTime {
				startX = dalek.VisualPos.X
				startY = dalek.VisualPos.Y
			}

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

	for i := range g.daleks {
		dalek := &g.daleks[i]

		// Calculate direction toward player
		playerPos := FloatPosition{X: float64(g.player.X), Y: float64(g.player.Y)}

		dx := playerPos.X - dalek.VisualPos.X
		dy := playerPos.Y - dalek.VisualPos.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > 0.1 { // Still has distance to cover (increased threshold)
			// Normalize direction
			dx /= dist
			dy /= dist

			// Move toward player at current speed
			moveDistance := g.lastStandSpeed * deltaTime
			dalek.VisualPos.X += dx * moveDistance
			dalek.VisualPos.Y += dy * moveDistance

			// Update grid position for collision detection
			dalek.GridPos.X = int(math.Round(dalek.VisualPos.X))
			dalek.GridPos.Y = int(math.Round(dalek.VisualPos.Y))

			// Clamp to grid bounds
			if dalek.GridPos.X < 0 {
				dalek.GridPos.X = 0
				dalek.VisualPos.X = 0
			} else if dalek.GridPos.X >= gridWidth {
				dalek.GridPos.X = gridWidth - 1
				dalek.VisualPos.X = float64(gridWidth - 1)
			}

			if dalek.GridPos.Y < 0 {
				dalek.GridPos.Y = 0
				dalek.VisualPos.Y = 0
			} else if dalek.GridPos.Y >= gridHeight {
				dalek.GridPos.Y = gridHeight - 1
				dalek.VisualPos.Y = float64(gridHeight - 1)
			}

			anyMoving = true
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
		if g.level%3 == 0 { // Bonus screwdriver every 3 levels
			g.screwdrivers++
		}
		if g.level%5 == 0 { // Bonus last stand every 5 levels
			g.lastStands++
		}
		if g.level > 10 {
			g.state = StateWin
			g.gameOverMessage = "Congratulations! You survived all levels!"
		} else {
			g.startLevel()
		}
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
