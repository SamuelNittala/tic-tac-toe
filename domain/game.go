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

func (t TicTacToe) isGridFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

func (t TicTacToe) checkRowsForSameValue() CellValue {
	for r := 0; r < 3; r++ {
		firstValue := t[r][0]
		if firstValue == Empty {
			continue
		}
		valid := true
		for c := 1; c < 3; c++ {
			if t[r][c] != firstValue {
				valid = false
			}
		}
		if valid {
			return firstValue
		}
	}
	return Empty
}

func (t TicTacToe) checkColsForSameValue() CellValue {
	for c := 0; c < 3; c++ {
		firstValue := t[0][c]
		if firstValue == Empty {
			continue
		}
		valid := true
		for r := 1; r < 3; r++ {
			if t[r][c] != firstValue {
				valid = false
			}
		}
		if valid {
			return firstValue
		}
	}
	return Empty
}

func (t TicTacToe) checkLeftDiagonalsForSameValue() CellValue {
	firstValue := t[0][0]
	secondValue := t[1][1]
	thirdValue := t[2][2]
	if firstValue == secondValue && firstValue == thirdValue && secondValue == thirdValue {
		return firstValue
	}
	return Empty
}

func (t TicTacToe) checkRightDiagonalForSameValue() CellValue {
	// TODO: make it dynamic in future, if we increase the grid size
	first_value := t[0][2]
	second_value := t[1][1]
	third_value := t[2][0]
	if first_value == second_value && first_value == third_value && second_value == third_value {
		return first_value
	}
	return Empty
}

func (t TicTacToe) checkForWin() GameResult {
	// check each row for win
	// row - 0
	rowCheck := t.checkRowsForSameValue()
	if rowCheck != Empty {
		if rowCheck == Cross {
			return PlayerOneWon
		}
		return PlayerTwoWon
	}
	colCheck := t.checkColsForSameValue()
	if colCheck != Empty {
		if colCheck == Cross {
			return PlayerOneWon
		}
		return PlayerTwoWon
	}
	leftDiagCheck := t.checkLeftDiagonalsForSameValue()
	if leftDiagCheck != Empty {
		if leftDiagCheck == Cross {
			return PlayerOneWon
		}
		return PlayerTwoWon
	}
	rightDiagCheck := t.checkRightDiagonalForSameValue()
	if rightDiagCheck != Empty {
		if rightDiagCheck == Cross {
			return PlayerOneWon
		}
		return PlayerTwoWon
	}
	if t.isGridFull() {
		return Tie
	}
	return Continue
}

type GameResult int

const (
	Continue GameResult = iota
	PlayerOneWon
	PlayerTwoWon
	Tie
)

func GetGameState(gameState string) GameResult {
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
	return parsedGameState.checkForWin()
}
