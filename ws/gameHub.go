package ws

import (
	"github.com/google/uuid"
)

type ToMove int

const (
	PlayerOne ToMove = iota
	PlayerTwo
)

type GameHub struct {
	id        string
	player1   *Player
	player2   *Player
	broadcast chan []byte
	toMove    ToMove
	gameState string
}

func NewGameHub(player1, player2 *Player) *GameHub {
	gameId := uuid.New().String()
	return &GameHub{
		id:        gameId,
		player1:   player1,
		player2:   player2,
		toMove: PlayerOne,,
		broadcast: make(chan []byte),
		gameState: gameId + "/created/.../.../.../",
	}
}

func (g *GameHub) run() {
	for {
		select {
		case message := <-g.broadcast:
			select {
			case g.player1.send <- message:
			default:
				close(g.player1.send)
			}
			select {
			case g.player2.send <- message:
			default:
				close(g.player2.send)
			}

		}
	}

}

func (g *GameHub) validateGame(id, gameState string) {
	// checking if the correct player is making the move
	if g.toMove == PlayerOne && g.player1.id == id {
		g.toMove = PlayerTwo
		g.gameState = gameState
	}
	if g.toMove == PlayerTwo && g.player2.id == id {
		g.toMove = PlayerOne
		g.gameState = gameState
	}
}
