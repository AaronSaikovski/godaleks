package main

import (
	"log"

	"github.com/AaronSaikovski/godaleks/daleks"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// main - Main function
func main() {

	game := daleks.NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Daleks")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
