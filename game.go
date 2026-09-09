package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	resetDelay = 50
	scoreToWin = 5
)

type Game struct {
	pitch  *Pitch
	p1     *Player
	p2     *Player
	ball   *Ball
	score1 int
	score2 int
	freeze int
	winner int // 0 none, 1 blue, 2 red
}

func NewGame() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}
	pitch := newPitch(assets)
	center := pitch.Center()
	half := float64(pitch.X + pitch.W/2)
	r := playerRadius

	p1 := newPlayer(
		assets.Blue[:2],
		Vec{float64(pitch.X) + 160, center.Y},
		float64(pitch.X)+r,
		half-r,
		0,
		ebiten.KeyW, ebiten.KeyS, ebiten.KeyA, ebiten.KeyD, ebiten.KeyShiftLeft,
	)
	p2 := newPlayer(
		assets.Red[:2],
		Vec{float64(pitch.X+pitch.W) - 160, center.Y},
		half+r,
		float64(pitch.X+pitch.W)-r,
		math.Pi,
		ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyShiftRight,
	)
	ball := newBall(assets.Ball, center)
	g := &Game{pitch: pitch, p1: p1, p2: p2, ball: ball}
	g.kickoff()
	return g, nil
}

func (g *Game) kickoff() {
	g.p1.Reset()
	g.p2.Reset()
	g.ball.Reset()
	g.p1.Facing = 0
	g.p2.Facing = math.Pi
}

func (g *Game) restartMatch() {
	g.score1 = 0
	g.score2 = 0
	g.winner = 0
	g.freeze = 0
	g.kickoff()
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.restartMatch()
		return nil
	}
	if g.winner != 0 {
		return nil
	}
	if g.freeze > 0 {
		g.freeze--
		if g.freeze == 0 {
			g.kickoff()
		}
		return nil
	}

	g.p1.Update(g.pitch)
	g.p2.Update(g.pitch)
	separatePlayers(g.p1, g.p2)
	g.p1.clamp(g.pitch)
	g.p2.clamp(g.pitch)

	g.ball.Update(g.pitch)
	g.ball.CollidePlayer(g.p1)
	g.ball.CollidePlayer(g.p2)

	g.checkGoal()
	return nil
}

func (g *Game) checkGoal() {
	lx, ly, lw, lh := g.pitch.LeftGoal()
	rx, ry, rw, rh := g.pitch.RightGoal()
	if inRect(g.ball.Pos, lx, ly, lw, lh) {
		g.score2++
		g.afterGoal()
		return
	}
	if inRect(g.ball.Pos, rx, ry, rw, rh) {
		g.score1++
		g.afterGoal()
	}
}

func (g *Game) afterGoal() {
	if g.score1 >= scoreToWin {
		g.winner = 1
		return
	}
	if g.score2 >= scoreToWin {
		g.winner = 2
		return
	}
	g.freeze = resetDelay
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.pitch.Image, nil)
	g.p1.Draw(screen)
	g.p2.Draw(screen)
	g.ball.Draw(screen)
	g.drawHUD(screen)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, screenW, 28, color.RGBA{0, 0, 0, 130}, false)
	vector.FillRect(screen, 0, screenH-24, screenW, 24, color.RGBA{0, 0, 0, 130}, false)
	msg := fmt.Sprintf("%d   -   %d", g.score1, g.score2)
	ebitenutil.DebugPrintAt(screen, msg, screenW/2-28, 8)
	ebitenutil.DebugPrintAt(screen, "P1 WASD + LShift     first to 5     P2 Arrows + RShift     R restart", 280, screenH-18)

	if g.winner != 0 {
		vector.FillRect(screen, 0, 0, screenW, screenH, color.RGBA{0, 0, 0, 140}, false)
		who := "BLUE WINS"
		if g.winner == 2 {
			who = "RED WINS"
		}
		ebitenutil.DebugPrintAt(screen, who+"  —  press R to play again", screenW/2-110, screenH/2-8)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}
