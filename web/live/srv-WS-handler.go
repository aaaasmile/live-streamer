package live

import (
	"log"
	"net/http"

	"github.com/aaaasmile/live-streamer/web/live/ws"
	"github.com/gorilla/websocket"
)

var (
	upgrader  websocket.Upgrader
	wsClients *ws.WsClients
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS error", err)
		return
	}

	wsClients.AddConn(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error ", err)
			wsClients.CloseConn(conn)
			return
		}
		if messageType == websocket.TextMessage {
			log.Println("Message rec: ", string(p))
		}
	}
}

func WsHandlerShutdown() {
	wsClients.EndWS()
}

func init() {
	wsClients = ws.NewWsClients()
	wsClients.StartWS()
}
