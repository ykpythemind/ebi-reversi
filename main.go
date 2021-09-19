package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	img *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Fill")

	img, _, err := ebitenutil.NewImageFromFile("images/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{img: img}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
