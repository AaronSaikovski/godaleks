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
)

// Helper function to calculate centered sprite position
func getCenteredSpritePosition(gridX, gridY, offsetX, offsetY int, spriteImage *ebiten.Image) (float64, float64) {
	// Get sprite dimensions
	spriteBounds := spriteImage.Bounds()
	spriteWidth := spriteBounds.Dx()
	spriteHeight := spriteBounds.Dy()

	// Calculate center position within the grid cell
	cellCenterX := float64(offsetX + gridX*cellSize + cellSize/2)
	cellCenterY := float64(offsetY + gridY*cellSize + cellSize/2)

	// Subtract half sprite size to center it
	x := cellCenterX - float64(spriteWidth)/2
	y := cellCenterY - float64(spriteHeight)/2

	return x, y
}

// createScrapImage creates a Dalek debris pile sprite like in the classic game
func createScrapImage() *ebiten.Image {
	img := ebiten.NewImage(cellSize-2, cellSize-2)

	size := cellSize - 2
	centerX := size / 2
	centerY := size / 2

	// Create Dalek debris - looks like scattered Dalek parts and metal fragments
	// Main debris cluster in center
	debrisPositions := []Position{
		// Central cluster
		{centerX, centerY},
		{centerX - 1, centerY},
		{centerX + 1, centerY},
		{centerX, centerY - 1},
		{centerX, centerY + 1},
		{centerX - 1, centerY - 1},
		{centerX + 1, centerY + 1},

		// Scattered pieces around the main cluster
		{centerX - 2, centerY},
		{centerX + 2, centerY},
		{centerX, centerY - 2},
		{centerX, centerY + 2},
		{centerX - 1, centerY + 2},
		{centerX + 1, centerY - 2},

		// Outer scattered debris
		{centerX - 3, centerY + 1},
		{centerX + 3, centerY - 1},
		{centerX - 2, centerY - 2},
		{centerX + 2, centerY + 2},

		// Small fragments
		{centerX - 4, centerY},
		{centerX + 4, centerY},
		{centerX, centerY - 3},
		{centerX, centerY + 3},
		{centerX - 3, centerY - 1},
		{centerX + 3, centerY + 1},
	}

	// Draw the debris pieces
	for _, pos := range debrisPositions {
		if pos.X >= 0 && pos.X < size && pos.Y >= 0 && pos.Y < size {
			img.Set(pos.X, pos.Y, color.Black)
		}
	}

	// Add some additional random scattered bits to make it look more chaotic
	additionalDebris := []Position{
		{centerX - 1, centerY + 3},
		{centerX + 1, centerY - 3},
		{centerX - 4, centerY + 1},
		{centerX + 4, centerY - 1},
		{centerX - 2, centerY + 3},
		{centerX + 2, centerY - 3},
	}

	for _, pos := range additionalDebris {
		if pos.X >= 0 && pos.X < size && pos.Y >= 0 && pos.Y < size {
			img.Set(pos.X, pos.Y, color.Black)
		}
	}

	return img
}
