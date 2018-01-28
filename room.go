package main

type room struct {
    forward chan []byte //他のクライアントへのメッセージを保持するチャネル
    join chan *client //チャットルームに参加しようとするクライアントを扱う
    leave chan *client //チャットルームに退出しようとするクライアントを扱う
    clients map[*client]bool //ルーム内に在室中のクライアントを保持する
}

func (r *room) run() {
    for {
        switch {
        case client := <-r.join:
            r.clients[client] = true
        case client := <-r.leave:
            delete(r.clients, client)
            close(client.send())
        case msg := <-forward:
            for client := range r.clients {
                select {
                case client.send <- msg: //メッセージ送信
                default: //メッセージ送信失敗
                    delete(r.clients[client])
                    close(client.send)
                }
            }
        }
    }
}
