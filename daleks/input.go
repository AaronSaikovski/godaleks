package daleks

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct{}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Update() (ebiten.Key, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return ebiten.KeyArrowUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return ebiten.KeyArrowLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return ebiten.KeyArrowRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return ebiten.KeyArrowDown, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		return ebiten.KeyN, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		//fmt.Print("T pressed")
		return ebiten.KeyT, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return ebiten.KeyS, true
	}

	return 0, false
}
