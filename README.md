# SignalIO Golang Server Setup

This section provides instructions on setting up the server side of your WebSocket application using Golang with the `signal-io` library. The library enables real-time, bidirectional communication between the server and clients.

## Description
A lightweight library designed to enable low-latency, bidirectional, and event-based communication between client and server using WebSockets. This library simplifies real-time data exchange, making it ideal for applications requiring fast and reliable communication. It offers an intuitive API, seamless integration, and efficient handling of various data types, ensuring smooth and responsive interactions across your applications.

Note: For client-side communication, use the `socket-io-client` JavaScript library to connect to your server. This library provides the necessary functionality to handle WebSocket connections, events, and interactions from the client side.

## Installation

First, install the `signal-io` package using `go get`:

```bash
go get github.com/Syntax0xError/signal.io-golang
```

## Initialization

To initialize the WebSocket server in your Golang application, use the following code:

```go
package main

import (
    "github.com/Syntax0xError/signal.io-golang"
    "log"
)

func main() {
    // Initialize the WebSocket server on port 8080
    socket := signal.IOServer("8080")

    // Event handler for 'message' events
    socket.On("message", func(payload signal.Payload, client signal.Client) {
        // Handle incoming messages
        log.Printf("Received message from client %v: %v", client.ConnectionId(), payload)
        
        // Example: Send a response back to the client
        socket.Emit("response", "Message received")
    })
}

```
This will create a WebSocket server that listens on port `8080`, allowing clients to connect and communicate in real-time.

## Event Handling

### Event Registration
You can register event handlers using the `On` method. For example, to handle incoming messages:
```go
socket.On("message", func(payload signal.Payload, client signal.Client) {
    // `payload` is of type `interface{}` (or `any`), representing the data sent by the client
    // `client` provides information about the connected client

    // Handle the message here
    log.Printf("Received message from client %v: %v", client.ConnectionId(), payload)
    
    // Example: Send a response back to the client
    socket.Emit("response", "Message received")
})
```
### Payload Type
- `signal.Payload`: Represents the data sent from the client. It is of type interface{}, which is equivalent to any in other languages. This allows for flexible handling of various data types.

### Client Information
- `signal.Client`: Provides detailed information about the connected client, including:

    - `ConnectionId`: The unique identifier for the client connection.
    - `Auth`: The authentication token or credentials associated with the client.
    - `Query`: A map of query parameters sent during the connection initialization.
    - `Socket`: The WebSocket connection object (*websocket.Conn).
    - `HTTPRequest`: The HTTP request associated with the WebSocket connection (*http.Request).

### Emitting Messages
To send messages back to a client, use the Emit method on the client object:
```go
client.Emit("eventName", payload)
```
- `eventName`: The name of the event to send.
- `payload`: The data to send with the event. This can be any data type supported by interface{}.

Example usage:
```go
socket.On("message", func(payload signal.Payload, client signal.Client) {
    // Send a response back to the client
    client.Emit("response", "Message received")
})
```
This allows you to communicate specific responses or notifications back to the individual client.

### Broadcasting Messages to All Clients

To send a message to all connected clients, use the Broadcast method:

```go
socket.Broadcast(eventName, payload)
```
- `eventName`: The name of the event to send.
- `payload`: The data to send with the event.

Example usage:
```go
socket.On("broadcastToAll", func(eventName string, payload signal.Payload) {
    // Broadcast a message to all connected clients
    socket.Broadcast(eventName, payload)
    log.Printf("Broadcast message: %v", payload)
})
```
This method allows you to send messages to every connected client, useful for global updates or notifications.

## Room Management

### Joining a Room
To add a client to a specific room, use the JoinRoom method:
```go
socket.JoinRoom(roomId, client)
```
- `roomId`: The identifier for the room you want the client to join.
- `client`: The signal.Client object representing the client to be added to the room.

Example usage:
```go
socket.On("joinRoom", func(roomId string, client signal.Client) {
    // Add the client to the specified room
    socket.JoinRoom(roomId, client)
    log.Printf("Client %v joined room %v", client.ConnectionId, roomId)
})
```
This method allows you to manage rooms or groups of clients, facilitating organized communication within the WebSocket server.

### Emitting Messages to a Room

To send a message to all clients in a specific room, use the EmitTo method:
```go
socket.EmitTo(room, "message", payload)
```
- `room`: The identifier for the room where the message should be sent.
- `message`: The name of the event to send.
- `payload`: The data to send with the event.

Example usage:
```go
socket.On("broadcastToRoom", func(roomId string, message string, payload signal.Payload) {
    // Send a message to all clients in the specified room
    socket.EmitTo(roomId, message, payload)
    log.Printf("Message sent to room %v: %v", roomId, payload)
})
```
This method allows you to broadcast messages to all clients within a specific room or group, making it easy to send updates or notifications to multiple clients simultaneously.

## Donations and Sponsorships

If you find this library useful and want to support its ongoing development, you can contribute through donations or sponsorships. Your support helps me maintain and improve the library, add new features, and provide better support to the community.

### How to Donate

You can make a donation via the following platforms:

- **[GitHub Sponsors](https://github.com/sponsors/Syntax0xError)**: Support me directly through GitHub Sponsors. Contributions on GitHub help me fund development and cover project costs.
- **Cryptocurrency**: You can also support me through cryptocurrency. If you prefer to donate this way, you can use the following wallet addresses:

  - **Binance Coin (BNB):** `0x1E9890ac2f04F0B446D16ad1A26519c0a9535938` (BNB Chain network assets)
  - **Bitcoin (BTCB):** `0x1E9890ac2f04F0B446D16ad1A26519c0a9535938` (BNB Chain network assets)
  - **Ethereum (ETH):** `0x1E9890ac2f04F0B446D16ad1A26519c0a9535938` (BNB Chain network assets)

### Sponsorship Opportunities

We also offer sponsorship opportunities for organizations and businesses interested in supporting our work. As a sponsor, you can benefit from:

- **Prominent Recognition:** Your company's name and logo will be featured on our GitHub repository and project website.
- **Priority Support:** Receive priority support and dedicated assistance for integrating and using the library.
- **Custom Features:** Request custom features or enhancements tailored to your organization's needs.

For more information about sponsorship opportunities and benefits, please contact me at [oussemamri2013@gmail.com](mailto:oussemamri2013@gmail.com)

### Why Support Me?

Contributions and sponsorships help me to:

- **Improve the Library:** Develop new features, fix bugs, and enhance performance.
- **Provide Better Support:** Offer timely support and resolve issues more efficiently.
- **Maintain the Project:** Cover costs associated with hosting, development, and maintenance.

Thank you for considering supporting our project. Your contributions make a significant difference!