package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)


type Client struct {
	id string
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}


func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection
	// create the client
	// listen to the hub
}