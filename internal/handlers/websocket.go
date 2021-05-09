package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketConn struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Username string        `json:"username"`
	Action   string        `json:"action"`
	Message  string        `json:"message"`
	Conn     WebSocketConn `json:"-"`
}

var payloadChan = make(chan WsPayload)
var clients = make(map[WebSocketConn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	var res WsJsonResponse
	res.Message = "connection successfully upgraded"

	wsConn := WebSocketConn{Conn: conn}
	clients[wsConn] = true

	err = conn.WriteJSON(res)

	if err != nil {
		log.Println(err)
	}
	go LoadPayload(&wsConn)
}

func LoadPayload(conn *WebSocketConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//no payload
		} else {
			payload.Conn =*conn
			payloadChan <- payload
		}
	}
}
