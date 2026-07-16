// MIT License

// Copyright (c) 2025 - Aaron Saikovski <asaikovski@outlook.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package cmd

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X, Y int
}

type FloatPosition struct {
	X, Y float64
}

type CollisionEffect struct {
	Pos      FloatPosition // Position of the collision
	Timer    float64       // Animation timer
	Duration float64       // Total duration
}

type Dalek struct {
	GridPos   Position      // Current grid position
	VisualPos FloatPosition // Interpolated visual position (current frame)
	StartPos  FloatPosition // Starting position for animation
	TargetPos FloatPosition // Target visual position
	IsMoving  bool          // Whether currently animating
	MoveTimer float64       // Animation timer
	IsEmperor bool          // Whether this is the Dalek Emperor (can only be killed when all normal daleks are dead)
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

	playerImage  *ebiten.Image
	dalekImage   *ebiten.Image
	emperorImage *ebiten.Image
	scrapImage   *ebiten.Image
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
	// Level transition
	levelCompleteTimer float64
	// Emperor warning
	emperorWarningMessage  string
	emperorWarningTimer    float64
	emperorWarningDuration float64
	// Last Stand smooth movement
	lastStandSpeed        float64 // Speed in cells per second during Last Stand
	lastStandAcceleration float64 // Acceleration multiplier per second
	lastStandMaxSpeed     float64 // Maximum speed cap
	// Mouse support
	lastClickTime time.Time
	soundPlayer   *SoundPlayer

	// Performance optimizations - pre-allocated objects
	drawOp           ebiten.DrawImageOptions // Reusable draw options
	playerHalfWidth  float64                 // Cached sprite dimensions
	playerHalfHeight float64
	dalekHalfWidth   float64
	dalekHalfHeight  float64
	scrapHalfWidth   float64
	scrapHalfHeight  float64
	hudStatusText    string    // Cached HUD text
	lastHUDUpdate    time.Time // Track when HUD needs refresh

	// Cached draw strings to avoid per-frame fmt.Sprintf
	cachedFinalScore   string
	cachedGridStatus   string
	cachedLevelMsg     string
	cachedLevelNextMsg string

	// Collision animations
	collisionEffects []CollisionEffect // Active collision effects

	// Grid-based occupancy array for O(1) lookups with zero GC pressure
	occupancyGrid [gridWidth][gridHeight]bool
	// Grid-based scrap lookup for O(1) checks
	scrapGrid [gridWidth][gridHeight]bool

	// Reusable buffers to reduce allocations
	dalekRemoveMap     map[int]bool       // Reusable for dalek removal
	scrapMap           map[Position]bool  // Reusable for scrap lookups in collision
	daleksByPosition   map[Position][]int // Reusable for collision grouping
	newScrapsBuf       []Position         // Reusable scrap buffer for Last Stand
	survivingDaleksBuf []Dalek            // Reusable dalek buffer for filtering

	// Pre-computed trig lookup table for particle effects
	sinTable [trigTableSize]float64
	cosTable [trigTableSize]float64
}
