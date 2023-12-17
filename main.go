package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	player1     Player
	player2     Player
	gameState   string
	sendChan    chan string
	receiveChan chan string
}

type Player struct {
	conn      *websocket.Conn
	inGame    bool
	searching bool
	id        string
}

type PlayerPool []Player

func createGame(player1, player2 Player) Game {
	return Game{
		id:          "1",
		player1:     player1,
		player2:     player2,
		gameState:   "1/" + player1.id + "/.../.../.../" + player1.id,
		sendChan:    make(chan string),
		receiveChan: make(chan string),
	}
}

func (game *Game) listenForMoves() {
	defer func() {
		game.player1.conn.Close()
		game.player2.conn.Close()
	}()
	for {
		_, byteState, err := game.player1.conn.ReadMessage()
		_, byteState2, err := game.player2.conn.ReadMessage()
		p1msg := string(byteState[:])
		p2msg := string(byteState2[:])
		fmt.Println(p1msg, p2msg)
		if err != nil {
			log.Println(err)
		}
		game_split_state := strings.Split(p1msg, "/")

		current_player := game_split_state[1]
		player_to_move := game_split_state[len(game_split_state)-1]

		// allow only if the current player is to move
		if current_player == player_to_move {
			if current_player == game.player1.id {
				game.gameState = strings.Join(game_split_state[:len(game_split_state)-1], "/") + "/" + game.player2.id
			} else {
				game.gameState = strings.Join(game_split_state[:len(game_split_state)-1], "/") + "/" + game.player1.id
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

func (game *Game) start(w http.ResponseWriter, r *http.Request) {
	go game.listenForMoves()
}

func search(w http.ResponseWriter, r *http.Request, playerPool []Player) Game {
	var avail_players []Player
	// find players
	for _, player := range playerPool {
		if !player.inGame && player.searching {
			avail_players = append(avail_players, player)
		}
	}
	return createGame(avail_players[0], avail_players[1])
}

func createPlayer(conn *websocket.Conn) Player {
	return Player{
		conn:   conn,
		id:     uuid.New().String(),
		inGame: false,
	}
}

func main() {
	flag.Parse()
	player_pool := make([]Player, 0)
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		player := createPlayer(conn)
		player_pool = append(player_pool, player)
		fmt.Println(player_pool)
		//listen
		_, message, err := conn.ReadMessage()
		if string(message) == "search" {
			player.searching = true
			fmt.Println(player_pool)
			// game := search(w, r, player_pool)
			// fmt.Println(game)
		}
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
