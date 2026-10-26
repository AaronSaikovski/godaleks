package cmd

import (
	"fmt"
	"math"
	"testing"
)

func newLevelMovementGame(level int) *Game {
	return &Game{
		state:                 StatePlaying,
		level:                 level,
		player:                Position{30, 10},
		daleks:                []Dalek{{GridPos: Position{5, 10}, VisualPos: FloatPosition{5, 10}}},
		moveAnimationDuration: 0.7,
		lastStandSpeed:        2,
		lastStandAcceleration: 1.1,
		lastStandMaxSpeed:     8,
		dalekRemoveMap:        make(map[int]bool),
		scrapMap:              make(map[Position]bool),
		daleksByPosition:      make(map[Position][]int),
	}
}

func TestLevelCompletionThroughFifteen(t *testing.T) {
	for _, tc := range []struct {
		level int
		state GameState
	}{{10, StateLevelComplete}, {14, StateLevelComplete}, {15, StateWin}} {
		t.Run(fmt.Sprint(tc.level), func(t *testing.T) {
			g := newLevelMovementGame(tc.level)
			// Clear the final normal Dalek through actual scrap collision handling.
			g.scraps = []Position{g.daleks[0].GridPos}
			g.checkCollisions()
			if g.state != tc.state || g.level != tc.level+1 {
				t.Fatalf("state/level = %v/%d, want %v/%d", g.state, g.level, tc.state, tc.level+1)
			}
			if g.score != 2+tc.level*10 || g.teleports != 2 || g.screwdrivers != 2 {
				t.Fatalf("completion rewards incorrect: score=%d teleports=%d screwdrivers=%d", g.score, g.teleports, g.screwdrivers)
			}
		})
	}
}

func TestLevelMovementScaling(t *testing.T) {
	for level := 1; level <= 15; level++ {
		t.Run(fmt.Sprint(level), func(t *testing.T) {
			wantMultiplier := 1.0
			if level > 10 {
				wantMultiplier += float64(level-10) / 5
			}
			if got := levelSpeedMultiplier(level); math.Abs(got-wantMultiplier) > 1e-12 {
				t.Fatalf("multiplier = %g, want %g", got, wantMultiplier)
			}
			g := newLevelMovementGame(level)
			g.moveDaleks()
			if g.daleks[0].GridPos != (Position{6, 10}) {
				t.Fatal("normal movement must remain one grid step")
			}
			g.updateNormalMovement(0.1)
			progress := 0.1 * wantMultiplier / 0.7
			wantX := 5 + progress*progress*progress*(progress*(progress*6-15)+10)
			if math.Abs(g.daleks[0].VisualPos.X-wantX) > 1e-12 {
				t.Fatalf("normal visual X = %g, want %g", g.daleks[0].VisualPos.X, wantX)
			}
			g.updateNormalMovement(0.7/wantMultiplier - 0.1 + 1e-12)
			if g.daleksMoving || g.daleks[0].IsMoving || g.daleks[0].VisualPos != (FloatPosition{6, 10}) {
				t.Fatal("normal animation did not finish at its single-cell target")
			}
			for _, initialSpeed := range []float64{2, 8} {
				g = newLevelMovementGame(level)
				g.lastStandSpeed = initialSpeed
				g.isLastStandActive = true
				g.updateLastStandMovement(0.1)
				wantSpeed := math.Min(initialSpeed*1.01, 8)
				wantX = 5 + wantSpeed*0.1*wantMultiplier
				if math.Abs(g.lastStandSpeed-wantSpeed) > 1e-12 || math.Abs(g.daleks[0].VisualPos.X-wantX) > 1e-12 {
					t.Fatalf("Last Stand speed/X = %g/%g, want %g/%g", g.lastStandSpeed, g.daleks[0].VisualPos.X, wantSpeed, wantX)
				}
			}
		})
	}
}

func TestResetRestoresEarlyLevelMovement(t *testing.T) {
	g := newLevelMovementGame(15)
	g.moveDaleks()
	g.updateNormalMovement(0.1)
	lateNormalX := g.daleks[0].VisualPos.X
	g.resetGame()
	if g.level != 1 || g.lastStandSpeed != 2 {
		t.Fatalf("reset level/speed = %d/%g", g.level, g.lastStandSpeed)
	}
	// Keep reset settings, replacing only the randomized board with a fixed lane.
	g.player = Position{30, 10}
	g.daleks = []Dalek{{GridPos: Position{5, 10}, VisualPos: FloatPosition{5, 10}}}
	g.moveDaleks()
	g.updateNormalMovement(0.1)
	baseline := newLevelMovementGame(1)
	baseline.moveDaleks()
	baseline.updateNormalMovement(0.1)
	if g.daleks[0].VisualPos != baseline.daleks[0].VisualPos || g.daleks[0].VisualPos.X >= lateNormalX {
		t.Fatal("reset did not restore original normal animation speed")
	}
	g.daleks = []Dalek{{GridPos: Position{5, 10}, VisualPos: FloatPosition{5, 10}}}
	g.isLastStandActive = true
	g.updateLastStandMovement(0.1)
	baseline = newLevelMovementGame(1)
	baseline.isLastStandActive = true
	baseline.updateLastStandMovement(0.1)
	if g.daleks[0].VisualPos != baseline.daleks[0].VisualPos {
		t.Fatal("reset did not restore original Last Stand movement speed")
	}
}
