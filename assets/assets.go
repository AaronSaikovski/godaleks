package assets

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var assets embed.FS

var HeroImage = mustLoadImage("hero.png")
var Dalek = mustLoadImage("robot.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// func mustLoadImages(path string) []*ebiten.Image {
// 	matches, err := fs.Glob(assets, path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	images := make([]*ebiten.Image, len(matches))
// 	for i, match := range matches {
// 		images[i] = mustLoadImage(match)
// 	}

// 	return images
// }

// func mustLoadFont(name string) font.Face {
// 	f, err := assets.ReadFile(name)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tt, err := opentype.Parse(f)
// 	if err != nil {
// 		panic(err)
// 	}

// 	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
// 		Size:    48,
// 		DPI:     72,
// 		Hinting: font.HintingVertical,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return face
// }
