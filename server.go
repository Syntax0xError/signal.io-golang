/**
 * @module signal.io
 * @description
 * A lightweight library designed to enable low-latency, bidirectional, and event-based communication between client and server using WebSockets.
 * This library simplifies real-time data exchange, making it ideal for applications requiring fast and reliable communication.
 * It offers an intuitive API, seamless integration, and efficient handling of various data types, ensuring smooth and responsive interactions across your applications.
 *
 * @note Server-side support is currently limited to Golang.
 *
 * @author Oussema Amri <oussemamri2013@gmail.com>
 * @see {@link https://github.com/Syntax0xError} GitHub Repository
 */
package signal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define an upgrader to upgrade HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

func (socket *signalIO) On(eventName string, callback Event) {
	if socket.listeners == nil {
		socket.listeners = make(map[string]Event)
	}
	socket.listeners[eventName] = callback
}

func (socket *signalIO) init() {
	socket.connections = make([]Client, 0)
	socket.rooms = make(map[string][]Client)
}

func (socket *signalIO) Start() {
	socket.init()
	log.Println("SignalIO service has been started on port", socket.wsPort)
	// Define the WebSocket route
	http.HandleFunc("/", socket.handleConnections)

	err := http.ListenAndServe(":"+socket.wsPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func (socket *signalIO) GetTotalConnections() int {
	return len(socket.connections)
}
func createClient(ws *websocket.Conn, r *http.Request) (Client, error) {
	r.URL.Query()
	client := Client{
		ConnectionId: CreateConnectionId(),
		Socket:       ws,
		HTTPRequest:  r,
	}
	queryParams := r.URL.Query()

	client.Auth = queryParams.Get("auth")

	queryData := queryParams.Get("queryData")
	query, err := DecodeQueryData(queryData)
	if err != nil {
		return client, err
	}
	client.Query = query

	return client, nil
}

func (socket *signalIO) cleanup(connectionId string) {
	var wg sync.WaitGroup
	for roomId, clients := range socket.rooms {
		wg.Add(1)
		go func(roomId, connectionId string, clients []Client) {
			defer wg.Done()

			position := IndexOf(connectionId, clients)
			if position == -1 {
				return
			}

			// Lock the mutex before modifying the map
			socket.mu.Lock()
			defer socket.mu.Unlock()

			// Remove the client from the slice by index
			clients = append(clients[:position], clients[position+1:]...)

			// If the room is now empty, delete the room
			if len(clients) == 0 {
				delete(socket.rooms, roomId)
			} else {
				// Otherwise, update the room's clients
				socket.rooms[roomId] = clients
			}

		}(roomId, connectionId, clients)
	}
	wg.Wait()
}

func (socket *signalIO) removeConnection(connectionId string) {
	count := len(socket.connections)
	for index, client := range socket.connections {
		if client.ConnectionId == connectionId {
			// remove connection
			socket.connections[index] = socket.connections[count-1]
			socket.connections = socket.connections[:count-1]
			// disconnect user from rooms
			socket.cleanup(connectionId)
			return
		}
	}
}

func (socket *signalIO) onConnect(client Client) {
	socket.connections = append(socket.connections, client)
	onConnect := socket.listeners["connect"]
	if onConnect != nil {
		onConnect(nil, client)
	}
}

func (socket *signalIO) onDisconnect(client Client) {
	socket.removeConnection(client.ConnectionId)
	onDisconnect := socket.listeners["disconnect"]
	if onDisconnect != nil {
		onDisconnect(nil, client)
	}
}

func (socket *signalIO) onError(client Client, err error) {
	socket.removeConnection(client.ConnectionId)
	onError := socket.listeners["error"]
	if onError != nil {
		onError(err, client)
	}
}

func (socket *signalIO) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	client, err := createClient(ws, r)
	if err != nil {
		socket.onError(client, err)
		return
	}

	socket.onConnect(client)

	for {
		// Read a message from the client
		_, message, err := ws.ReadMessage()
		if err != nil {
			socket.onDisconnect(client)
			break
		}

		var msg Message

		err = json.Unmarshal(message, &msg)
		if err != nil {
			socket.onError(client, err)
			break
		}

		socket.processMessage(msg, client)
	}
}

func (socket *signalIO) processMessage(message Message, client Client) {
	if event, exists := socket.listeners[message.EventName]; exists {
		event(message.Payload, client)
	}
}

func (client *Client) Emit(eventName string, payload Payload) error {
	// Create the message struct with the event name and payload
	msg := Message{
		EventName: eventName,
		Payload:   payload,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return err
	}

	// Send the message to the client
	err = client.Socket.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		log.Printf("WriteMessage error: %v", err)
		return err
	}

	return nil
}

func (socket *signalIO) Broadcast(eventName string, payload Payload) {
	for _, client := range socket.connections {
		client.Emit(eventName, payload)
	}
}

func (socket *signalIO) JoinRoom(roomId string, client Client) {
	if socket.rooms[roomId] == nil {
		socket.rooms[roomId] = []Client{client}
		return
	}

	if IndexOf(client.ConnectionId, socket.rooms[roomId]) != -1 {
		return
	}

	socket.rooms[roomId] = append(socket.rooms[roomId], client)
}

func (socket *signalIO) EmitTo(roomId, eventName string, payload Payload) {
	room := socket.rooms[roomId]

	if room == nil {
		return
	}

	for _, client := range room {
		client.Emit(eventName, payload)
	}
}

func IOServer(WS_PORT string) *signalIO {
	server := signalIO{
		wsPort: WS_PORT,
	}
	return &server
}
