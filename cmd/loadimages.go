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
	"bytes"
	"embed"
	"image"
	_ "image/png" // Import image format support

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

// Daleks Images
type DalekGameImages struct {
	Human *ebiten.Image
	Dalek *ebiten.Image
}

// loadImage loads an image from the assets directory
func loadImage(filename string) (*ebiten.Image, error) {
	data, err := assets.ReadFile("assets/" + filename)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// / LoadImages initializes all crossword images and returns an error if any image fails to load
func (images *DalekGameImages) LoadImages() error {
	var err error

	if images.Human, err = loadImage("human.png"); err != nil {
		return err
	}
	if images.Dalek, err = loadImage("dalek.png"); err != nil {
		return err
	}

	return nil
}

// Loads images
func loadImages() *DalekGameImages {
	gameImages := &DalekGameImages{}
	gameImages.LoadImages()
	return gameImages
}
