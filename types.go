package signal

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Payload any

type Message struct {
	EventName string  `json:"eventName"`
	Payload   Payload `json:"payload"`
}

type Event = func(Payload, Client)

type signalIO struct {
	wsPort      string
	listeners   map[string]Event
	connections []Client
	rooms       map[string][]Client

	mu sync.Mutex
}

type Client struct {
	connectionId string
	auth         string
	query        map[string]string
	socket       *websocket.Conn
	request      *http.Request
}
