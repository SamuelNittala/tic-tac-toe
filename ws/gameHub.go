package ws

import (
	"github.com/google/uuid"
)

type GameHub struct {
	id        string
	player1   *Player
	player2   *Player
	broadcast chan string
	toMove    *Player
	gameState string
}

func NewGameHub(player1, player2 *Player) *GameHub {
	gameId := uuid.New().String()
	return &GameHub{
		id:        gameId,
		player1:   player1,
		player2:   player2,
		toMove:    player1,
		broadcast: make(chan string),
		gameState: gameId + "/.../.../.../",
	}
}

func (g *GameHub) isPlayerToMove(playerId string) string {
	if g.toMove.id == playerId {
		return "true"
	}
	return "false"
}

func (g *GameHub) run() {
	for {
		select {
		case message := <-g.broadcast:
			select {
			case g.player1.send <- []byte(message + g.isPlayerToMove(g.player1.id)):
			default:
				close(g.player1.send)
			}
			select {
			case g.player2.send <- []byte(message + g.isPlayerToMove(g.player2.id)):
			default:
				close(g.player2.send)
			}

		}
	}

}

func (g *GameHub) validateGame(id, gameState string) {
	// checking if the correct player is making the move
	if g.toMove.id == g.player1.id && g.player1.id == id {
		g.toMove = g.player2
		g.gameState = gameState
	}
	if g.toMove.id == g.player2.id && g.player2.id == id {
		g.toMove = g.player1
		g.gameState = gameState
	}
}
