package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

type Client struct {
	ID               string
	Username         string
	ShowOnlineStatus string
	ShowUsername     bool
	Online           string
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action           string              `json:"action"`
	Message          string              `json:"message"`
	MessageType      string              `json:"message_type"`
	ConnectedClients []string            `json:"connected_clients"`
	Clients          map[string][]string `json:"clients"`
	From             string              `json:"from"`
	To               string              `json:"to"`
	ID               string              `json:"id"`
	Title            string              `json:"title"`
	Level            string              `json:"level"`
	Online           bool                `json:"online"`
	ShowOnlineStatus bool                `json:"show_online_status"`
	ShowUsername     bool                `json:"show_username"`
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
	Action           string              `json:"action"`
	Message          string              `json:"message"`
	Username         string              `json:"username"`
	ID               string              `json:"id"`
	Conn             WebSocketConnection `json:"-"`
	From             string              `json:"from"`
	To               string              `json:"to"`
	ShowOnlineStatus bool                `json:"showonlinestatus"`
	ShowUsername     bool                `json:"show_username"`
	Online           bool                `json:"online"`
}

func (m *Respository) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)

	println("Client connected to endpoint")

	if err != nil {
		log.Printf("Endpoint Error:\t%s", err.Error())
		return
	}

	var response WsJsonResponse

	response.Action = "initialconnectionconfirmed"

	conn := WebSocketConnection{Conn: ws}
	strMap := make(map[string]string)
	strMap["id"] = fmt.Sprintf("%v", conn.RemoteAddr())
	strMap["online"] = fmt.Sprintf("%t", false)
	strMap["username"] = ""
	strMap["showonlinestatus"] = fmt.Sprintf("%t", true)
	strMap["showusername"] = fmt.Sprintf("%t", true)
	clients[conn] = strMap

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
			delete(clients, e.Conn)

			clients := getConnectedClients()
			response.Clients = clients
			response.Action = "clients"
			broadcastToAll(response)

		case "thankyou":
			clients := getConnectedClients()
			response.Clients = clients
			response.Action = "clients"
			broadcastToAll(response)

		case "username":
			username := e.Message
			id := e.ID
			fmt.Printf("Client submitted their username and ID\n\t%s\n\t%s\n", username, id)
			checkUsernameExists(e.Conn, username, id)

		case "chat-users":
			getOnlineClients()

		case "broadcast":
			msg := e.Message
			from := e.From
			id := e.ID
			var client Client
			client = GetClient(id)

			fmt.Printf("\nReceived broadcast from client\n\tFrom: %s with ID %s\n\tMessage:%s\n", from, msg, id)

			response.Action = e.Action
			response.From = from
			response.Message = msg
			response.ShowUsername = client.ShowUsername
			broadcastToAll(response)

		case "leftroom":
			id := e.ID
			userLeft(id)
			clients := getConnectedClients()
			response.Clients = clients
			response.Action = "clients"
			broadcastToAll(response)

		case "showonlinestatus":
			id := e.ID
			showOnlineStatus := e.ShowOnlineStatus

			fmt.Printf("Client %s changed show online status to %t\n", id, showOnlineStatus)

			UpdateClientOnlineStatus(id, showOnlineStatus)

			clients := getConnectedClients()
			response.Clients = clients
			response.Action = "clients"
			broadcastToAll(response)
		}
	}

}

func UpdateClientOnlineStatus(id string, showOnlineStatus bool) {
	fmt.Printf("\n\tUpdateClientOnlineStatus method invoked with ID: %s\n", id)

	for c := range clients {
		dict := clients[c]

		if dict["id"] == id {
			dict["showonlinestatus"] = fmt.Sprintf("%t", showOnlineStatus)
			break
		}
	}
}

func GetClient(id string) Client {
	var client Client
	fmt.Printf("\n\tgetClient method invoked with ID: %s\n", id)

	for c := range clients {
		dict := clients[c]

		if dict["id"] == id {
			userId := dict["id"]
			username := dict["username"]
			online := dict["online"]
			showOnlineStatus := dict["showonlinestatus"]
			showUsername := dict["showusername"]
			show, _ := strconv.ParseBool(showUsername)

			client.ID = userId
			client.Username = username
			client.Online = online
			client.ShowOnlineStatus = showOnlineStatus
			client.ShowUsername = show
			break
		}
	}

	return client
}

func userLeft(id string) {
	for c := range clients {
		client := clients[c]

		if client["id"] == id {
			client["online"] = fmt.Sprintf("%t", false)
			return
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
	dict["online"] = fmt.Sprintf("%t", true)

	response.Action = "goodusername"
	response.Message = username
	sendToClient(conn, response)
}

/* func removeClient(conn WebSocketConnection, from, userId string) {
	// var response WsJsonResponse
	// delete(clients, conn)

	for client := range clients {
		c := clients[client]

		id := strings.ToLower(strings.TrimSpace(c["id"]))
		username := strings.ToLower(strings.TrimSpace(c["username"]))

		from := strings.ToLower(strings.TrimSpace(from))
		userId := strings.ToLower(strings.TrimSpace(userId))

		fmt.Printf("%s === %s and %s === %s\n", from, username, userId, id)

		if id == userId && username == from {
			fmt.Printf("User %s with ID %s left\n", username, id)
			delete(clients, conn)
			return
		}
	}
} */

func sendToClient(conn WebSocketConnection, response WsJsonResponse) {
	err := conn.WriteJSON(response)

	if err != nil {
		log.Println("Error send response to client:\t", err.Error())

		_ = conn.Close()

		delete(clients, conn)
	}
}

func getOnlineClients() {
	var response WsJsonResponse
	onlineClients := make(map[string][]string)

	for client := range clients {
		dict := clients[client]
		onlineClients[dict["id"]] = []string{dict["username"], dict["id"], dict["online"], dict["showonlinestatus"], dict["showusername"]}
	}

	response.Clients = onlineClients
	response.Action = "online-clients"
	broadcastToAll(response)
}

func getConnectedClients() map[string][]string {
	onlineClients := make(map[string][]string)

	for client := range clients {
		dict := clients[client]
		onlineClients[dict["id"]] = []string{dict["username"], dict["id"], dict["online"], dict["showonlinestatus"], dict["showusername"]}
	}

	return onlineClients
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
