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

	var rows []string
	for _, r := range strings.Split(str, "\n") {
		if r != "" { // skip blank line
			rows = append(rows, r)
		}
	}

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
			case "-", "!":
				board.Content[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}}
			case "a":
				board.Content[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}, state: PlayerAFilled}
			case "b":
				board.Content[colindex][rowindex] = &Square{pos: SquarePosition{X: colindex, Y: rowindex}, state: PlayerBFilled}
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

func TestBoardCheck(t *testing.T) {
	temp := `
-,-,-,-,-,-,-,-
-,-,-,-,-,-,-,-
-,-,-,-,-,-,-,-
-,-,-,a,a,a,-,a
-,-,-,b,b,a,a,a
-,-,-,-,-,a,a,a
-,-,-,-,-,-,-,a
-,-,-,-,-,-,!,-`

	board, err := parseBoard(temp)
	if err != nil {
		t.Fatal(err)
	}
	sq := &Square{pos: SquarePosition{X: 6, Y: 7}}
	result := board.Check(sq, PlayerB)

	if len(result) > 0 {
		t.Fatalf("want: cannot reverse, got: %+v", result)
	}
}

func TestBoardCheck2(t *testing.T) {
	temp := `
-,-,-,-,-,-,-,-
-,-,-,-,-,-,-,-
-,-,-,-,-,-,-,-
-,-,-,a,a,a,-,-
-,-,-,b,b,a,-,-
-,-,-,-,-,a,-,-
-,-,-,-,-,-,-,-
-,-,-,-,-,-,!,-`

	board, err := parseBoard(temp)
	if err != nil {
		t.Fatal(err)
	}
	sq := &Square{pos: SquarePosition{X: 6, Y: 7}}
	result := board.Check(sq, PlayerB)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) > 0 {
		t.Fatalf("want: cannot reverse, got: %+v", result)
	}
}

func TestBoardCheck3(t *testing.T) {
	temp := `
-,-,-,-,-,-,-,-
-,-,-,-,-,-,-,-
-,-,-,-,-,!,-,-
-,-,-,a,a,a,-,-
-,-,-,b,b,a,-,-
-,-,-,-,-,a,-,-
-,-,-,-,-,b,-,-
-,-,-,-,-,b,-,-`

	board, err := parseBoard(temp)
	if err != nil {
		t.Fatal(err)
	}
	sq := &Square{pos: SquarePosition{X: 5, Y: 2}}
	{
		result := board.Check(sq, PlayerA)
		if len(result) > 0 {
			t.Fatalf("want: playerA cant reverse. got: %v", result)
		}
	}

	{
		result := board.Check(sq, PlayerB)
		expect := []struct {
			x int
			y int
		}{
			{x: 4, y: 3},
			{x: 5, y: 3},
			{x: 5, y: 4},
			{x: 5, y: 5},
		}

		for _, e := range expect {
			find := false
			for _, r := range result {
				find = r.pos.X == e.x && r.pos.Y == e.y
				if find {
					break
				}
			}
			if !find {
				t.Errorf("want (%d,%d), but not in results", e.x, e.y)
			}
		}
	}
}
