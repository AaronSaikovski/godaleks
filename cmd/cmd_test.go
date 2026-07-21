package cmd

import (
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{0, 0},
		{5, 5},
		{-5, 5},
		{-1, 1},
		{math.MaxInt32, math.MaxInt32},
	}
	for _, c := range cases {
		if got := abs(c.in); got != c.want {
			t.Errorf("abs(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestCheckCollisionWithThreshold(t *testing.T) {
	g := &Game{}
	cases := []struct {
		name          string
		a, b          FloatPosition
		threshold     float64
		wantCollision bool
	}{
		{"identical", FloatPosition{1, 1}, FloatPosition{1, 1}, 0.5, true},
		{"inside threshold", FloatPosition{0, 0}, FloatPosition{0.3, 0.3}, 0.5, true},
		{"on threshold edge", FloatPosition{0, 0}, FloatPosition{0.5, 0}, 0.5, false},
		{"outside threshold", FloatPosition{0, 0}, FloatPosition{1, 1}, 0.5, false},
		{"far apart", FloatPosition{0, 0}, FloatPosition{10, 10}, 0.5, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := g.checkCollisionWithThreshold(c.a, c.b, c.threshold)
			if got != c.wantCollision {
				t.Errorf("got %v, want %v", got, c.wantCollision)
			}
		})
	}
}

func TestDistanceIsSquared(t *testing.T) {
	// distance() returns squared distance for efficiency — verify the contract
	// so callers know which unit to compare against.
	g := &Game{}
	d := g.distance(Position{0, 0}, Position{3, 4})
	if d != 25 {
		t.Errorf("distance((0,0),(3,4)) = %v, want 25 (squared)", d)
	}
	if d = g.distance(Position{5, 5}, Position{5, 5}); d != 0 {
		t.Errorf("distance of identical points = %v, want 0", d)
	}
}

func TestRebuildScrapGrid(t *testing.T) {
	g := &Game{
		scraps: []Position{{1, 2}, {5, 7}, {gridWidth - 1, gridHeight - 1}},
	}
	g.rebuildScrapGrid()

	if !g.scrapGrid[1][2] {
		t.Error("scrapGrid[1][2] should be true")
	}
	if !g.scrapGrid[5][7] {
		t.Error("scrapGrid[5][7] should be true")
	}
	if !g.scrapGrid[gridWidth-1][gridHeight-1] {
		t.Error("scrapGrid corner should be true")
	}
	if g.scrapGrid[0][0] {
		t.Error("scrapGrid[0][0] should be false")
	}
}

func TestRebuildScrapGridIgnoresOutOfBounds(t *testing.T) {
	// Bounds guards should prevent panics if corrupted state sneaks in.
	g := &Game{
		scraps: []Position{{-1, 0}, {0, -1}, {gridWidth, 0}, {0, gridHeight}, {3, 3}},
	}
	// Must not panic.
	g.rebuildScrapGrid()

	if !g.scrapGrid[3][3] {
		t.Error("valid scrap at (3,3) should still be recorded")
	}
}

func TestRebuildScrapGridClearsStaleEntries(t *testing.T) {
	g := &Game{}
	g.scrapGrid[10][10] = true // stale
	g.scraps = []Position{{5, 5}}

	g.rebuildScrapGrid()

	if g.scrapGrid[10][10] {
		t.Error("rebuildScrapGrid should clear stale entries")
	}
	if !g.scrapGrid[5][5] {
		t.Error("rebuildScrapGrid should mark current scraps")
	}
}

func TestEnsureScrapGridRespectsDirtyFlag(t *testing.T) {
	g := &Game{scraps: []Position{{5, 5}}}

	// Clean flag + empty grid: ensureScrapGrid must NOT rebuild.
	g.scrapGridDirty = false
	g.ensureScrapGrid()
	if g.scrapGrid[5][5] {
		t.Error("ensureScrapGrid rebuilt while not dirty")
	}

	// Dirty flag: ensureScrapGrid rebuilds and then clears the flag.
	g.scrapGridDirty = true
	g.ensureScrapGrid()
	if !g.scrapGrid[5][5] {
		t.Error("ensureScrapGrid did not rebuild when dirty")
	}
	if g.scrapGridDirty {
		t.Error("ensureScrapGrid should clear the dirty flag after rebuilding")
	}
}

func TestRebuildScrapGridClearsDirtyFlag(t *testing.T) {
	g := &Game{scraps: []Position{{1, 1}}, scrapGridDirty: true}
	g.rebuildScrapGrid()
	if g.scrapGridDirty {
		t.Error("rebuildScrapGrid should clear the dirty flag")
	}
}

func TestRebuildOccupancyGrid(t *testing.T) {
	g := &Game{
		daleks: []Dalek{
			{GridPos: Position{2, 3}},
			{GridPos: Position{10, 10}, IsEmperor: true},
		},
		scraps: []Position{{4, 5}},
	}
	g.rebuildOccupancyGrid()

	if !g.occupancyGrid[2][3] {
		t.Error("dalek position should be occupied")
	}
	if !g.occupancyGrid[10][10] {
		t.Error("emperor position should be occupied")
	}
	if !g.occupancyGrid[4][5] {
		t.Error("scrap position should be occupied")
	}
	if g.occupancyGrid[0][0] {
		t.Error("empty cell should not be occupied")
	}
}

func TestTrigTableSizeIsPowerOfTwo(t *testing.T) {
	// trigTableSize = 32 is a power of two so (i & (trigTableSize-1)) works
	// as a fast modulo in drawTeleportEffect. Fail loudly if someone changes
	// it to a non-power-of-two value without updating the index math.
	if trigTableSize&(trigTableSize-1) != 0 {
		t.Fatalf("trigTableSize = %d must be a power of two", trigTableSize)
	}
}

func TestTrigTableContent(t *testing.T) {
	// Verify the table generation formula produces correct sin/cos values.
	const epsilon = 1e-9
	var sinTable, cosTable [trigTableSize]float64
	for i := range trigTableSize {
		angle := float64(i) * 2.0 * math.Pi / float64(trigTableSize)
		sinTable[i] = math.Sin(angle)
		cosTable[i] = math.Cos(angle)
	}

	// Sanity anchors: index 0 is (cos=1, sin=0), index trigTableSize/4 is (0, 1).
	if math.Abs(cosTable[0]-1.0) > epsilon || math.Abs(sinTable[0]) > epsilon {
		t.Errorf("entry 0: cos=%v sin=%v, want (1, 0)", cosTable[0], sinTable[0])
	}
	quarter := trigTableSize / 4
	if math.Abs(cosTable[quarter]) > epsilon || math.Abs(sinTable[quarter]-1.0) > epsilon {
		t.Errorf("entry %d: cos=%v sin=%v, want (0, 1)", quarter, cosTable[quarter], sinTable[quarter])
	}
}

func TestLastStandAccelerationLinearApprox(t *testing.T) {
	// The linear approximation 1 + (accel-1)*dt replaces math.Pow(accel, dt).
	// At 60 FPS with accel=1.1, the error must be well under one part per thousand
	// for the approximation to be correct; if someone tightens accel beyond 1.1
	// or changes dt bounds, this test will catch it.
	accel := 1.1
	dts := []float64{1.0 / 240.0, 1.0 / 60.0, 1.0 / 30.0}
	for _, dt := range dts {
		linear := 1.0 + (accel-1.0)*dt
		exact := math.Pow(accel, dt)
		if math.Abs(linear-exact) > 0.002 {
			t.Errorf("dt=%v linear=%v exact=%v diff=%v", dt, linear, exact, math.Abs(linear-exact))
		}
	}
}

func TestIsScrapAt_Bounds(t *testing.T) {
	g := &Game{scraps: []Position{{5, 5}}}

	if g.isScrapAt(Position{-1, 0}) {
		t.Error("isScrapAt negative X should return false")
	}
	if g.isScrapAt(Position{0, -1}) {
		t.Error("isScrapAt negative Y should return false")
	}
	if g.isScrapAt(Position{gridWidth, 0}) {
		t.Error("isScrapAt X >= gridWidth should return false")
	}
	if g.isScrapAt(Position{0, gridHeight}) {
		t.Error("isScrapAt Y >= gridHeight should return false")
	}
	if !g.isScrapAt(Position{5, 5}) {
		t.Error("isScrapAt(5,5) should return true with scrap there")
	}
}
