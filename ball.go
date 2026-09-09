package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ballFriction     = 0.985
	ballRestitution  = 0.75
	ballMaxSpeed     = 14.0
	ballCarry        = 1.15
	ballBounceBoost  = 2.2
	ballSpinPerSpeed = 0.08
)

type Ball struct {
	Pos    Vec
	Vel    Vec
	Angle  float64
	Radius float64
	Frames []*ebiten.Image
	Home   Vec
}

func newBall(frames []*ebiten.Image, home Vec) *Ball {
	native := float64(frames[0].Bounds().Dx())
	return &Ball{
		Pos:    home,
		Radius: native * spriteScale / 2,
		Frames: frames,
		Home:   home,
	}
}

func (b *Ball) Reset() {
	b.Pos = b.Home
	b.Vel = Vec{}
	b.Angle = 0
}

func (b *Ball) Update(pitch *Pitch) {
	b.Pos = b.Pos.Add(b.Vel)
	b.Vel = b.Vel.Mul(ballFriction)
	if spd := b.Vel.Len(); spd > ballMaxSpeed {
		b.Vel = b.Vel.Mul(ballMaxSpeed / spd)
	}
	if spd := b.Vel.Len(); spd > 0.15 {
		b.Angle += spd * ballSpinPerSpeed
	} else if spd < 0.05 {
		b.Vel = Vec{}
	}
	b.bounceWalls(pitch)
}

func (b *Ball) bounceWalls(pitch *Pitch) {
	r := b.Radius
	top := float64(pitch.Y)
	bot := float64(pitch.Y + pitch.H)
	left := float64(pitch.X)
	right := float64(pitch.X + pitch.W)
	gt := float64(pitch.GoalTop())
	gb := gt + float64(pitch.GoalH)
	gw := float64(pitch.GoalW)

	if b.Pos.Y-r < top {
		b.Pos.Y = top + r
		b.Vel.Y = -b.Vel.Y * ballRestitution
	}
	if b.Pos.Y+r > bot {
		b.Pos.Y = bot - r
		b.Vel.Y = -b.Vel.Y * ballRestitution
	}

	inLeftMouth := b.Pos.Y > gt+r && b.Pos.Y < gb-r
	inRightMouth := b.Pos.Y > gt+r && b.Pos.Y < gb-r

	if b.Pos.X-r < left {
		if inLeftMouth {
			if b.Pos.X-r < left-gw {
				b.Pos.X = left - gw + r
				b.Vel.X = -b.Vel.X * ballRestitution
			}
			if b.Pos.Y-r < gt {
				b.Pos.Y = gt + r
				b.Vel.Y = -b.Vel.Y * ballRestitution
			}
			if b.Pos.Y+r > gb {
				b.Pos.Y = gb - r
				b.Vel.Y = -b.Vel.Y * ballRestitution
			}
		} else {
			b.Pos.X = left + r
			b.Vel.X = -b.Vel.X * ballRestitution
		}
	}
	if b.Pos.X+r > right {
		if inRightMouth {
			if b.Pos.X+r > right+gw {
				b.Pos.X = right + gw - r
				b.Vel.X = -b.Vel.X * ballRestitution
			}
			if b.Pos.Y-r < gt {
				b.Pos.Y = gt + r
				b.Vel.Y = -b.Vel.Y * ballRestitution
			}
			if b.Pos.Y+r > gb {
				b.Pos.Y = gb - r
				b.Vel.Y = -b.Vel.Y * ballRestitution
			}
		} else {
			b.Pos.X = right - r
			b.Vel.X = -b.Vel.X * ballRestitution
		}
	}
}

func (b *Ball) CollidePlayer(p *Player) {
	delta := b.Pos.Sub(p.Pos)
	dist := delta.Len()
	minDist := b.Radius + p.Radius
	if dist == 0 {
		delta = p.facingDir()
		if delta.Len() == 0 {
			delta = Vec{1, 0}
		}
		dist = 0.001
	}
	if dist >= minDist {
		return
	}
	n := delta.Mul(1 / dist)
	b.Pos = p.Pos.Add(n.Mul(minDist))

	rel := b.Vel.Sub(p.Vel)
	closing := rel.Dot(n)
	if closing < 0 {
		b.Vel = b.Vel.Sub(n.Mul(closing * (1 + ballRestitution)))
		b.Vel = b.Vel.Add(p.Vel.Mul(ballCarry))
		b.Vel = b.Vel.Add(n.Mul(ballBounceBoost))
	}
	if p.kicking() {
		b.Vel = b.Vel.Add(p.facingDir().Mul(kickImpulse))
	}
	if spd := b.Vel.Len(); spd > ballMaxSpeed {
		b.Vel = b.Vel.Mul(ballMaxSpeed / spd)
	}
}

func (b *Ball) Draw(dst *ebiten.Image) {
	n := len(b.Frames)
	idx := 0
	if n > 0 {
		idx = int(math.Abs(b.Angle)/0.6) % n
	}
	drawSprite(dst, b.Frames[idx], b.Pos, spriteScale, b.Angle, 0)
}
