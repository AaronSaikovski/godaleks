package main

import (
	"log"

	"github.com/AaronSaikovski/gorobots/pkg/daleks"
	"github.com/hajimehoshi/ebiten/v2"
)

// main - Main function
func main() {

	game := daleks.NewGame()

	ebiten.SetWindowSize(daleks.ScreenWidth, daleks.ScreenHeight)
	ebiten.SetWindowTitle("GoRobots")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
