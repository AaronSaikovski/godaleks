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

import "image/color"

// Pre-allocated color values to avoid per-frame allocations
var (
	colorOverlay = color.RGBA{0, 0, 0, 128}
	colorRed     = color.RGBA{255, 0, 0, 255}
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridWidth    = 50
	gridHeight   = 33 // Reduced to ensure grid fits within screen (50 + 33*16 = 578 pixels)
	cellSize     = 16

	// Pre-computed layout offsets (derived from screen/grid constants)
	gridOffsetX = (screenWidth - gridWidth*cellSize) / 2
	gridOffsetY = 50
)

const (
	StateMenu GameState = iota
	StatePlaying
	StateLevelComplete
	StateGameOver
	StateWin
)

const levelTransitionDelay = 1.5 // seconds to wait between levels

const trigTableSize = 32 // Number of pre-computed sin/cos entries
