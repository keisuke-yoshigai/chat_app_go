package main

import (
    "github.com/gorilla/websocket"
)

type client struct {
    socket *websocket.Conn // このクライアントのためのWebSocket
    send chan []byte //メッセージが送るチャネル
    room *room //このクライアントが参加しているチャットルーム
}

func (c *client) read() {
    for {
        if _, msg , err := c.socket.ReadMessage(); err == nil {
            c.room.forward <- msg
        } else {
            break
        }
    }
    c.socket.Close()
}

func (c *client) write() {
    for msg := range c.send {
        if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
            break
        }
    }
    c.socket.Close()
}
