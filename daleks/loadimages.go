package daleks

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
