package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 1280
	screenH = 720
)

var lineWhite = color.RGBA{R: 245, G: 245, B: 245, A: 230}

type Pitch struct {
	X, Y, W, H float32
	GoalW      float32
	GoalH      float32
	LineW      float32
	Image      *ebiten.Image
}

func newPitch(assets *Assets) *Pitch {
	p := &Pitch{
		X:     96,
		Y:     56,
		W:     1088,
		H:     608,
		GoalW: 64,
		GoalH: 220,
		LineW: 4,
	}
	p.Image = p.render(assets)
	return p
}

func (p *Pitch) Center() Vec {
	return Vec{float64(p.X + p.W/2), float64(p.Y + p.H/2)}
}

func (p *Pitch) GoalTop() float32 {
	return p.Y + (p.H-p.GoalH)/2
}

func (p *Pitch) LeftGoal() (x, y, w, h float64) {
	return float64(p.X - p.GoalW), float64(p.GoalTop()), float64(p.GoalW), float64(p.GoalH)
}

func (p *Pitch) RightGoal() (x, y, w, h float64) {
	return float64(p.X + p.W), float64(p.GoalTop()), float64(p.GoalW), float64(p.GoalH)
}

func inRect(pos Vec, x, y, w, h float64) bool {
	return pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h
}

func (p *Pitch) render(assets *Assets) *ebiten.Image {
	img := ebiten.NewImage(screenW, screenH)
	tileGrass(img, assets.Grass)
	p.stampNets(img, assets.Net)
	p.drawMarkings(img)
	p.stampFlags(img, assets.Flag)
	return img
}

func tileGrass(dst, grass *ebiten.Image) {
	gw := grass.Bounds().Dx()
	gh := grass.Bounds().Dy()
	for y := 0; y < screenH; y += gh {
		for x := 0; x < screenW; x += gw {
			op := &ebiten.DrawImageOptions{}
			op.Filter = ebiten.FilterNearest
			op.GeoM.Translate(float64(x), float64(y))
			dst.DrawImage(grass, op)
		}
	}
}

func (p *Pitch) stampNets(dst, net *ebiten.Image) {
	nw := net.Bounds().Dx()
	nh := net.Bounds().Dy()
	stamp := func(gx, gy, gw, gh float32) {
		for y := gy; y < gy+gh; y += float32(nh) {
			for x := gx; x < gx+gw; x += float32(nw) {
				op := &ebiten.DrawImageOptions{}
				op.Filter = ebiten.FilterNearest
				op.GeoM.Translate(float64(x), float64(y))
				dst.DrawImage(net, op)
			}
		}
	}
	gt := p.GoalTop()
	stamp(p.X-p.GoalW, gt, p.GoalW, p.GoalH)
	stamp(p.X+p.W, gt, p.GoalW, p.GoalH)
}

func (p *Pitch) stampFlags(dst, flag *ebiten.Image) {
	const scale = 2.0
	w := float64(flag.Bounds().Dx()) * scale
	h := float64(flag.Bounds().Dy()) * scale
	corners := []Vec{
		{float64(p.X), float64(p.Y)},
		{float64(p.X + p.W), float64(p.Y)},
		{float64(p.X), float64(p.Y + p.H)},
		{float64(p.X + p.W), float64(p.Y + p.H)},
	}
	for _, c := range corners {
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(c.X-w/2, c.Y-h)
		dst.DrawImage(flag, op)
	}
}

func (p *Pitch) drawMarkings(dst *ebiten.Image) {
	lw := p.LineW
	aa := true

	vector.StrokeRect(dst, p.X, p.Y, p.W, p.H, lw, lineWhite, aa)

	midX := p.X + p.W/2
	midY := p.Y + p.H/2
	vector.StrokeLine(dst, midX, p.Y, midX, p.Y+p.H, lw, lineWhite, aa)

	centerR := float32(88)
	vector.StrokeCircle(dst, midX, midY, centerR, lw, lineWhite, aa)
	vector.FillCircle(dst, midX, midY, 5, lineWhite, aa)

	penW := p.W * 0.16
	penH := p.H * 0.58
	penY := p.Y + (p.H-penH)/2
	vector.StrokeRect(dst, p.X, penY, penW, penH, lw, lineWhite, aa)
	vector.StrokeRect(dst, p.X+p.W-penW, penY, penW, penH, lw, lineWhite, aa)

	goalAreaW := p.W * 0.055
	goalAreaH := p.H * 0.28
	goalAreaY := p.Y + (p.H-goalAreaH)/2
	vector.StrokeRect(dst, p.X, goalAreaY, goalAreaW, goalAreaH, lw, lineWhite, aa)
	vector.StrokeRect(dst, p.X+p.W-goalAreaW, goalAreaY, goalAreaW, goalAreaH, lw, lineWhite, aa)

	spot := p.W * 0.11
	vector.FillCircle(dst, p.X+spot, midY, 4, lineWhite, aa)
	vector.FillCircle(dst, p.X+p.W-spot, midY, 4, lineWhite, aa)

	p.strokeArc(dst, p.X+penW, midY, 56, float32(-math.Pi/2), float32(math.Pi/2))
	p.strokeArc(dst, p.X+p.W-penW, midY, 56, float32(math.Pi/2), float32(3*math.Pi/2))

	cr := float32(26)
	p.strokeArc(dst, p.X, p.Y, cr, 0, float32(math.Pi/2))
	p.strokeArc(dst, p.X+p.W, p.Y, cr, float32(math.Pi/2), float32(math.Pi))
	p.strokeArc(dst, p.X+p.W, p.Y+p.H, cr, float32(math.Pi), float32(3*math.Pi/2))
	p.strokeArc(dst, p.X, p.Y+p.H, cr, float32(3*math.Pi/2), float32(2*math.Pi))

	gt := p.GoalTop()
	gb := gt + p.GoalH
	// Left goal frame (open toward the pitch).
	vector.StrokeLine(dst, p.X-p.GoalW, gt, p.X, gt, lw, lineWhite, aa)
	vector.StrokeLine(dst, p.X-p.GoalW, gb, p.X, gb, lw, lineWhite, aa)
	vector.StrokeLine(dst, p.X-p.GoalW, gt, p.X-p.GoalW, gb, lw, lineWhite, aa)
	// Right goal frame.
	vector.StrokeLine(dst, p.X+p.W, gt, p.X+p.W+p.GoalW, gt, lw, lineWhite, aa)
	vector.StrokeLine(dst, p.X+p.W, gb, p.X+p.W+p.GoalW, gb, lw, lineWhite, aa)
	vector.StrokeLine(dst, p.X+p.W+p.GoalW, gt, p.X+p.W+p.GoalW, gb, lw, lineWhite, aa)
}

func (p *Pitch) strokeArc(dst *ebiten.Image, cx, cy, r, start, end float32) {
	path := &vector.Path{}
	path.Arc(cx, cy, r, start, end, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: p.LineW, LineJoin: vector.LineJoinRound}
	drawOp := &vector.DrawPathOptions{}
	drawOp.AntiAlias = true
	drawOp.ColorScale.ScaleWithColor(lineWhite)
	vector.StrokePath(dst, path, strokeOp, drawOp)
}
