package ports

type GameService interface {
	CreateGame() (string, error)
	updateGame(gameData string) string
	validateGame(gameData string) string
}

type SocketHandler struct {
	gameService GameService
}

func NewSocketHandler(gs GameService) *SocketHandler {
	return &SocketHandler{
		gameService: gs,
	}
}

func (srv *SocketHandler) CreateGame(send chan []byte) {
	game, err := srv.gameService.CreateGame()
	if err != nil {
		return
	}
	send <- []byte(game)
}
