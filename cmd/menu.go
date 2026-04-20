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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// drawBoldText draws smooth, non-bold text
func drawBoldText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	// Draw single clean text for smooth appearance
	text.Draw(screen, str, basicfont.Face7x13, x, y, clr)
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	title := "GoDaleks - v1.3.0"
	drawBoldText(screen, title, screenWidth/2-len(title)*3, 100, color.Black)

	gameDesc := "Based on the original 1984 'Daleks' Macintosh Classic game by Johan Strandberg."
	drawBoldText(screen, gameDesc, screenWidth/2-len(gameDesc)*7/2, 120, color.Black)

	instructions := []string{
		"Use arrow keys or mouse to move",
		"Q, E, Z, C for diagonal movement",
		"N To start a new game",
		"SPACE or . to wait",
		"T to teleport randomly",
		"R to teleport safely",
		"S to use sonic screwdriver",
		"L for Last Stand (all daleks rush you)",
		"G to turn game grid On/Off",
		"",
		"MOUSE: Click adjacent cell to move there",
		"Click on player to wait in place",
		"",
		"Avoid the Daleks!",
		"Make them crash into each other!",
		"Watch out for the Emperor Dalek",
		"Sonic Screwdriver destroys adjacent Daleks!",
		"Last Stand forces all daleks to move!",
		"",
		"Press SPACE or click to start",
	}

	for i, line := range instructions {
		drawBoldText(screen, line, 50, 200+i*15, color.Black)
	}

	authorInfo := "Ported to Go by: Aaron Saikovski - https://github.com/AaronSaikovski/godaleks"
	drawBoldText(screen, authorInfo, screenWidth/2-len(authorInfo)*7/2, screenHeight-30, color.Black)
}
