package main

import "math"

type Vec struct {
	X, Y float64
}

func (v Vec) Add(o Vec) Vec     { return Vec{v.X + o.X, v.Y + o.Y} }
func (v Vec) Sub(o Vec) Vec     { return Vec{v.X - o.X, v.Y - o.Y} }
func (v Vec) Mul(s float64) Vec { return Vec{v.X * s, v.Y * s} }
func (v Vec) Dot(o Vec) float64 { return v.X*o.X + v.Y*o.Y }
func (v Vec) Len2() float64     { return v.X*v.X + v.Y*v.Y }
func (v Vec) Len() float64      { return math.Hypot(v.X, v.Y) }
func (v Vec) Angle() float64    { return math.Atan2(v.Y, v.X) }

func (v Vec) Normalize() Vec {
	l := v.Len()
	if l == 0 {
		return Vec{}
	}
	return v.Mul(1 / l)
}

func clamp(v, lo, hi float64) float64 {
	return min(hi, max(lo, v))
}

func snap8(angle float64) float64 {
	const step = math.Pi / 4
	return math.Round(angle/step) * step
}
