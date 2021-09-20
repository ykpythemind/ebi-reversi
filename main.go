package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	img   *ebiten.Image
	Board *Board
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	drawX := 0
	drawY := 0

	b := g.img.Bounds()
	imgScaleX, imgScaleY := getScaleToAdjustRect(&b, 40, 40)

	for i := 0; i < 8; i++ {
		drawX = i * 40
		for j := 0; j < 8; j++ {
			drawY = j * 40
			ebitenutil.DrawRect(screen, float64(drawX), float64(drawY), float64(40), float64(40), color.Black)
			ebitenutil.DrawRect(screen, float64(drawX+1), float64(drawY+1), float64(40-2), float64(40-2), color.White)

			if g.Board[i][j].has {
				geom := ebiten.GeoM{}
				geom.Scale(imgScaleX, imgScaleY)
				geom.Translate(float64(drawX), float64(drawY))
				screen.DrawImage(g.img, &ebiten.DrawImageOptions{GeoM: geom})
			}
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

type Square struct {
	has bool // tmp
}

// Board is 8x8 reversi board
type Board [8][8]Square

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Fill")

	img, _, err := ebitenutil.NewImageFromFile("images/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	img.Bounds()

	board := &Board{}
	board[1][2].has = true
	board[6][7].has = true
	game := &Game{img: img, Board: board}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getScaleToAdjustRect(rect *image.Rectangle, cx, cy float64) (sx, sy float64) {
	x := float64(rect.Max.X - rect.Min.X)
	y := float64(rect.Max.Y - rect.Min.Y)

	return cx / x, cy / y
}
