package ws

type Game struct {
	playerOne *Client
	playerTwo *Client
	gameData  []byte
}

type Hub struct {
	clients    map[*Client]bool
	games      []*Game
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if len(h.clients) >= 2 {
				var free_players []*Client
				for client := range h.clients {
					if !client.inGame {
						free_players = append(free_players, client)
					}
				}
				// get first two free players
				fp := free_players[0]
				sp := free_players[1]
				h.games = append(h.games, &Game{
					fp,
					sp,
					make([]byte, 0),
				})
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// find it the player is in a game
				for idx, game := range h.games {
					// remove if in game
					if game.playerOne == client || game.playerTwo == client {
						h.games = RemoveIndex(h.games, idx)
					}
				}
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
