package main

import (
    "log"
	"net/http"
    "github.com/gorilla/websocket"
)

type room struct {
    forward chan []byte //他のクライアントへのメッセージを保持するチャネル
    join chan *client //チャットルームに参加しようとするクライアントを扱う
    leave chan *client //チャットルームに退出しようとするクライアントを扱う
    clients map[*client]bool //ルーム内に在室中のクライアントを保持する
}

const (
    socketBufferSize = 1024
    messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{//HTTP接続をアップグレードするためにwebsocket.Upgrader型を作成。websocket利用に必要。 
    ReadBufferSize: socketBufferSize,
    WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    socket, err := upgrader.Upgrade(w, req, nil) //websocketのコネクションを取得
    if err != nil {
        log.Fatal("ServeHTTP:", err)
        return
    }
    client := &client{
        socket: socket,
        send: make(chan []byte, messageBufferSize),
        room: r,
    }
    r.join <- client
    defer func() { r.leave <- client }()
    go client.write()
    client.read()
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
//func (r *room) run() {
//    for {
//        switch {
//        case client := <-r.join:
//            r.clients[client] = true
//        case client := <-r.leave:
//            delete(r.clients, client)
//            close(client.send)
//        case msg := <-r.forward:
//            for client := range r.clients {
//                select {
//                case client.send <- msg: //メッセージ送信
//                default: //メッセージ送信失敗
//                    delete(r.clients[client])
//                    close(client.send)
//                }
//            }
//        }
//    }
//}

func newRoom() *room {
    return &room{
        forward: make(chan []byte),
        join: make(chan *client),
        leave: make(chan *client),
        clients: make(map[*client]bool),
    }
}
