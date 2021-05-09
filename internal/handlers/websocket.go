package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)
type WsJsonResponse struct{
	Action string `json:"action"`
	Message string `json:"message"`
	MessageType string `json:"message_type"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool { return true},
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err!=nil{
		log.Println(err)
	}

	var res WsJsonResponse
	res.Message = "connection successfully upgraded"

	err = conn.WriteJSON(res)

	if err!=nil{
		log.Println(err)
	}

}