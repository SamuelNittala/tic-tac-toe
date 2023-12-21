package ws

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Player struct {
	id      string
	conn    *websocket.Conn
	gameHub *GameHub
	send    chan []byte
}
type PlayerPool []*Player

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (p *Player) findOpponent(searchPool *[]*Player) *Player {
	// Check if there is at least one other player in the pool
	// fmt.Println("search pool", searchPool)
	if len(*searchPool) > 0 {
		for i, playerTwo := range *searchPool {
			fmt.Println(playerTwo, "p2")
			if playerTwo.id != p.id {
				// Found a match, remove playerTwo from the pool
				tempPool := *searchPool
				*searchPool = append(tempPool[:i], tempPool[i+1:]...)
				return playerTwo
			}
		}
	}
	// If no match found, add this player to the pool and return nil
	*searchPool = append(*searchPool, p)
	// fmt.Println("search pool", *searchPool)
	return nil
}

func (p *Player) readPump(searchPool *[]*Player) {
	defer func() {
		p.conn.Close()
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		stringMessage := string(message)
		if stringMessage == "search" && p.gameHub == nil {
			// fmt.Println(p.gameHub, " - gamehub for player")
			// fmt.Println("searching for" + p.id)
			opponent := p.findOpponent(searchPool)
			// fmt.Println(opponent, "opponent for"+p.id)
			if opponent != nil {
				// fmt.Println("oppenent for" + p.id + "is " + opponent.id)
				newGame := NewGameHub(p, opponent)
				p.gameHub = newGame
				opponent.gameHub = newGame
				go newGame.run()
				// fmt.Println("Game created" + newGame.id)
				p.gameHub.broadcast <- []byte(newGame.gameState) 
			}
		} else {
			if p.gameHub != nil {
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				currentPlayerToMove := p.gameHub.toMove
				// fmt.Println(currentPlayerToMove, p)
				p.gameHub.validateGame(p.id, string(message))
				p.gameHub.broadcast <- []byte(p.gameHub.gameState) 
			}
		}
		// c.hub.broadcast <- message
	}
}

func (p *Player) writePump() {
	defer func() {
		p.conn.Close()
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := p.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func ServeWs(searchPool *[]*Player, gamePool []*GameHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := &Player{
		conn: conn,
		id:   uuid.New().String(),
		send: make(chan []byte, 256),
	}
	// fmt.Println("Added player " + player.id)
	go player.readPump(searchPool)
	go player.writePump()
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	// go pl.writePump()
	// go client.readPump()
}
