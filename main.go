.
package main

import (
	"log"

	"github.com/AaronSaikovski/godaleks/cmd"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// main - Main function
func main() {

	game := cmd.NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GoDaleks v1.2.2")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
