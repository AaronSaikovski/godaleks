package daleks

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

func (g *Game) isSafePosition(pos Position) bool {
	// Check if any dalek can reach this position in one move
	for _, dalek := range g.daleks {
		if g.distance(pos, dalek.GridPos) <= 2 { // Within one move
			return false
		}
	}
	return true
}
