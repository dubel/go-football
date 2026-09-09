package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	playerAccel    = 0.55
	playerMaxSpeed = 5.2
	playerDamp     = 0.82
	playerRadius   = 36.0
	kickImpulse    = 9.5
)

type Player struct {
	Pos    Vec
	Vel    Vec
	Facing float64
	Radius float64
	Frames []*ebiten.Image
	Home   Vec
	MinX   float64
	MaxX   float64
	Anim   int
	Kick   ebiten.Key
	Up     ebiten.Key
	Down   ebiten.Key
	Left   ebiten.Key
	Right  ebiten.Key
}

func newPlayer(frames []*ebiten.Image, home Vec, minX, maxX float64, facing float64, up, down, left, right, kick ebiten.Key) *Player {
	return &Player{
		Pos:    home,
		Facing: facing,
		Radius: playerRadius,
		Frames: frames,
		Home:   home,
		MinX:   minX,
		MaxX:   maxX,
		Kick:   kick,
		Up:     up,
		Down:   down,
		Left:   left,
		Right:  right,
	}
}

func (p *Player) Reset() {
	p.Pos = p.Home
	p.Vel = Vec{}
	p.Anim = 0
}

func (p *Player) kicking() bool {
	return inpututil.IsKeyJustPressed(p.Kick)
}

func (p *Player) Update(pitch *Pitch) {
	in := Vec{}
	if ebiten.IsKeyPressed(p.Up) {
		in.Y -= 1
	}
	if ebiten.IsKeyPressed(p.Down) {
		in.Y += 1
	}
	if ebiten.IsKeyPressed(p.Left) {
		in.X -= 1
	}
	if ebiten.IsKeyPressed(p.Right) {
		in.X += 1
	}
	moving := in.Len() > 0
	if moving {
		in = in.Normalize()
		p.Vel = p.Vel.Add(in.Mul(playerAccel))
		p.Facing = snap8(in.Angle())
		p.Anim++
	} else {
		p.Vel = p.Vel.Mul(playerDamp)
		p.Anim = 0
	}
	if spd := p.Vel.Len(); spd > playerMaxSpeed {
		p.Vel = p.Vel.Mul(playerMaxSpeed / spd)
	}
	p.Pos = p.Pos.Add(p.Vel)
	p.clamp(pitch)
}

func (p *Player) clamp(pitch *Pitch) {
	r := p.Radius
	p.Pos.X = clamp(p.Pos.X, p.MinX, p.MaxX)
	p.Pos.Y = clamp(p.Pos.Y, float64(pitch.Y)+r, float64(pitch.Y+pitch.H)-r)
}

func (p *Player) facingDir() Vec {
	return Vec{math.Cos(p.Facing), math.Sin(p.Facing)}
}

func (p *Player) Draw(dst *ebiten.Image) {
	frame := p.Frames[0]
	if p.Vel.Len() > 0.4 && len(p.Frames) > 1 {
		frame = p.Frames[(p.Anim/8)%len(p.Frames)]
	}
	bob := 0.0
	if p.Vel.Len() > 0.4 {
		bob = math.Sin(float64(p.Anim)*0.45) * 1.5
	}
	// Kenney standing poses face east; snap Facing to 8 directions and rotate.
	drawSprite(dst, frame, p.Pos, spriteScale, p.Facing, bob)
}

func separatePlayers(a, b *Player) {
	delta := b.Pos.Sub(a.Pos)
	dist := delta.Len()
	minDist := a.Radius + b.Radius
	if dist == 0 {
		delta = Vec{1, 0}
		dist = 1
	}
	if dist >= minDist {
		return
	}
	n := delta.Mul(1 / dist)
	overlap := (minDist - dist) / 2
	a.Pos = a.Pos.Sub(n.Mul(overlap))
	b.Pos = b.Pos.Add(n.Mul(overlap))
}
