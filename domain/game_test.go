package domain

import (
	"testing"
)

func TestGridRows(t *testing.T) {
	//row test
	gameState := "id/ooo/.../.../"
	parsedGameState := GetGameState(gameState)
	if parsedGameState != PlayerTwoWon {
		t.Fail()
	}
	// row 2 test
	gameState = "id/.../***/.../"
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
	// row 3 test
	gameState = "id/.../.../***/"
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
}

func TestGridColumns(t *testing.T) {
	// col-1 test
	gameState := "id/*../*.../*..."
	parsedGameState := GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
	// col-2 test
	gameState = "id/.o./.o./.o*"
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerTwoWon {
		t.Fail()
	}
	// col-3 test
	gameState = "id/.o*/oo*/..*"
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
}

func TestLeftDiagonal(t *testing.T) {
	gameState := "id/*../.*./..*"
	parsedGameState := GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
	gameState = "id/o../.o./..o"
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerTwoWon {
		t.Fail()
	}
}

func TestRightDiagonal(t *testing.T) {
	gameState := "id/..*/.*./*.."
	parsedGameState := GetGameState(gameState)
	if parsedGameState != PlayerOneWon {
		t.Fail()
	}
	gameState = "id/..o/.o./o.."
	parsedGameState = GetGameState(gameState)
	if parsedGameState != PlayerTwoWon {
		t.Fail()
	}
}
