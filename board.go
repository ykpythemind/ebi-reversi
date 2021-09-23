package main

import (
	"bytes"
	"errors"
	"sync"
)

var NeedPassError = errors.New("need pass")
var GameEndError = errors.New("game finished")

// Board is 8x8 reversi board
type Board struct {
	Content [8][8]*Square
	mutex   sync.Mutex
}

func (b *Board) String() string {
	var buf bytes.Buffer
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sq := b.Content[i][j]
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

func (b *Board) gameCheck(nextPlayer Player) error {
	found := false

loop:
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			// すべてのマスにチェックして、次置ける場所があるか確認
			sq := b.Content[i][j]
			squares := b.Check(sq, nextPlayer)
			if len(squares) > 0 {
				found = true
				break loop
			}
		}
	}

	if found {
		return nil
	}

	gameEnd := true
	var nextnextPlayer Player
	if nextPlayer == PlayerB {
		nextnextPlayer = PlayerA
	} else {
		nextnextPlayer = PlayerB
	}

loop2:
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			// スキップされた結果、もう一度やることになるが置けるか。置けない場合はゲーム終了
			sq := b.Content[i][j]
			squares := b.Check(sq, nextnextPlayer)
			if len(squares) > 0 {
				gameEnd = false
				break loop2
			}
		}
	}

	if gameEnd {
		return GameEndError
	} else {
		return NeedPassError
	}
}

// Place places square
func (b *Board) Place(placeSquare *Square, player Player) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	squares := b.Check(placeSquare, player)
	if squares == nil || len(squares) == 0 {
		// ignore
		return nil
	}

	if placeSquare.state != Blank {
		// fmt.Printf("already placed... to place (%d,%d)\n", placeSquare.pos.X, placeSquare.pos.Y)
		// fmt.Printf("%+v\n", placeSquare)
		// fmt.Println(b)
		return errors.New("already placed")
	}

	for _, sq := range squares {
		if player == PlayerA {
			sq.state = PlayerAFilled
		} else { // B
			sq.state = PlayerBFilled
		}
	}

	if player == PlayerA {
		placeSquare.state = PlayerAFilled
	} else { // B
		placeSquare.state = PlayerBFilled
	}

	return nil
}

// Check checks whether player can place given square to the current board and return target squares to reverse.
func (b *Board) Check(input *Square, player Player) []*Square {
	if !input.IsBlank() {
		return nil
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
		return nil
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

	return result
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
			goto exit
		}

		sq := b.Content[px][py]
		if player == PlayerA {
			switch sq.state {
			case PlayerAFilled:
				if len(result) == 0 {
					goto exit
				} else {
					// 挟んだ
					return result
				}
			case PlayerBFilled:
				result = append(result, sq)
			case Blank:
				goto exit
			default:
				panic("unreachable")
			}
		} else if player == PlayerB {
			switch sq.state {
			case PlayerAFilled:
				result = append(result, sq)
			case PlayerBFilled:
				if len(result) == 0 {
					goto exit
				} else {
					// 挟んだ
					return result
				}
			case Blank:
				return nil
			default:
				panic("unreachable")
			}
		}
	}

exit:
	return nil
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
			sq := b.Content[px][py] // check target
			sqs = append(sqs, sq)
		}
	}

	return sqs
}
