package main

import (
	"bytes"
	"log"
	"text/template"
)

type Hub struct {
	clients map[*Client]bool

	messages []*Message

	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    map[*Client]bool{},
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			log.Printf("user with ID: %s registred", client.id)

			for _, msg := range h.messages {
				client.send <- getMessageTemplate(msg)
			}


		case client := <- h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Printf("user with ID: %s unregistred", client.id)
			}

		case msg := <- h.broadcast:
			h.messages = append(h.messages, msg)

			for client := range h.clients {
				select {
					case client.send <- getMessageTemplate(msg):
						default:
							close(client.send)
							delete(h.clients, client) 
				}
			}

		}
	}
}


func getMessageTemplate(msg *Message) []byte {

	// This function loads the template file (message.html) and parses it.
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("Template parsing error: %s", err)
	}

	// This creates a buffer (a temporary in-memory storage area) to hold the rendered template.
	// Once the placeholders in the template are replaced with the actual message data,
	// the final HTML will be written to this buffer.
	var renderedMessage bytes.Buffer

	// This line actually processes the template and replaces 
	// placeholders (like {{ .UserId }} and {{ .Message }}) with values from the msg struct
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("execution error: %s", err)
	}

	
	return renderedMessage.Bytes()

}
