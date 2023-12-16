package ws

import (
	"time"

	"golang.org/x/net/websocket"
)

const (
	writeWait = 10 * time.Second
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	inGame bool
}
