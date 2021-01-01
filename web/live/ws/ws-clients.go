package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsClients struct {
	clients     map[*websocket.Conn]bool
	broadcastCh chan string
	mux         sync.Mutex
}

func NewWsClients() *WsClients {
	res := WsClients{
		clients:     make(map[*websocket.Conn]bool),
		broadcastCh: make(chan string),
	}
	return &res
}

func (wh *WsClients) Broadcast(msg string) {
	log.Println("Broadcast  ws msg", msg)
	wh.broadcastCh <- msg
}

func (wh *WsClients) AddConn(conn *websocket.Conn) {
	wh.mux.Lock()
	wh.clients[conn] = true
	wh.mux.Unlock()
	log.Println("New connection. Connected clients", conn.RemoteAddr(), len(wh.clients))
}

func (wh *WsClients) CloseConn(conn *websocket.Conn) {
	conn.Close()
	wh.RemoveConn(conn)
}

func (wh *WsClients) RemoveConn(conn *websocket.Conn) {
	wh.mux.Lock()
	delete(wh.clients, conn)
	wh.mux.Unlock()
	log.Println("Clients still connected ", len(wh.clients))
}

func (wh *WsClients) closeAllConn() {
	for conn := range wh.clients {
		wh.CloseConn(conn)
	}
}

func (wh *WsClients) listenBroadcastMsg() {
	log.Println("WS Waiting for broadcast")
	for {
		msg := <-wh.broadcastCh

		for conn := range wh.clients {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				wh.CloseConn(conn)
				log.Println("Socket error: ", err)
			}
		}
	}
}

func (wh *WsClients) StartWS() {
	go wh.listenBroadcastMsg()
}

func (wh *WsClients) EndWS() {
	log.Println("End of websocket service")
	wh.closeAllConn()
}
