package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

type WebSocketConn struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Username string        `json:"username"`
	Action   string        `json:"action"`
	Message  string        `json:"message"`
	Conn     WebSocketConn `json:"-"`
}

var payloadChan = make(chan WsPayload)
var clients = make(map[WebSocketConn]string)

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
	clients[wsConn] = ""

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
			payload.Conn = *conn
			payloadChan <- payload
		}
	}
}

func ListenToPayloadChan() {
	var response WsJsonResponse

	for {
		payload := <-payloadChan

		//response.Action = "Got here"
		//response.Message = fmt.Sprintf("got some message and action was %s", payload.Action)
		switch payload.Action {
		case "username":
			clients[payload.Conn] = payload.Username
			userList:=getUserList()

			response.Action = "connected_users"
			response.ConnectedUsers = userList
			broadcastToAll(response)
		case "left":
			delete(clients, payload.Conn)

			response.Action = "connected_users"
			userList := getUserList()
			response.ConnectedUsers = userList
			broadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s",payload.Username,payload.Message)
			broadcastToAll(response)
		}
	}
}
func getUserList()[]string{
	var userList []string
	for _,val := range clients{
		if val!=""{
			userList = append(userList,val)
		}
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJsonResponse){
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("error occured in sending response")
			_ = client.Close()
			delete(clients, client)
		}
	}
}