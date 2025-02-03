package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	red  = color.RGBA{255, 0, 0, 255}
	blue = color.RGBA{0, 0, 255, 255}
)

type Game struct {
	Width        float32
	Height       float32
	Data         chan [][2]float32
	current      [][2]float32
	GenerationNb int
	Distance     float64
}

func NewGame(width, height float32) (*Game, func()) {
	game := &Game{
		Width:  width,
		Height: height,
		Data:   make(chan [][2]float32, 200),
	}
	return game, func() { close(game.Data) }
}

func (g *Game) AddData(data [][2]float32) {
	g.Data <- data
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("escape key pressed")
	}
	return nil
}

func (g *Game) Draw(img *ebiten.Image) {
	// Save current generation
	data := <-g.Data
	if len(data) > 0 {
		g.current = data
	}

	// Draw background (white)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw lines
	n := len(g.current)
	for i := range n {
		var x1, y1 float32
		if i == 0 {
			x1 = g.current[n-1][0]
			y1 = g.current[n-1][1]
		} else {
			x1 = g.current[i-1][0]
			y1 = g.current[i-1][1]
		}
		x0 := g.current[i][0]
		y0 := g.current[i][1]
		vector.StrokeLine(img, x0, y0, x1, y1, 1, blue, true)
	}

	// Draw circles
	for _, row := range g.current {
		cx := row[0]
		cy := row[1]
		vector.DrawFilledCircle(img, cx, cy, 3, red, true)
	}

	// Draw text
	//	ft := basicfont.Face7x13
	ebitenutil.DebugPrintAt(img, fmt.Sprintf("Generation #%d", g.GenerationNb), 10, 560)
	ebitenutil.DebugPrintAt(img, fmt.Sprintf("Distance: %.2f", g.Distance), 10, 580)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.Width), int(g.Height)
}

func Run(game *Game) error {
	// Initialize the Ebiten game loop
	ebiten.SetWindowSize(int(game.Height), int(game.Width))
	ebiten.SetWindowTitle("example - traveling salesman")
	if err := ebiten.RunGame(game); err != nil {
		return err
	}
	return nil
}
