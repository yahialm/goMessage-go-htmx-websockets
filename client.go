package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

const (

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Max message size allowed from a peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := uuid.New()

	// create the client
	client := &Client{
		id:   id.String(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()

}

func (c *Client) readPump() {

	defer func() {
		c.conn.Close()
		c.hub.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	/*
		-> This function handles pong messages from the client.
		When the server sends a ping to check if the client is still connected,
		the client must respond with a pong message.
		-> If the server doesn't receive the pong within the time limit,
		it will assume the connection is dead.
		-> The handler function resets the read deadline whenever a pong is received, which effectively keeps the connection alive.
	*/

	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// While loop
	for {

		// Get the text received
		_, text, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &WSMessage{}
		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(msg); err != nil {
			log.Printf("error: %v", err)
		}

		c.hub.broadcast <- &Message{userId: c.id, message: msg.Text}

	}

}

func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.conn.Close()
	}()

	for {
		select {

		case msg, ok := <-c.send:

			// --> Sets a write deadline to ensure that if writing a message takes too long,
			// (e.g., if the network is slow), it times out after writeWait
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			
			// If the ok flag is false, the channel has been closed (meaning the client likely disconnected),
			// and the connection is terminated by sending a Close message to the WebSocket (websocket.CloseMessage)
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err := w.Write(msg); err != nil {
				log.Printf("error: %v", err)
			}

			w.Write(msg)

			n := len(c.send)
			for i:=0; i<n; i++ {
				nextMsg := <-c.send
				w.Write(nextMsg)
			}

			if err:=w.Close(); err != nil {
				return
			}

		// You are receiving a value from ticker.C, but you are not storing that value in any variable.
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}
