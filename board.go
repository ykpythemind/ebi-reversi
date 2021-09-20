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

// Place places square
func (b *Board) Place(placeSquare *Square, player Player) error {
	squares, err := b.Check(placeSquare, player)
	if err != nil {
		return err
	}
	if len(squares) == 0 {
		// ignore
		return nil
	}

	if placeSquare.state != Blank {
		return errors.New("already placed")
	}

	if player == PlayerA {
		placeSquare.state = PlayerAFilled
	} else { // B
		placeSquare.state = PlayerBFilled
	}

	for _, sq := range squares {
		if player == PlayerA {
			sq.state = PlayerAFilled
		} else { // B
			sq.state = PlayerBFilled
		}
	}

	return nil
}

// Check checks whether player can place given square to the current board and return target squares to reverse.
func (b *Board) Check(input *Square, player Player) ([]*Square, error) {
	if !input.IsBlank() {
		return nil, errors.New("already exists")
	}

	allblank := true
	aroundSquares := b.around(input)
	for _, aroundSquare := range aroundSquares {
		if !aroundSquare.IsBlank() {
			allblank = false
		}
	}
	if allblank {
		// around is all blank. so early break
		return nil, nil
	}

	var result []*Square

	// check
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue // self
			}
			r := b.findTargetWithDirection(input, player, i, j)
			if r != nil && len(r) > 0 {
				result = append(result, r...)
			}
		}
	}

	fmt.Println("targets:", result)
	return result, nil
}

func (b *Board) findTargetWithDirection(input *Square, player Player, dx, dy int) []*Square {
	// dx,dyを探索方向としてリバースできるマスを集めてくる
	var result []*Square
	px, py := input.pos.X, input.pos.Y

	for {
		px += dx
		py += dy
		// out of range, so finish
		if px < 0 || py < 0 || 7 < px || 7 < py {
			return nil
		}

		sq := b[px][py]
		if player == PlayerA {
			switch sq.state {
			case PlayerAFilled:
				if len(result) == 0 {
					return nil
				} else {
					// 挟んだ
					return result
				}
			case PlayerBFilled:
				result = append(result, sq)
			case Blank:
				return nil
			default:
				panic("unreachable")
			}
		} else if player == PlayerB {
			switch sq.state {
			case PlayerAFilled:
				result = append(result, sq)
			case PlayerBFilled:
				if len(result) == 0 {
					return nil
				} else {
					// 挟んだ
					return result
				}
			case Blank:
				return nil
			default:
				panic("unreachable")
			}
		} else {
			panic("unreachable")
		}
	}

	//FIXME: unreachable code
	return result
}

func (b *Board) around(square *Square) []*Square {
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
