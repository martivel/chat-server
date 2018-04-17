package main

import (
	"log"
)

// The Hub contains all the connected clients and the broadcast
// Broadcast is the channel that contains the messages between the clients
// In the future we can reuse the Hub struct to make multiple connections
type Hub struct {
	// All the connected clients
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message
}

// Define the message object. This is how the broadcast messages will look like
// These are filled via the client side
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan Message),
	}
}

func (h *Hub) handleMessages() {
	for {
		msg := <-h.broadcast

		for client := range h.clients {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.conn.Close()
				delete(h.clients, client)
			}
		}
	}
}
