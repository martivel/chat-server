package main

import "github.com/gorilla/websocket"

// Struct for containing the connection
// In the future it might contain other data
type Client struct {
	conn *websocket.Conn
}
