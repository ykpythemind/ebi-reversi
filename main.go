package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/gofont/goitalic"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var normalFont font.Face

const (
	BOARDX, BOARDY = 320, 320
	SQUARE         = BOARDX / 8
)

func init() {

	tt, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	img           *ebiten.Image
	Board         *Board
	CurrentSquare *Square // CurrentSquare is hovered square
}

func (g *Game) Update() error {
	// fmt.Println(ebiten.CursorPosition())
	err, ix, iy := getCurrentSquare(ebiten.CursorPosition())
	if err == nil {
		// square detected
		g.CurrentSquare = g.Board[ix][iy]
	} else {
		g.CurrentSquare = nil
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	drawX := 0
	drawY := 0

	b := g.img.Bounds()
	imgScaleX, imgScaleY := getScaleToAdjustRect(&b, SQUARE, SQUARE)

	for i := 0; i < 8; i++ {
		drawX = i * SQUARE
		for j := 0; j < 8; j++ {
			drawY = j * SQUARE
			ebitenutil.DrawRect(screen, float64(drawX), float64(drawY), float64(SQUARE), float64(SQUARE), color.Black)
			ebitenutil.DrawRect(screen, float64(drawX+1), float64(drawY+1), float64(SQUARE-2), float64(SQUARE-2), color.White)

			if cr := g.CurrentSquare; cr != nil {
				if cr.pos.X == i && cr.pos.Y == j {
					// selected square
					col := color.RGBA{85, 165, 34, 255}
					ebitenutil.DrawRect(screen, float64(drawX+1), float64(drawY+1), float64(SQUARE-2), float64(SQUARE-2), col)
				}
			}

			sq := g.Board[i][j]
			if sq.has {
				geom := ebiten.GeoM{}
				geom.Scale(imgScaleX, imgScaleY)
				geom.Translate(float64(drawX), float64(drawY))
				screen.DrawImage(g.img, &ebiten.DrawImageOptions{GeoM: geom})
			}
		}
	}

	// debug message
	msg := fmt.Sprintf("TPS: %0.2f / FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, normalFont, 50, 50, color.Black)
	if g.CurrentSquare != nil {
		sq := g.CurrentSquare
		sqMsg := fmt.Sprintf("square: (%d,%d)", sq.pos.X, sq.pos.Y)
		text.Draw(screen, sqMsg, normalFont, 20, 20, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

type SquarePosition struct {
	X int
	Y int
}

type Square struct {
	has bool // tmp
	pos SquarePosition
}

// Board is 8x8 reversi board
type Board [8][8]*Square

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Fill")

	img, _, err := ebitenutil.NewImageFromFile("images/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	img.Bounds()

	// init board
	board := &Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			s := &Square{pos: SquarePosition{X: i, Y: j}}
			board[i][j] = s
		}
	}

	board[1][2].has = true
	board[6][7].has = true
	game := &Game{img: img, Board: board}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getScaleToAdjustRect(rect *image.Rectangle, cx, cy float64) (sx, sy float64) {
	return cx / float64(rect.Dx()), cy / float64(rect.Dy())
}

func getCurrentSquare(curX, curY int) (err error, ix, iy int) {
	if curX < 0 || BOARDY < curY {
		return errors.New("out of board"), 0, 0
	}

	point := image.Point{curX, curY}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			min := image.Point{i * SQUARE, j * SQUARE}
			max := image.Point{i*SQUARE + SQUARE - 1, j*SQUARE + SQUARE - 1}
			rect := image.Rectangle{Min: min, Max: max}
			if point.In(rect) {
				return nil, i, j
			}
		}
	}

	return errors.New("not found"), 0, 0
}
