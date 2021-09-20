package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
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
	imgPlayerA    playerImage
	imgPlayerB    playerImage
	Board         *Board
	CurrentSquare *Square // CurrentSquare is hovered square
	Player        Player
	canPlace      bool
}

func (g *Game) Update() error {
	// fmt.Println(ebiten.CursorPosition())
	err, ix, iy := getCurrentSquare(ebiten.CursorPosition())

	if err == nil {
		// square detected
		sq := g.Board[ix][iy]
		g.CurrentSquare = sq
	} else {
		g.CurrentSquare = nil
	}

	clicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// フォーカス中マスのチェック
	if sq := g.CurrentSquare; sq != nil {
		squares := g.Board.Check(sq, g.Player)

		if len(squares) > 0 {
			// you can place
			g.canPlace = true
			if clicked {
				// clicked. update state
				err := g.Board.Place(sq, g.Player)
				if err != nil {
					// - skipが必要な場合
					// - 試合が終わった場合
					if errors.Is(err, GameEndError) {
						// todo result and reset
						fmt.Printf("game end!\n")
						return nil
					}
					if errors.Is(err, NeedPassError) {
						fmt.Printf("player %s skipped!\n", g.Player)
						return nil
					}

					return err
				}
				g.switchTurn()
			}
		} else {
			g.canPlace = false
			// ignore
		}
	}

	return nil
}

func (g *Game) switchTurn() {
	if g.Player == PlayerA {
		g.Player = PlayerB
	} else {
		g.Player = PlayerA
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sq := g.Board[i][j]
			sq.Draw(g, screen)
		}
	}

	// draw prompt
	var st string
	if g.Player == PlayerA {
		st = "player A turn"
	} else if g.Player == PlayerB {
		st = "player B turn"
	}
	if st != "" {
		text.Draw(screen, st, normalFont, BOARDX-100, BOARDY-10, color.Black)
	}

	// debug message
	msg := fmt.Sprintf("TPS: %0.2f / FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, normalFont, 50, 50, color.Black)
	if g.CurrentSquare != nil {
		sq := g.CurrentSquare
		sqMsg := fmt.Sprintf("square: (%d,%d)", sq.pos.X, sq.pos.Y)
		text.Draw(screen, sqMsg, normalFont, 20, 20, color.Black)

		if g.canPlace {
			text.Draw(screen, "can place", normalFont, 20, 40, color.Black)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

type playerImage struct {
	image  *ebiten.Image
	scaleX float64
	scaleY float64
}

type GameState int

type Player int

func (p Player) String() string {
	switch p {
	case PlayerA:
		return "A"
	case PlayerB:
		return "B"
	default:
		return "unknown"
	}
}

const (
	PlayerA Player = iota
	PlayerB
)

type SquarePosition struct {
	X int
	Y int
}

type SquareState int

const (
	Blank SquareState = iota
	PlayerAFilled
	PlayerBFilled
)

type Square struct {
	pos   SquarePosition
	state SquareState
}

func (s *Square) IsBlank() bool {
	return s.state == Blank
}

func (s *Square) Eq(sq *Square) bool {
	return s.pos.X == sq.pos.X && s.pos.Y == sq.pos.Y
}

func (s *Square) ScreenPos() (x, y float64) {
	return float64(s.pos.X * SQUARE), float64(s.pos.Y * SQUARE)
}

func (s *Square) Draw(g *Game, screen *ebiten.Image) {
	drawX, drawY := s.ScreenPos()
	ebitenutil.DrawRect(screen, float64(drawX), float64(drawY), float64(SQUARE), float64(SQUARE), color.Black)
	ebitenutil.DrawRect(screen, float64(drawX+1), float64(drawY+1), float64(SQUARE-2), float64(SQUARE-2), color.White)

	if cr := g.CurrentSquare; cr != nil {
		if s.Eq(cr) {
			// selected square
			col := color.RGBA{85, 165, 34, 255}
			ebitenutil.DrawRect(screen, float64(drawX+1), float64(drawY+1), float64(SQUARE-2), float64(SQUARE-2), col)
		}
	}

	if s.state != Blank {
		geom := ebiten.GeoM{}
		var player *playerImage
		switch s.state {
		case PlayerAFilled:
			player = &g.imgPlayerA
		case PlayerBFilled:
			player = &g.imgPlayerB
		default:
			panic("not reachable")
		}
		geom.Scale(player.scaleX, player.scaleY)
		geom.Translate(s.ScreenPos())
		screen.DrawImage(player.image, &ebiten.DrawImageOptions{GeoM: geom})
	}
}

// NewGame initialize game initial state
func NewGame() *Game {
	// init board
	board := &Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			s := &Square{pos: SquarePosition{X: i, Y: j}}
			board[i][j] = s
		}
	}

	// game init
	board[3][3].state = PlayerAFilled
	board[3][4].state = PlayerBFilled
	board[4][3].state = PlayerBFilled
	board[4][4].state = PlayerAFilled

	return &Game{Board: board}
}

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Fill")

	timgA, _, err := ebitenutil.NewImageFromFile("images/gopher.png")
	if err != nil {
		log.Fatal(err)
	}
	a := timgA.Bounds()
	imgScaleX, imgScaleY := getScaleToAdjustRect(&a, SQUARE, SQUARE)
	imgPlayerA := playerImage{image: timgA, scaleX: imgScaleX, scaleY: imgScaleY}

	timgB, _, err := ebitenutil.NewImageFromFile("images/ykpythemind.jpg")
	if err != nil {
		log.Fatal(err)
	}
	b := timgB.Bounds()
	imgScaleX, imgScaleY = getScaleToAdjustRect(&b, SQUARE, SQUARE)
	imgPlayerB := playerImage{image: timgB, scaleX: imgScaleX, scaleY: imgScaleY}

	game := NewGame()
	game.imgPlayerA = imgPlayerA
	game.imgPlayerB = imgPlayerB

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
