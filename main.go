package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"tictactoe.com/m/ws"
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

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	searchPool := make([]*ws.Player, 0)
	gamePool := make([]*ws.GameHub, 0)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(&searchPool, gamePool, w, r)
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
