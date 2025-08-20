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
package daleks

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
