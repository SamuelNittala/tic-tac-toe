package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/*
	Game = [ update_chan, move_chan ]

	// player 1 makes a move through receive_chan
	// player 1 should not be able to make a move until player 2 makes a move and vice-versa

	// need a way to find who is making the move -> through a playerId?
	// game_id/player_id/row1/row2/row3

	// step - 1 -> player 1 makes a move
		move_chan <- gameState
		// check if player 1 is the one that needs to make a move(how?)
		// store it in the gameState like (game_id/player_id/row_1/row_2/row_3/move_player_id)

	p1 -> rc
*/

type Game struct {
	id          string
	player1     string
	player2     string
	gameState   string
	sendChan    chan string
	receiveChan chan string
	conn        *websocket.Conn
}

func createGame(player1, player2 string) Game {
	return Game{
		id:          "1",
		player1:     player1,
		player2:     player2,
		gameState:   "1/" + player1 + "/.../.../.../" + player1,
		sendChan:    make(chan string),
		receiveChan: make(chan string),
	}
}

func (game *Game) listenForMoves() {
	defer func() {
		game.conn.Close()
	}()
	for {
		_, byteState, err := game.conn.ReadMessage()
		gameState := string(byteState[:])
		fmt.Println(gameState)
		if err != nil {
			log.Println(err)
		}
		game_split_state := strings.Split(gameState, "/")

		current_player := game_split_state[1]
		player_to_move := game_split_state[len(game_split_state)-1]

		// allow only if the current player is to move
		if current_player == player_to_move {
			if current_player == game.player1 {
				game.gameState = strings.Join(game_split_state[:len(game_split_state)-1], "/") + "/" + game.player2
			} else {
				game.gameState = strings.Join(game_split_state[:len(game_split_state)-1], "/") + "/" + game.player1
			}
		}
		game.sendChan <- game.gameState
	}
}

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func gameLoop(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	game := createGame("a", "b")
	game.conn = conn
	go game.listenForMoves()
	x, ok := <-game.sendChan
	if ok {
		fmt.Println(x)
	}

}

func main() {
	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gameLoop(w, r)
	})
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
