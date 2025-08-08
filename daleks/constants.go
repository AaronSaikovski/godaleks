package daleks

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 50
	gridHeight   = 35 // Reduced from 37 to 35 to ensure sprites stay in bounds
	cellSize     = 16
)

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
	StateWin
)
