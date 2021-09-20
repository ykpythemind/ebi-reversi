package main

import (
	"bytes"
	"errors"
	"fmt"
)

// Board is 8x8 reversi board
type Board [8][8]*Square

func (b *Board) String() string {
	var buf bytes.Buffer
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sq := b[i][j]
			switch sq.state {
			case PlayerAFilled:
				buf.Write([]byte("x"))
			case PlayerBFilled:
				buf.Write([]byte("o"))
			case Blank:
				buf.Write([]byte("-"))
			default:
				panic("unreachable")
			}
			if j != 7 {
				buf.Write([]byte(","))
			}
		}
		buf.Write([]byte("\n"))
	}

	return buf.String()
}

// Check checks whether player can place given square to the current board
func (b *Board) Check(checkTarget *Square) error {
	if !checkTarget.IsBlank() {
		return errors.New("already exists")
	}

	allblank := true
	aroundSquares := b.Around(checkTarget)
	fmt.Println(aroundSquares)
	for _, aroundSquare := range aroundSquares {
		if !aroundSquare.IsBlank() {
			allblank = false
		}
	}
	if allblank {
		return errors.New("around is all blank")
	}

	return nil
}

func (b *Board) Around(square *Square) []*Square {
	type pos struct {
		x int
		y int
	}

	sqs := []*Square{}

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			px, py := square.pos.X+i, square.pos.Y+j
			// out of range
			if px < 0 || py < 0 || 7 < px || 7 < py {
				continue
			}
			sq := b[px][py] // check target
			sqs = append(sqs, sq)
		}
	}

	return sqs
}
