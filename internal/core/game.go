package tictactoe

type CellValue int

const (
	Empty CellValue = iota
	Cross
	Circle
)

type Cell struct {
	x     int
	y     int
	value CellValue
}

type TicTacToe [][]Cell

func CreateGrid() TicTacToe {
	var grid TicTacToe
	for i := 0; i < 3; i++ {
		grid_row := [3]Cell{}
		for j := 0; j < 3; j++ {
			grid_row[j] = Cell{
				x:     i,
				y:     j,
				value: Empty,
			}
		}
		grid = append(grid, grid_row[:])
	}
	return grid
}

func (t TicTacToe) set(i int, j int, value CellValue) {
	t[i][j].value = value
}

func (t TicTacToe) isGridFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t[i][j].value != Empty {
				return false
			}
		}
	}
	return true
}

func (t TicTacToe) checkRowForSameValue(row int) bool {
	first_value := t[row][0].value
	for i := 1; i < 3; i++ {
		if t[row][i].value != first_value {
			return false
		}
	}
	return true
}

func (t TicTacToe) checkColForSameValue(col int) bool {
	first_value := t[0][col].value
	for i := 1; i < 3; i++ {
		if t[i][col].value != first_value {
			return false
		}
	}
	return true
}

func (t TicTacToe) checkLeftDiagonalsForSameValue() bool {
	first_value := t[0][0].value
	for i := 1; i < 3; i++ {
		for j := 1; j < 3; j++ {
			if t[i][j].value != first_value {
				return false
			}
		}
	}
	return true
}

func (t TicTacToe) checkRightDiagonalForSameValue() bool {
	// TODO: make it dynamic in future, if we increase the grid size
	first_value := t[0][2].value
	second_value := t[1][1].value
	third_value := t[2][0].value
	return first_value == second_value && first_value == third_value && second_value == third_value
}

func (t TicTacToe) checkForWin(CellValue) {

}
