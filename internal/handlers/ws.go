package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)

// var clients = make(map[WebSocketConnection][]string)
var clients = make(map[WebSocketConnection]map[string]string)

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
	Title            string                 `json:"title"`
	Level            string                 `json:"level"`
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
	ID       string              `json:"id"`
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
	strMap := make(map[string]string)
	strMap["id"] = fmt.Sprintf("%v", conn.RemoteAddr())
	clients[conn] = strMap

	for client := range clients {
		arr := clients[client]

		fmt.Println("Client ID:\t", arr["id"])
	}

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

		case "username":
			username := e.Message
			id := e.ID

			fmt.Printf("Client submitted their username and ID\n\t%s\n\t%s", username, id)

			checkUsernameExists(e.Conn, username, id)
		}
	}

}

func checkUsernameExists(conn WebSocketConnection, username string, id string) {
	var response WsJsonResponse

	for client := range clients {
		dict := clients[client]

		if dict["id"] != id {
			if username == dict["username"] {
				response.Action = "badusername"
				response.Title = "Username Error"
				response.Level = "error"
				response.Message = fmt.Sprintf("Username %s is already taken", username)
				sendToClient(conn, response)
				return
			}

		}
	}

	dict := clients[conn]
	dict["username"] = username

	response.Action = "goodusername"
	response.Message = username
	sendToClient(conn, response)
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
		err := client.WriteJSON(response)

		if err != nil {
			log.Println("Web socket error")

			_ = client.Close()

			delete(clients, client)
		}
	}
}
