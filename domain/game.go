package domain

import (
	"strings"
)

type CellValue int

const (
	Empty CellValue = iota
	Cross
	Circle
)

type TicTacToe [][]CellValue

// func CreateGrid() TicTacToe {
// 	var grid TicTacToe
// 	for i := 0; i < 3; i++ {
// 		grid_row := [3]Cell{}
// 		for j := 0; j < 3; j++ {
// 			grid_row[j] =
// 		}
// 		grid = append(grid, grid_row[:])
// 	}
// 	return grid
// }

func (t TicTacToe) isGridFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t[i][j] != Empty {
				return false
			}
		}
	}
	return true
}

func (t TicTacToe) checkRowForSameValue(row int) bool {
	first_value := t[row][0]
	for i := 1; i < 3; i++ {
		if t[row][i] != first_value {
			return false
		}
	}
	return true
}

func (t TicTacToe) checkColForSameValue(col int) bool {
	first_value := t[0][col]
	for i := 1; i < 3; i++ {
		if t[i][col] != first_value {
			return false
		}
	}
	return true
}

func (t TicTacToe) checkLeftDiagonalsForSameValue() bool {
	first_value := t[0][0]
	for i := 1; i < 3; i++ {
		for j := 1; j < 3; j++ {
			if t[i][j] != first_value {
				return false
			}
		}
	}
	return true
}

func (t TicTacToe) checkRightDiagonalForSameValue() bool {
	// TODO: make it dynamic in future, if we increase the grid size
	first_value := t[0][2]
	second_value := t[1][1]
	third_value := t[2][0]
	return first_value == second_value && first_value == third_value && second_value == third_value
}

func (t TicTacToe) checkForWin(CellValue) {}

func ParseGameState(gameState string) TicTacToe {
	splitGameState := strings.Split(gameState, "/")

	firstRow := splitGameState[1]
	secondRow := splitGameState[2]
	thirdRow := splitGameState[3]

	parsedGameState := make(TicTacToe, 0)
	matrix := []string{firstRow, secondRow, thirdRow}

	for _, row := range matrix {
		currCellRow := []CellValue{}
		for _, col := range row {
			switch col {
			case '.':
				currCellRow = append(currCellRow, Empty)
			case '*':
				currCellRow = append(currCellRow, Cross)
			case 'o':
				currCellRow = append(currCellRow, Circle)
			}
		}
		parsedGameState = append(parsedGameState, currCellRow)
	}
	return parsedGameState
}
