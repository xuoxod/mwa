package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
	"github.com/xuoxod/mwa/pkg/utils"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)
var clientList = make(map[string]interface{})
var println = utils.Print

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
	clients[conn] = ""

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
		case "username":
			// get a list of all users and send it back via broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>:\t%s", e.Username, e.Message)
			response.From = e.Username
			broadcastToAll(response)

		case "initialconnection":
			response.Action = "clients"
			users := getUserList()
			response.ConnectedUsers = users
			response.Message = "You're connected"
			broadcastToAll(response)

		case "ipaddress":
			ip := e.Message
			clientDetails := make(map[string]interface{})
			clientDetails["address"] = ip
			clientDetails["conn"] = e.Conn
			clientList[ip] = clientDetails
			fmt.Println("Received Client's IP address:\t", ip)

			response.Action = "confirmed"
			response.ID = ip
			broadcastToClient(&e.Conn, response, ip)

			clients[e.Conn] = ip
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		}

		// response.Action = "Got here"
		// response.Message = fmt.Sprintf("Some message and action was %s ", e.Action)
		// broadcastToAll(response)
	}
}

func getUserList() []string {
	var userList []string

	for _, c := range clients {
		if c != "" {
			userList = append(userList, c)
		}
	}

	sort.Strings(userList)

	return userList
}

func broadcastToClient(client *WebSocketConnection, response WsJsonResponse, id string) {
	err := client.WriteJSON(response)

	if err != nil {

		log.Println("Web socket error")

		_ = client.Close()

		delete(clientList, id)
	}
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
