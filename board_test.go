package main

import (
	"fmt"
	"strings"
	"testing"
)

func parseBoard(str string) (*Board, error) {
	// -,-,-,-,-,-,-,-
	// -,a,-,-,-,-,-,-
	// -,-,a,-,-,-,-,-
	// -,-,-,a,b,a,-,-
	// -,-,-,b,a,-,-,-
	// -,-,a,-,-,-,-,-
	// -,-,a,-,-,-,-,-
	// -,-,-,-,-,-,-,-

	board := &Board{}

	rows := strings.Split(str, "\n")
	if len(rows) != 8 {
		return nil, fmt.Errorf("invalid row num %d", len(rows))
	}

	for rowindex, row := range rows {
		cols := strings.Split(row, ",")
		if len(cols) != 8 {
			return nil, fmt.Errorf("invalid col num %d", len(cols))
		}

		for colindex, elm := range cols {
			switch elm {
			case "-":
				board[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}}
			case "a":
				board[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}, state: PlayerAFilled}
			case "b":
				board[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}, state: PlayerBFilled}
			default:
				return nil, fmt.Errorf("invalid entity %s", elm)
			}
		}
	}

	return board, nil
}

func TestBoardParse(t *testing.T) {
	temp := `-,-,-,-,-,-,-,-
-,a,-,-,-,-,-,-
-,-,a,-,-,-,-,-
-,-,-,a,b,a,-,-
-,-,-,b,a,-,-,-
-,-,a,-,-,-,-,-
-,-,a,-,-,-,-,-
-,-,-,-,-,-,-,-`

	board, err := parseBoard(temp)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(board)
}

func TestReversiTarget(t *testing.T) {
}
