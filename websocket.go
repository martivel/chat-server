package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Improve cross origin fix below
		// For now expecting the app to run on localhost:8080
		// A better (CI/CD supportive) & multiple origins approach
		// might be added later
		if r.Header.Get("Origin") == "http://localhost:8080" {
			return true
		}
		return false
	},
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Making an struct that will inherit the Handler interface via
// the method pointer ServeHTTP
type Websocket struct {
	Hub *Hub
}

func (ws *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Make new client instance and add to clients map
	client := &Client{conn}
	ws.Hub.clients[client] = true
	fmt.Println("Client made a connection successfully")

	for {
		var msg Message

		// Read a string and pass to channel
		// _, msg, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Printf("error: %v", err)
		// 	delete(ws.Hub.clients, client)
		// 	break
		// }
		// ws.Hub.broadcast <- string(msg)

		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(ws.Hub.clients, client)
			break
		}

		// Send the received message to the broadcast channel
		fmt.Println(msg)
		ws.Hub.broadcast <- msg
	}
}

func newWSConnection(hub *Hub) *Websocket {
	return &Websocket{
		Hub: hub,
	}
}
