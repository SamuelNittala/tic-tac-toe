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
}

func NewGameHub(player1, player2 *Player) *GameHub {
	gameId := uuid.New().String()
	return &GameHub{
		id:        gameId,
		player1:   player1,
		player2:   player2,
		toMove:    PlayerOne,
		broadcast: make(chan []byte),
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
