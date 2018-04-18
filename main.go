package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Setting flags with default parameters
	// Can also be parsed while building/deploying the app
	var url, port string
	flag.StringVar(&url, "url", "localhost", "URL")
	flag.StringVar(&port, "port", ":9000", "Port")
	flag.Parse()

	// Setting up the hub & setup
	// Only 1 hub is active (for now) and shared across all connections
	hub := newHub()

	// Routing
	http.Handle("/ws", newWSConnection(hub))

	// Watching for messages
	go hub.handleMessages()

	// Starting the server
	log.Println("http server started on", url+port)

	err := http.ListenAndServe(url+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
