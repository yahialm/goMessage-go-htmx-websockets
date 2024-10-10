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

	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("Template parsing error: %s", err)
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("execution error: %s", err)
	}

	return renderedMessage.Bytes()

}
