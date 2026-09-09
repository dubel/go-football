package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var assetFS embed.FS

const spriteScale = 4

type Assets struct {
	Blue  []*ebiten.Image
	Red   []*ebiten.Image
	Ball  []*ebiten.Image
	Grass *ebiten.Image
	Net   *ebiten.Image
	Flag  *ebiten.Image
}

func loadAssets() (*Assets, error) {
	a := &Assets{}
	var err error
	a.Blue, err = loadSeries("assets/players/blue/characterBlue_%d.png", 1, 14)
	if err != nil {
		return nil, err
	}
	a.Red, err = loadSeries("assets/players/red/characterRed_%d.png", 1, 14)
	if err != nil {
		return nil, err
	}
	a.Ball, err = loadSeries("assets/ball/ball_soccer%d.png", 1, 4)
	if err != nil {
		return nil, err
	}
	a.Grass, err = loadImage("assets/pitch/grass.png")
	if err != nil {
		return nil, err
	}
	a.Net, err = loadImage("assets/pitch/net.png")
	if err != nil {
		return nil, err
	}
	a.Flag, err = loadImage("assets/pitch/flag.png")
	if err != nil {
		return nil, err
	}
	return a, nil
}

func loadSeries(pattern string, from, to int) ([]*ebiten.Image, error) {
	out := make([]*ebiten.Image, 0, to-from+1)
	for i := from; i <= to; i++ {
		img, err := loadImage(fmt.Sprintf(pattern, i))
		if err != nil {
			return nil, err
		}
		out = append(out, img)
	}
	return out, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := assetFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return ebiten.NewImageFromImage(img), nil
}

func drawSprite(dst, src *ebiten.Image, pos Vec, scale, angle, bob float64) {
	w := float64(src.Bounds().Dx())
	h := float64(src.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(pos.X, pos.Y+bob)
	dst.DrawImage(src, op)
}
