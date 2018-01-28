package main

import (
    "github.com/gorilla/websocket"
)

const (
    socketBufferSize = 1024
    messageBufferSize = 256
)

var upgradeer = &websocket.Upgrader{ReadBufferSize: //HTTP接続をアップグレードするためにwebsocket.Upgrader型を作成。websocket利用に必要。 
socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    socket, err := upgrader.Upgrader(w, req, nil) //websocketのコネクションを取得
    if err != nil {
        log.Fatal("ServeHTTP:", err)
        return
    }
    client := &client{
        socket: socket,
        send: make(chan []byte, messageBufferSizea),
        room: r,
    }
    r.join <- client
    defer func() { r.leave <- client }()
    go client.write()
    client.read()
}

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
    socket.Close()
}

func (c *client) write() {
    for msg := range c.send {
        if err := c.socket.WriteMessage(socket.TextMessage, msg); err != nil {
            break
        }
    }
    socket.Close()
}
