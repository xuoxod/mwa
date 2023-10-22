package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection][]string)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action           string                 `json:"action"`
	Message          string                 `json:"message"`
	MessageType      string                 `json:"message_type"`
	ConnectedClients []string               `json:"connected_clients"`
	Clients          map[string]interface{} `json:"clients"`
	From             string                 `json:"from"`
	To               string                 `json:"to"`
	ID               string                 `json:"id"`
}

type ClientResponse struct {
	Action         string                 `json:"action"`
	Message        string                 `json:"message"`
	MessageType    string                 `json:"message_type"`
	ConnectedUsers []string               `json:"connected_users"`
	Clients        map[string]interface{} `json:"clients"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	ID             string                 `json:"id"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	Username string              `json:"username"`
	Conn     WebSocketConnection `json:"-"`
}

func (m *Respository) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)

	println("Client connected to endpoint")

	if err != nil {
		log.Println(err)
	}

	var response WsJsonResponse

	response.Action = "initialconnection"

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = []string{fmt.Sprintf("%v", conn.RemoteAddr())}

	err = ws.WriteJSON(response)

	if err != nil {
		log.Println("Could not send initial response to client:\t", err.Error())
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "initialconnection":
			response.Action = "confirmed"
			response.ID = fmt.Sprintf("%v", e.Conn.RemoteAddr())
			sendToClient(e.Conn, response)

		case "left":
			removeClient(e.Conn)
			clients := getConnectedClients()
			response.ConnectedClients = clients
			response.Action = "clients"
			broadcastToAll(response)

		case "thankyou":
			clients := getConnectedClients()
			response.ConnectedClients = clients
			response.Action = "clients"
			broadcastToAll(response)
		}
	}

}

func removeClient(conn WebSocketConnection) {
	delete(clients, conn)
}

func sendToClient(conn WebSocketConnection, response WsJsonResponse) {
	err := conn.WriteJSON(response)

	if err != nil {
		log.Println("Error send response to client:\t", err.Error())

		_ = conn.Close()

		delete(clients, conn)
	}
}

func getConnectedClients() []string {
	members := []string{}

	for c := range clients {
		members = append(members, c.RemoteAddr().String())
	}

	return members
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		fmt.Println("client:\t", client)
		err := client.WriteJSON(response)

		if err != nil {
			log.Println("Web socket error")

			_ = client.Close()

			delete(clients, client)
		}
	}
}
